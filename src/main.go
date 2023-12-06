package main

import (
	"IM-Service/src/util"
	"fmt"
)

func main() {
	bs, e := util.EncryptAes2([]byte{56, 57, 97, 100, 0, 100, 0, 247, 0, 0, 163, 197, 170, 162, 196, 169}, "7e33184588b1e604e4e6665ce2ebeb2c")
	fmt.Println(bs, e)
	bs, e = util.DecryptAes2(bs, "7e33184588b1e604e4e6665ce2ebeb2c")
	fmt.Println(bs, e)

}
