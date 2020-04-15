package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	showDotFiles := flag.Bool("a", false, "show dot files")
	reverse := flag.Bool("r", false, "list in reverse order")
	flag.Parse()

	dirname := os.Args[len(os.Args)-1]

	// . .. fileを取得したい
	c, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	if *reverse {
		for i := 0; i <= len(c)/2-1; i++ {
			opp := len(c) - i - 1
			c[i], c[opp] = c[opp], c[i]
		}
	}

	for _, file := range c {
		if !*showDotFiles && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fmt.Println(file.Name())
	}
}
