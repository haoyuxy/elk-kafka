package main

import (
	//"fmt"
	"github.com/haoyuxy/elk-kafka/filter"
)
func main() {
	
	c := filter.Config()
	filter.KafkaOut(c.Topic, c.Group, c.Kafka)
	
}
