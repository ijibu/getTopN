//golang内存分配研究
package main

import (
	"fmt"
)

var (
	ipFileMap map[string]int
	globIpMap map[string]int
	ijibu     []byte
	ijibu1    []int
	ijibu2    []Item
)

func main() {
	fmt.Printf("%p\n", ipFileMap)
	fmt.Printf("%p\n", ijibu)
	fmt.Printf("%p\n", ijibu1)
	fmt.Printf("%p\n", ijibu2)
	ijibu = []byte{1, 2, 3}
	fmt.Printf("%p\n", ijibu)
	globIpMap = map[string]int{}
	fmt.Printf("%p\n", globIpMap)
	globIpMap = map[string]int{}
	fmt.Printf("%p\n", globIpMap)
}

type Item struct {
	Key string
	Val int
}

/*
	执行输出：
		0x0
		0x0
		0x0
		0x0
		0x11700088
		0x11706400
	可见在定义变量的时候，golang只是默认给变量一个默认的地址0x0，不会有内存分配，只有在初始化的时候才会有内存分配。
*/
