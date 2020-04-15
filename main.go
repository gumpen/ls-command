package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"syscall"
)

type Options struct {
	all     bool
	long    bool
	reverse bool
}

type Infos struct {
	permissions string
	hardLinkNum string
	owner       string
	group       string
	byteSize    string
	timeStamp   string
	name        string
}

type Width struct {
	permissions int
	hardLinkNum int
	owner       int
	group       int
	byteSize    int
	timeStamp   int
	name        int
}

func main() {
	options := Options{}
	flag.BoolVar(&options.all, "a", false, "show dot files")
	flag.BoolVar(&options.long, "l", false, "show long file status")
	flag.BoolVar(&options.reverse, "r", false, "list in reverse order")
	flag.Parse()

	dirname := os.Args[len(os.Args)-1]

	// . .. fileを取得したい
	c, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	if options.reverse {
		for i := 0; i <= len(c)/2-1; i++ {
			opp := len(c) - i - 1
			c[i], c[opp] = c[opp], c[i]
		}
	}

	width := Width{0, 0, 0, 0, 0, 0, 0}

	infoList := []Infos{}
	for _, file := range c {
		if !options.all && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		info := Infos{}
		if options.long {
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

			info.permissions = file.Mode().String()
			info.hardLinkNum = strconv.Itoa(int(nlink))
			info.owner = ownerName
			info.group = groupName
			info.byteSize = strconv.Itoa(int(file.Size()))
			info.timeStamp = file.ModTime().Format("Jan 2 15:04")

			width.permissions = max(width.permissions, len(info.permissions))
			width.hardLinkNum = max(width.hardLinkNum, len(info.hardLinkNum))
			width.owner = max(width.owner, len(info.owner))
			width.group = max(width.group, len(info.group))
			width.byteSize = max(width.byteSize, len(info.byteSize))
			width.timeStamp = max(width.timeStamp, len(info.timeStamp))

		}
		fileName := file.Name()
		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			realPath, err := os.Readlink(path.Join(dirname, file.Name()))
			if err != nil {
				panic(err)
			}
			fileName = fileName + " -> " + realPath
		}

		if file.IsDir() {
			fileName = fmt.Sprintf("\x1b[34m%s\x1b[0m", fileName)
			if err != nil {
				panic(err)
			}
		}
		info.name = fileName

		width.name = max(width.name, len(info.name))

		infoList = append(infoList, info)
	}
	for _, info := range infoList {
		if options.long {
			permissions := info.permissions
			hardLinkNum := strings.Repeat(" ", (width.hardLinkNum-len(info.hardLinkNum))+2) + info.hardLinkNum
			owner := strings.Repeat(" ", (width.owner-len(info.owner))+2) + info.owner
			group := strings.Repeat(" ", (width.group-len(info.group))+2) + info.group
			byteSize := strings.Repeat(" ", (width.byteSize-len(info.byteSize))+2) + info.byteSize
			timeStamp := strings.Repeat(" ", (width.timeStamp-len(info.timeStamp))+2) + info.timeStamp
			name := strings.Repeat(" ", 2) + info.name

			fmt.Println(permissions + hardLinkNum + owner + group + byteSize + timeStamp + name)

		} else {
			fmt.Println(info.name)
		}
	}

}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
