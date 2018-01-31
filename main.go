package main

import (
	//"fmt"
	"github.com/haoyuxy/gotest/filter"
)
func main() {
	
	c := filter.Config()
	filter.KafkaOut(c.Topic, c.Group, c.Kafka)
	
}