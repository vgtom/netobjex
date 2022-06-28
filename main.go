package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Unit struct {
	iMap map[int]int
}

var Data map[int]Unit

var mtx sync.Mutex

func Add(c *gin.Context, s, k, v int) {
	go func() {
		mtx.Lock()
		defer mtx.Unlock()

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
	}()

	c.String(http.StatusOK, "0")
}

func Remove(c *gin.Context, s, k int) {
	go func() {
		mtx.Lock()
		defer mtx.Unlock()

		v1, ok := Data[s]
		if ok {
			_, ok := v1.iMap[k]
			if ok {
				delete(v1.iMap, k)
			}
		}
	}()

	c.String(http.StatusOK, "0")
}

func GetSize(c *gin.Context, s int) {
	mtx.Lock()
	defer mtx.Unlock()

	v1, ok := Data[s]
	if ok {
		ln := len(v1.iMap)
		str := fmt.Sprintf("1 %d", ln)
		c.String(http.StatusOK, str)
		return
	}
	str := fmt.Sprintf("1 %d", 0)
	c.String(http.StatusOK, str)
}

func GetValue(c *gin.Context, s, k int) {
	mtx.Lock()
	defer mtx.Unlock()

	v1, ok := Data[s]
	if ok {
		v2, ok := v1.iMap[k]
		if ok {
			str := fmt.Sprintf("1 %d", v2)
			c.String(http.StatusOK, str)
			return
		}
	}

	str := fmt.Sprintf("1 %d", 0)
	c.String(http.StatusOK, str)
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
		Add(c, set, key, val)
	case 2:
		set := ints[2]
		key := ints[3]
		Remove(c, set, key)
	case 3:
		set := ints[2]
		GetSize(c, set)
	case 4:
		set := ints[2]
		key := ints[3]
		GetValue(c, set, key)
	case 6:
		c.Abort()
	}
}

func main() {
	Data = make(map[int]Unit)

	router := gin.Default()
	router.POST("/server", postServer)

	router.Run("localhost:8080")

}
