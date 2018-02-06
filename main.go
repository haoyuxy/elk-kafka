package main

import (
	"github.com/haoyuxy/elk-kafka/filter"
)

const (
	MaxCount = 2000  //max xxx count
)

func main() {
	
	cfg := filter.Config()
	filter.KafkaOut(MaxCount, &cfg)
	
}
