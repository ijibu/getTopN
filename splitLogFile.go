//分而治之
//该文件用来把ip日志文件分割成1000个小文件，方便后续处理。本来该用hash(IP)/1000，但是还没学会
//golang的hash用法。
package main

import (
	"bufio"
	//"fmt"
	"os"
	"strconv"
	"strings"
)

var fileCacheMap map[int]*os.File

func main() {
	//var ip string = "127.0.0.1"
	//ret := ip2long(ip)
	//fmt.Println(ret)
	//fmt.Println(ret / 1000)
	//fmt.Println(ret % 1000)
	fileCacheMap = map[int]*os.File{}
	splitFile("ip.log")
	destruct()
	return
}

/**
 * hash分割文件。
 */
func splitFile(fileName string) {
	logFile, _ := os.Open(fileName)
	defer logFile.Close()
	reader := bufio.NewReader(logFile)

	// 逐行读入分词
	for {
		sip, _ := reader.ReadString('\n')

		if len(sip) == 0 {
			// 文件结束
			break
		}

		ip := ip2long(sip)

		fileIndex := ip % 1000
		index := int(fileIndex)
		insertToFile(index, sip)
	}

	return
}

func insertToFile(fileIndex int, sip string) {
	//把文件句柄缓存起来，会减少文件的频繁打开和关闭，会大大提高IO效率。初步测试，IO由850kb提高到6.3M。
	f, ok := fileCacheMap[fileIndex]
	if !ok {
		fileName := "./log/" + strconv.Itoa(fileIndex) + ".log"
		f, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
		fileCacheMap[fileIndex] = f
	}

	//defer f.Close()
	buf := []byte(sip)
	f.Write(buf)
	return
}

//没有验证IP有效性。
func ip2long(ip1 string) int64 {
	ip := strings.Split(ip1, ".")
	r0, _ := strconv.Atoi(ip[0])
	r1, _ := strconv.Atoi(ip[1])
	r2, _ := strconv.Atoi(ip[2])
	r3, _ := strconv.Atoi(ip[3])

	h0 := int64(r0)
	h1 := int64(r1)
	h2 := int64(r2)
	h3 := int64(r3)

	return h0*16777216 + h1*65536 + h2*256 + h3*1
}

//析构函数，关闭所有打开的文件句柄。
func destruct() {
	for _, f := range fileCacheMap {
		f.Close()
	}
}
