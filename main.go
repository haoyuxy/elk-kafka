package main

import (
	"runtime"
	"github.com/haoyuxy/elk-kafka/filter"
)

const (
	MaxCount = 2000  //max xxx count
)
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	
	cfg := filter.Config()
	filter.KafkaOut(MaxCount, &cfg)
	
}
