/**
 * @Author: YMBoom
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2022/05/31 16:30
 */
package main

import (
	"github.com/ymboom0042/wskey2Cookie/src"
	"os"
)

func main() {
	src.ReadConf()
	b, err := os.ReadFile("./wskey.txt")
	if err != nil {
		panic("err : not found wskey.txt")
	}

	wskey := string(b)
	if wskey == "" {
		panic("err := wskey is null")
	}

	w2c, err := src.NewWskey2Cookie(wskey)
	if err != nil {
		panic("err : not found wskey.txt")
		return
	}

	w2c.Do()
}
