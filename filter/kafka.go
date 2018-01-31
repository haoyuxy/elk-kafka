package filter

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	cluster "github.com/bsm/sarama-cluster"
)

func KafkaOut(topic ,group ,ip_port string) {
	
	s1 := make([]int64, 0, 1000)
	data := make(map[string][]int64)
	rulemap := make(map[string]*Rule)
	
	Ruleslice := Rules()
	for i, r := range Ruleslice {
		//*r.Rulel := make([]int64, 0, 1000)
		k := "als" + strconv.Itoa(i)
		data[k] = s1
		rulemap[k] = r
	
	}
	
	
	
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init consumer
	brokers := []string{ip_port}
	topics := []string{topic}
	consumer, err := cluster.NewConsumer(brokers, group, topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
			//	fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			//	fmt.Println(msg.Value)
				log := JsontoStr(msg.Value)
				for k,v := range rulemap {
					if v.Reg(log) {
						if v.CheckTime(time.Now().Hour()) {
							data[k] = append(data[k],time.Now().Unix() + v.Expired)
							data[k] = Expire(data[k])
							if v.CheckCount(len(data[k])) {
								fmt.Println("alarm",len(data[k]),k,data[k])
							} 
						}
					}
				}
		
			//	consumer.MarkOffset(msg, "")	// mark message as processed
				consumer.MarkOffset(msg, "")	// mark message as processed
				time.Sleep(1e9)
                                fmt.Println("------------------------------")
			}
		case <-signals:
			return
		}
	}

}

func Expire(se []int64) []int64 {
	t := time.Now().Unix()
LOOP:
	for i, item := range se {
		if item < t {
			se = append(se[:i], se[i+1:]...)
			goto LOOP
		}
	}
	return se
}
