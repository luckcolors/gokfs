package opt

import (
	"time"

	ldbOpt "github.com/syndtr/goleveldb/leveldb/opt"
	ldbUtil "github.com/syndtr/goleveldb/leveldb/util"
)

type Options struct {
	Path                string
	TablePath           string
	Rid                 []byte
	AutoCloseTimeout    time.Duration
	MaxTableSize        int64
	MaxSize             int64
	ChunkSize           int64
	PadLastChunk        bool
	KeySeparator        []byte
	StartKey            []byte
	EndKey              []byte
	LevelDBOpt          *ldbOpt.Options
	LevelDBIterRange    *ldbUtil.Range
	LevelDBIterReadOpt  *ldbOpt.ReadOptions
	LevelDBIterWriteOpt *ldbOpt.WriteOptions
}

func CreateDefault() *Options {
	return nil
}
