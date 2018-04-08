package kfs

import ldb "github.com/syndtr/goleveldb/leveldb/opt"
import "github.com/luckcolors/gokfs/opt"
import "github.com/luckcolors/gokfs/btable"
import "github.com/luckcolors/gokfs/constants"
import "github.com/luckcolors/gokfs/sbucket"
import "github.com/luckcolors/gokfs/blockstream"
import "github.com/luckcolors/gokfs/utils"


type Kfs struct {
	Opt *opt.Options
	BTable btable.BTable
	sBuckets []sbucket.SBucket
}


func Open(dbPath string, ldbOpt *ldb.Options, kfsOpt *opt.Options) *Kfs, error {

	if ldbOpt == nil {
		ldbOpt = &ldb.Options{
			OpenFilesCacheCapacity: 1000,
			Compression:            ldb.NoCompression,
			BlockCacheCapacity:     8 * ldb.MiB,
			WriteBuffer:            4 * ldb.MiB,
			ErrorIfMissing:         false,
			ErrorIfExist:           false,
			BlockSize:              4096,
			BlockRestartInterval:   16,
		}
	}

	if kfsOpt == nil {

		rid, err := utils.CreateReferenceId(nil)
		if err != nil {
			panic(err)
		}

		kfsOpt = &opt.Options{
			Rid:          rid,
			MaxTableSize: constants.S * constants.B,
			SBucketOpt: sbucket.Options{
				MaxSize:   constants.S,
				ChunkSize: constants.C,
			},
			BlockStreamOpt: blockstream.Options{
				ChunkSize:    constants.C,
				PadLastChunk: false,
			},
		}
	}
}
