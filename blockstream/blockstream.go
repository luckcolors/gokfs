package blockstream

import (
	"errors"
	"io"

	"github.com/luckcolors/gokfs/opt"
	"github.com/luckcolors/gokfs/sbucket"
	"github.com/luckcolors/gokfs/utils"
)

var FileNotFoundError = errors.New("BlockStream: File not found")

type BlockStream struct {
	SB      *sbucket.SBucket
	FileKey []byte
	Index   uint64
	Opt     *opt.Options
	pr      io.PipeReader
	pw      io.PipeWriter
}

func (bs BlockStream) New() {
}

// read() is an internal function for handling the write of a singular key, it's called by Write in a loop

func (bs BlockStream) read() ([]byte, error) {
	itemKey, err := utils.CreateItemKeyFromIndex(bs.FileKey, bs.Index)
	if err != nil {
		return nil, err
	}
	d, err := bs.SB.DB.Get(itemKey, nil)
	if err != nil {
		return nil, FileNotFoundError
	}
	bs.Index++
	return d, nil
}

// write() is an internal function for handling the write of a singular key, it's called by Write in a loop

func (bs BlockStream) write(d []byte) error {
	itemKey, err := utils.CreateItemKeyFromIndex(bs.FileKey, bs.Index)
	if err != nil {
		return err
	}
	err = bs.SB.DB.Put(itemKey, d, nil)
	if err != nil {
		return err
	}
	bs.Index++
	return nil
}

func (bs BlockStream) Write(r io.Reader) error {
	var untilEOFError = true
	for untilEOFError {
		buf := make([]byte, int(bs.SB.Opt.ChunkSize))

		n, err := io.ReadFull(r, buf)

		if err == io.EOF {
			untilEOFError = false

			if bs.SB.Opt.PadLastChunk {

				paddingLength := bs.SB.Opt.ChunkSize - n
				for paddingLength < bs.SB.Opt.ChunkSize {
					buf[n-1] = 0
					n++
				}
			}
		}

		err = bs.write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bs BlockStream) Read() (w io.Writer, err error) {
	var untilEOFError = true
	for untilEOFError {

	}
	return
}
