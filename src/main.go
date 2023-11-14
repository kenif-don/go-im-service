package main

import "IM-Service/src/util"

func main() {
	data, e := util.DecryptAes("PJY6TCp6H8ThtRN/vYtW4Er5AEibcqE/lYBIjP2DqhU=", "ccc1307504b8dd0c41a2cdf8743685c9")
	if e != nil {
		println(e)
	}
	println(data)
}
