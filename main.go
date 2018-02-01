package main

import (
	//"fmt"
	"github.com/haoyuxy/elk-kafka/filter"
)

const (
	MaxCount = 1000  //max xxx count
)
func main() {
	
	c := filter.Config()
	filter.KafkaOut(MaxCount, c.Topic, c.Group, c.Kafka)
	
}
