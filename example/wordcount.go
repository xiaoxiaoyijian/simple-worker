package main

import (
	"fmt"
	"github.com/xiaoxiaoyijian/simple-mapreduce/utils/file"
	"github.com/xiaoxiaoyijian/simple-worker/core"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(int64(time.Second))

	ret := make(map[string]int)
	worker := core.NewWorker(MyProcesser, MyErrorHandler, 10, 10)
	out_chan, _ := worker.Run(file.ReadLines("wordcount.go"))
	for mapResult := range out_chan {
		m, ok := mapResult.(map[string]int)
		if !ok {
			continue
		}

		for k, v := range m {
			_, ok := ret[k]
			if ok {
				ret[k] += v
			} else {
				ret[k] = v
			}
		}
	}

	for k, v := range ret {
		println(fmt.Sprintf("%v          %v", k, v))
	}

}

func MyProcesser(input interface{}) (output interface{}, err error) {
	ret := make(map[string]int)

	line := strings.TrimSpace(input.(string))
	tokens := strings.Split(line, " ")
	for _, value := range tokens {
		if value != "" {
			v, ok := ret[value]
			if ok {
				ret[value] = v + 1
			} else {
				ret[value] = 1
			}
		}
	}

	return ret, nil
}

func MyErrorHandler(input interface{}, err error) (err_output interface{}) {
	return strings.TrimSpace(input.(string)) + ":" + err.Error()
}
