package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Unit struct {
	iMap map[int]int
}

var Data map[int]Unit

func Add(s, k, v int) {
	v1, ok := Data[s]
	if ok {
		_, ok := v1.iMap[k]
		if !ok {
			v1.iMap[k] = v
		} else {
			v1.iMap[k] += v
		}
	} else {
		mp := make(map[int]int)
		mp[k] = v
		Ut := Unit{mp}
		Data[k] = Ut
	}
}

func Remove(s, k int) {
	v1, ok := Data[s]
	if ok {
		_, ok := v1.iMap[k]
		if ok {
			delete(v1.iMap, k)
		}
	}
}

func GetSize(s int) int {
	v1, ok := Data[s]
	if ok {
		ln := len(v1.iMap)
		return ln
	}
	return 0
}

func GetValue(s, k int) int {
	v1, ok := Data[s]
	if ok {
		v2, ok := v1.iMap[k]
		if ok {
			return v2
		} else {
			return 0
		}
	}

	return 0
}

func postServer(c *gin.Context) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	input := string(bodyBytes) // this empty

	inputData := strings.Split(input, " ")

	ints := make([]int, len(inputData))

	for i, s := range inputData {
		ints[i], _ = strconv.Atoi(s)
	}

	switch ints[1] {
	case 1:
		set := ints[2]
		key := ints[3]
		val := ints[4]
		go Add(set, key, val)
	case 2:
		set := ints[2]
		key := ints[3]
		go Remove(set, key)
	case 3:
		set := ints[2]
		fmt.Println("Value is %u", GetSize(set))
	case 4:
		set := ints[2]
		key := ints[3]
		fmt.Println("Value is %u", GetValue(set, key))
	case 5:

	}
}

func main() {
	Data = make(map[int]Unit)

	router := gin.Default()
	router.POST("/server", postServer)

	router.Run("localhost:8080")

}
