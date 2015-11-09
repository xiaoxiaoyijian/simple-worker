package core

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewPool(t *testing.T) {
	pool := NewPool(10)

	len1, len2 := pool.Len()
	fmt.Printf("Total cnt: %d, timer cnt:%d \n", len1, len2)
}

func Test_Add(t *testing.T) {
	pool := NewPool(10)
	pool.Add("1", "xxxxx", 0)
	pool.Add("2", "xxxxx", 0)
	pool.Add("3", "xxxxx", time.Second*3)
	len1, len2 := pool.Len()
	fmt.Printf("Total cnt: %d, timer cnt:%d \n", len1, len2)
	time.Sleep(time.Second * 5)
	len1, len2 = pool.Len()
	fmt.Printf("Total cnt2: %d, timer cnt2:%d \n", len1, len2)
}

func Test_Remove(t *testing.T) {
	pool := NewPool(10)
	pool.Add("1", "xxxxx", 0)
	pool.Add("2", "xxxxx", 0)
	pool.Add("3", "xxxxx", time.Second*3)
	len1, len2 := pool.Len()
	fmt.Printf("Total cnt: %d, timer cnt:%d \n", len1, len2)
	pool.Remove("1")
	len1, len2 = pool.Len()
	fmt.Printf("Total cnt2: %d, timer cnt2:%d \n", len1, len2)
}

func Test_Remove2(t *testing.T) {
	pool := NewPool(10)
	pool.Add("1", "xxxxx", 0)
	pool.Add("2", "xxxxx", 0)
	pool.Add("3", "xxxxx", time.Second*3)
	len1, len2 := pool.Len()
	fmt.Printf("Total cnt: %d, timer cnt:%d \n", len1, len2)
	pool.Remove("3")
	len1, len2 = pool.Len()
	fmt.Printf("Total cnt2: %d, timer cnt2:%d \n", len1, len2)
}

func Test_Get(t *testing.T) {
	pool := NewPool(10)
	pool.Add("1", "xxxxx", 0)
	pool.Add("2", "xxxxx", 0)
	pool.Add("3", "xxxxx", time.Second*3)
	len1, len2 := pool.Len()
	fmt.Printf("Total cnt: %d, timer cnt:%d \n", len1, len2)
	v, b, err := pool.Get("3")
	fmt.Printf("v: %v, b:%v, err:%v \n", v, b, err)
	v, b, err = pool.Get("2")
	fmt.Printf("v: %v, b:%v, err:%v \n", v, b, err)
}
