//分别对每个小文件进行hash_map统计，找到前10个。然后把1000个小文件中的前10个合并起来，再进行排序，取得前10个。
//目前情况下，内存占用517M,CPU占用4.16%,IO约5M。内存占用怎么会这么高啊？
//参考:http://blog.jobbole.com/48969/,采用内存复用，应该会解决问题。
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const chanLen = 4

var (
	c         chan int
	ipFileMap map[string]int
	globIpMap map[string]int
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置cpu的核的数量，从而实现高并发
	c = make(chan int, chanLen)          //简单采用goruntime缓冲，同时最多4个执行，和CPU的数量一致。
	path := "./log"
	globIpMap = map[string]int{}

	filepath.Walk(path, func(path string, f os.FileInfo, e error) error {
		if f == nil {
			return e
		}
		if f.IsDir() {
			return nil
		}

		parseIpLogs(path)

		return nil
	})

	ms := NewMapSorter(globIpMap)
	sort.Sort(ms)
	for i, item := range ms {
		if i >= 10 {
			break
		}
		fmt.Printf("%s:%d\n", item.Key, item.Val)
	}
	fmt.Println("scuess!")
}

//遍历每一个IP文件，生成一个map，统计每个IP出现的次数。然后进行排序。
func parseIpLogs(path string) {
	ipFileMap = map[string]int{}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	ss, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	sz := string(ss)
	content := strings.Split(sz, "\r\n")
	sl := len(content)
	//index := getFileName(path)
	for i := 0; i < sl; i++ {
		ip, ok := ipFileMap[content[i]]
		if !ok {
			ipFileMap[content[i]] = 1
		} else {
			ipFileMap[content[i]] = ip + 1
		}
	}

	//map排序
	ms := NewMapSorter(ipFileMap)
	sort.Sort(ms)
	for i, item := range ms {
		if i >= 10 {
			break
		}
		//获取前10个保存起来。
		globIpMap[item.Key] = item.Val
	}
}

/**
 * 根据路径名获取文件名
 */
func getFileName(fileFullPath string) string {
	fName := filepath.Base(fileFullPath)
	extName := filepath.Ext(fileFullPath)
	bName := fName[:len(fName)-len(extName)]
	return bName
}

//map排序相关的，参考网站：http://www.dotcoo.com/golang-sort，https://gist.github.com/ikbear/4038654
type MapSorter []Item

type Item struct {
	Key string
	Val int
}

func NewMapSorter(m map[string]int) MapSorter {
	ms := make(MapSorter, 0, len(m))

	for k, v := range m {
		ms = append(ms, Item{k, v})
	}

	return ms
}

func (ms MapSorter) Len() int {
	return len(ms)
}

func (ms MapSorter) Less(i, j int) bool {
	return ms[i].Val > ms[j].Val // 按值降序排序
	//return ms[i].Key < ms[j].Key // 按键升序排序
}

func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func sortIp(m map[string]int) {
	ms := NewMapSorter(m)
	sort.Sort(ms)

	for i, item := range ms {
		if i >= 10 {
			break
		}
		fmt.Printf("%s:%d\n", item.Key, item.Val)
	}
}
