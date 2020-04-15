package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	showDotFiles := flag.Bool("a", false, "show dot files")
	longStat := flag.Bool("l", false, "show long file status")
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
		if *longStat {
			permission := file.Mode().String()

			nlink := uint64(0)
			ownerName := ""
			groupName := ""
			if sys := file.Sys(); sys != nil {
				if stat, ok := sys.(*syscall.Stat_t); ok {
					nlink = uint64(stat.Nlink)

					u, err := user.LookupId(fmt.Sprint(stat.Uid))
					if err != nil {
						panic(err)
					}
					ownerName = u.Username

					g, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
					if err != nil {
						panic(err)
					}
					groupName = g.Name
				}
			}

			byteSize := file.Size()
			timeStamp := file.ModTime().Format("Jan 2 15:04")

			fmt.Println(permission + " " + strconv.Itoa(int(nlink)) + " " + ownerName + " " + groupName + " " + strconv.Itoa(int(byteSize)) + " " + timeStamp + " " + file.Name())

		} else {
			fmt.Println(file.Name())

		}
	}
}
