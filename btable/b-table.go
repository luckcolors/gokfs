package btable

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"path"
	"sync"

	"github.com/pkg/errors"

	"github.com/luckcolors/gokfs/constants"
	"github.com/luckcolors/gokfs/opt"
	"github.com/luckcolors/gokfs/sbucket"
	"github.com/luckcolors/gokfs/utils"
)

var ErrReadingRid = errors.New("Error while reading r.id file")
var ErrSBucketIndexGreaterThanB = errors.New("SBucket index bigger than B")
var ErrSBucketIndexNotLessthanZero = errors.New("SBucket index not less than zero")

type BTable struct {
	sBuckets  map[uint64]*sbucket.SBucket
	TablePath string
	destroy   chan uint64
	Opt       *opt.Options
	rwLock    *sync.RWMutex
}

type Stats []Stat

type Stat struct {
	SBucketIndex uint64
	Stat         *sbucket.Stats
}

type Lists []sbucket.List

func New(opt *opt.Options) *BTable {
	var bt = &BTable{
		sBuckets:  make(map[uint64]*sbucket.SBucket),
		TablePath: opt.TablePath,
		Opt:       opt,
		rwLock:    new(sync.RWMutex),
	}

	if opt.AutoCloseTimeout >= 0 {
		go bt.handleSbucketExpireLoop()
	}
	return bt
}

func (b BTable) Open() error {
	var ridPath = path.Join(b.Opt.TablePath, "r.id")
	if !utils.FileDoesExist(ridPath) {
		utils.InitBtableDirectory(b.Opt.TablePath, b.Opt.Rid)
	} else {
		utils.ValidateTablePath(b.Opt.TablePath)
	}
	rid, err := ioutil.ReadFile(ridPath)
	b.Opt.Rid = rid
	if err != nil {
		return errors.Wrap(err, ErrReadingRid.Error())
	}
	return nil
}

func (b BTable) HandleSBucketAtIndex(sBucketIndex uint64, destroyIndex bool) (*sbucket.SBucket, error) {
	if sBucketIndex > constants.B {
		return nil, ErrSBucketIndexGreaterThanB
	}
	if sBucketIndex < 0 {
		return nil, ErrSBucketIndexNotLessthanZero
	}

	b.rwLock.Lock()
	defer b.rwLock.Unlock()

	// Handle deletion
	if destroyIndex {
		delete(b.sBuckets, sBucketIndex)
		return nil, nil
	}

	if b.sBuckets[sBucketIndex] != nil {
		return b.sBuckets[sBucketIndex], nil
	}

	// Handle creation
	b.sBuckets[sBucketIndex] = sbucket.New(path.Join(b.TablePath, utils.CreateSbucketNameFromIndex(sBucketIndex)), b.Opt)
	b.sBuckets[sBucketIndex].Open()

	return b.sBuckets[sBucketIndex], nil
}

func (b BTable) handleSbucketExpireLoop() {
	for {
		select {
		case index, _ := <-b.destroy:
			b.HandleSBucketAtIndex(index, true)
		}
	}
}

func (b BTable) CalculateSBucketIndexForKey(key []byte) uint64 {
	//xorDist := make([]byte, constants.D/8)
	var xorDist [8]byte
	for i := 0; i < int(constants.D/8); i++ {
		xorDist[i] = b.Opt.Rid[i] ^ utils.HashKey(key)[i]
	}
	var out uint64
	_ = binary.Read(bytes.NewBuffer(xorDist[:]), binary.LittleEndian, &out)
	return out
}

func (b BTable) GetSBucketForKey(key []byte) (*sbucket.SBucket, error) {
	sBucketIndex := b.CalculateSBucketIndexForKey(key)

	sBucket, err := b.HandleSBucketAtIndex(sBucketIndex, false)
	if err != nil {
		return nil, err
	}
	//Todo Opening
	return sBucket, nil
}

func (b BTable) StatWithKey(key []byte) (*Stat, error) {
	return b.StatWithIndex(b.CalculateSBucketIndexForKey(key))
}

func (b BTable) StatWithIndex(sBucketIndex uint64) (*Stat, error) {
	s, err := b.HandleSBucketAtIndex(sBucketIndex, false)
	if err != nil {
		return nil, err
	}

	stat, err := s.Stat()
	if err != nil {
		return nil, err
	}

	return &Stat{
		SBucketIndex: sBucketIndex,
		Stat:         stat,
	}, nil
}

func (b BTable) Stat() (*Stats, error) {
	defer b.rwLock.RUnlock()
	b.rwLock.RLock()

	sBucketIndexes, err := utils.ExistingSbucketIndexes(b.TablePath)
	if err != nil {
		return nil, err
	}

	var stats = make(Stats, len(*sBucketIndexes))
	var i = 0
	for _, sbi := range *sBucketIndexes {
		s, err := b.sBuckets[sbi].Stat()
		if err != nil {
			return nil, err
		}
		stats[i].SBucketIndex = sbi
		stats[i].Stat = s
	}

	return &stats, nil
}

func (b BTable) ListWithKey(key []byte) (*sbucket.List, error) {
	defer b.rwLock.RUnlock()
	b.rwLock.RLock()

	sb, err := b.GetSBucketForKey(key)
	if err != nil {
		return nil, err
	}
	return sb.List()
}

func (b BTable) ListWithIndex(sBucketIndex uint64) (*sbucket.List, error) {
	defer b.rwLock.RUnlock()
	b.rwLock.RLock()

	sb, err := b.HandleSBucketAtIndex(sBucketIndex, false)
	if err != nil {
		return nil, err
	}
	return sb.List()
}

func (b BTable) List() (*Lists, error) {
	defer b.rwLock.RUnlock()
	b.rwLock.RLock()

	sBucketIndexes, err := utils.ExistingSbucketIndexes(b.TablePath)
	if err != nil {
		return nil, err
	}

	var lists = make(Lists, len(*sBucketIndexes))
	var i = 0
	for _, sbi := range *sBucketIndexes {
		l, err := b.sBuckets[sbi].List()
		if err != nil {
			return nil, err
		}
		lists[i] = *l
		i++
	}

	return &lists, nil
}

func (b BTable) Exist(key []byte) (bool, error) {
	sb, err := b.GetSBucketForKey(key)
	if err != nil {
		return false, nil
	}
	return sb.Exists(key)
}

func (b BTable) Unlink(key []byte) (bool, error) {
	sb, err := b.GetSBucketForKey(key)
	if err != nil {
		return false, nil
	}
	return true, sb.Unlink(key)
}

func (b BTable) Flush() error {
	defer b.rwLock.RUnlock()
	b.rwLock.RLock()

	sBucketIndexes, err := utils.ExistingSbucketIndexes(b.TablePath)
	if err != nil {
		return err
	}

	for _, sbi := range *sBucketIndexes {
		err := b.sBuckets[sbi].Flush()
		if err != nil {
			return err
		}
	}

	return nil
}
