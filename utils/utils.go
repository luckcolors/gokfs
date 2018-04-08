package utils

import (
	"crypto/rand"
	"encoding/binary"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/dustin/go-humanize"
	"github.com/luckcolors/gokfs/constants"
	"github.com/luckcolors/hashutil"
)

var ErrCreateItemFromIndexInvalidKey = errors.New("CreateItemFromIndex: Invalid key")
var ErrCreateItemFromIndexOutOfBounds = errors.New("CreateItemFromIndex: Index is out of bounds")
var ErrCreateReferenceIdInvalidLength = errors.New("CreateReferenceId: Invalid reference ID length")
var ErrTablePathIsNotAFolder = errors.New("Table path is not a folder")
var ErrTablePathIsInvalid = errors.New("Table path does not contain a valid r.id")
var ErrStatReadDir = errors.New("Error while reading the kfs folder")

//Tests if the string is a valid key
//Param: key - The file key
//Returns: boolean
func IsValidKey(key []byte) bool {
	// keyBytes, err := hex.DecodeString(key)
	// if err != nil {
	// 	return false
	// }

	if uint64(len(key)) == (constants.R / 8) {
		return true
	} else {
		return false //, errors.New("ValidateKey: Error key is invalid")
	}
}

//Hashes the given key
//Param: key - The file key
//Returns: string
func HashKey(key []byte) []byte {
	if IsValidKey(key) {
		return key
	}

	return hashutil.Sum(constants.HASH, key)
}

//Coerces input into a valid file key
//Param: key - The file key
//Returns: string
func CoerceKey(key []byte) []byte {
	if !IsValidKey(key) {
		return HashKey(key)
	}
	return key
}

//Get the key name for a data hash + index
//Param: key - Hash of the data
//Param: index - The index of the file chunk
func CreateItemKeyFromIndex(key []byte, index uint64) ([]byte, error) {
	var fileKey = HashKey(key)
	var indexLength = len(strconv.FormatFloat(math.Floor(float64(constants.S)/float64(constants.C)), 'f', -1, 64))
	var indexString = strconv.FormatUint(index, 10)

	var itemIndex string

	if !(uint64(len(fileKey))*8 == constants.R) {
		return nil, ErrCreateItemFromIndexInvalidKey
	}
	if !(len(indexString) <= indexLength) {
		return nil, ErrCreateItemFromIndexOutOfBounds
	}

	for i := 0; i < indexLength-len(indexString); i++ {
		itemIndex += "0"
	}

	itemIndex += indexString

	itemIndexN, err := strconv.ParseUint(itemIndex, 10, 64)
	if err != nil {
		return nil, err
	}

	fileKey = append(fileKey, []byte(" ")...)
	var b []byte
	binary.LittleEndian.PutUint64(b, itemIndexN)
	fileKey = append(fileKey, b...)

	return fileKey, nil
}

//Get the file name of an s bucket based on it's indexString
//Param: SbucketIndex - The index of the bucket in the B-table
func CreateSbucketNameFromIndex(sBucketIndex uint64) string {
	const indexLength = len(string(constants.B))
	var indexString = strconv.FormatUint(sBucketIndex, 10)

	var leadingZeroes string

	for i := 0; i < indexLength-len(indexString); i++ {
		leadingZeroes = leadingZeroes + "0"
	}
	return leadingZeroes + indexString + ".s"
}

//Creates a random reference ID
//Param: rid - An existing hex reference ID
//Returns: string
func CreateReferenceId(rid []byte) ([]byte, error) {
	if len(rid) == 0 {
		rnd := make([]byte, constants.R/8)
		rand.Read(rnd)
		return rnd, nil
	}

	if len(rid) != 40 {
		return nil, ErrCreateReferenceIdInvalidLength
	}
	return rid, nil
}

//Checks if the given path exists
//Param: filePath
//Returns bool
func FileDoesExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return true
}

func ToHumanReadableSize(n uint64) string {
	return humanize.Bytes(n)
}

//Ensures that the given path has a kfs exstension
//Param: tablePath
//Retuns: bool
func CoerceTablePath(tablePath string) string {
	if path.Ext(tablePath) != ".kfs" {
		return tablePath + ".kfs"
	}
	return tablePath
}

func ValidateSBucketPath(sbucketpath string) bool {
	if path.Ext(sbucketpath) != ".s" {
		return false
	}
	b := path.Base(sbucketpath)
	_, err := strconv.ParseUint(strings.Replace(b, ".s", "", -1), 10, 54)
	if err != nil {
		return false
	}
	return true
}

func ParseSBucketIndexFromPath(sbucketpath string) (o uint64) {
	o, _ = strconv.ParseUint(sbucketpath, 10, 64)
	return
}

func InitBtableDirectory(tablePath string, rid []byte) error {
	err := os.MkdirAll(tablePath, 0600)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tablePath, rid, 0600)
	if err != nil {
		return err
	}
	return nil
}

func ValidateTablePath(tablePath string) error {
	var path = path.Join(tablePath, "r.id")
	d, err := os.Stat(path)

	if err != nil {
		return err
	}
	if !d.IsDir() {
		return errors.Wrap(err, ErrTablePathIsNotAFolder.Error())
	}

	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		return errors.Wrap(err, ErrTablePathIsInvalid.Error())
	}

	return nil
}

func ExistingSbucketIndexes(tablePath string) (*[]uint64, error) {
	folders, err := ioutil.ReadDir(tablePath)
	if err != nil {
		return nil, errors.Wrap(err, ErrStatReadDir.Error())
	}

	var sbucketIndexes = make([]uint64, len(folders)-1) // Remove r.id
	var i = 0
	for _, f := range folders {
		if f.IsDir() && ValidateSBucketPath(f.Name()) {
			sbucketIndexes[i] = ParseSBucketIndexFromPath(f.Name())
			i++
		}
	}
	return &sbucketIndexes, nil
}
