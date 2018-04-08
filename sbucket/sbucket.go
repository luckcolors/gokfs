package sbucket

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/luckcolors/gokfs/opt"
	"github.com/luckcolors/gokfs/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type SBucket struct {
	DBPath string
	Opt    *opt.Options
	DB     *leveldb.DB
}

type Stats struct {
	UsedSpace int64
	FreeSpace int64
}

type List []KeyStat

type KeyStat struct {
	BaseKey         []byte
	ApproximateSize int64
}

var ErrSBucketOpeningError = errors.New("Error while creating leveldb instance")
var ErrSBucketOnStat = errors.New("Error while getting leveldb stats")

func New(DBPath string, opt *opt.Options) *SBucket {
	return &SBucket{
		DBPath: DBPath,
		Opt:    opt,
		DB:     nil,
	}
}

func (sb SBucket) Open() (err error) {
	sb.DB, err = leveldb.OpenFile(sb.DBPath, sb.Opt.LevelDBOpt)
	if err != nil {
		err = errors.Wrap(err, ErrSBucketOpeningError.Error())
	}
	return
}

func (sb SBucket) Close() error {
	return sb.DB.Close()
}

func (sb SBucket) Exists(key []byte) (bool, error) {
	key, err := utils.CreateItemKeyFromIndex(key, 0)
	if err != nil {
		return false, err
	}
	return sb.DB.Has(key, nil)
}

func (sb SBucket) Unlink(key []byte) error {
	key, err := utils.CreateItemKeyFromIndex(key, 0)
	if err != nil {
		return err
	}

	var index = 0
	var isFound = true

	for isFound {
		_, err = sb.DB.Get(key, nil)
		index++

		if err == nil {
			err = sb.DB.Delete(key, nil)
			if err != nil {
				return err
			}
		}
		if err == leveldb.ErrNotFound {
			isFound = false
		}
	}
	return nil
}

func (sb SBucket) Read() ([]byte, error) {
	return nil, nil
}

func (sb SBucket) WriteFile() {

}

func (sb SBucket) CreateFileReader() {

}

func (sb SBucket) CreateFileWriter() {

}

func (sb SBucket) Stat() (*Stats, error) {
	sizes, err := sb.DB.SizeOf([]util.Range{*sb.Opt.LevelDBIterRange})
	if err != nil {
		return nil, errors.Wrap(err, ErrSBucketOnStat.Error())
	}

	var i int64
	for _, s := range sizes {
		i += s
	}
	return &Stats{
		UsedSpace: i,
		FreeSpace: sb.Opt.MaxSize - i,
	}, nil
}

func (sb SBucket) List() (*List, error) {
	iterator := sb.DB.NewIterator(sb.Opt.LevelDBIterRange, sb.Opt.LevelDBIterReadOpt)

	var currentResult []byte
	var keys = make(map[string]int64)

	for !iterator.Next() && iterator.Error() == nil {
		currentResult = bytes.Split(iterator.Key(), sb.Opt.KeySeparator)[0]
		keys[string(currentResult)] += sb.Opt.ChunkSize
	}
	if iterator.Error() != nil {
		return nil, iterator.Error()
	}

	var list = make(List, len(keys))
	var i = 0
	for key, aprxSize := range keys {
		list[i] = KeyStat{
			BaseKey:         []byte(key),
			ApproximateSize: int64(aprxSize),
		}
		i++
	}
	return &list, nil
}

func (sb SBucket) Flush() error {
	return nil
}
