package main

import (
	"encoding/hex"
	"fmt"

	"github.com/luckcolors/gokfs/utils"
	"github.com/luckcolors/hashutil"
)

func main() {
	x, _ := hex.DecodeString("ddadef707ba62c166051b9e3cd0294c27515f2bc")
	fmt.Println(1, utils.IsValidKey([]byte("A")))  //1
	fmt.Printf("%x", utils.CoerceKey([]byte("A"))) //2
	fmt.Println()
	fmt.Println(3, utils.IsValidKey(x))          //3
	fmt.Printf("%x", utils.HashKey([]byte("A"))) //4
	fmt.Println()
	fmt.Println(utils.CreateItemKeyFromIndex([]byte("A"), 0)) //5
	fmt.Println(utils.CreateItemKeyFromIndex(x, 2213))        //6
	fmt.Println(7, utils.CreateSbucketNameFromIndex(2213))    //7
	fmt.Println(utils.CreateReferenceId([]byte("")))
	fmt.Println(utils.FileDoesExist("./main.go"))
	fmt.Println(utils.ToHumanReadableSize(100000))
	fmt.Println(utils.CoerceTablePath("a"))
	fmt.Println(utils.CoerceTablePath("a.kfs"))

	fmt.Println("=========")
	fmt.Printf("%x", hashutil.Sum("ripemd160", nil))
	fmt.Println()
}
