package filter

import (
	"fmt"
	"log"
	"os"
	"os/signal"
//	"strconv"
	"time"
	"strings"

	cluster "github.com/bsm/sarama-cluster"
)

//func KafkaOut(MaxCount int, topic ,group ,ip_port, apiurl string) {
func KafkaOut(MaxCount int, cfg *Cfg) {
/*
	s1 := make([]int64, 0, 1000)	
	expiredtime := make(map[string][]int64)
	lastalarmtime := make(map[string]int64)
	rulemap := make(map[string]*Rule)
*/	
	ApiUrl := cfg.Apiurl
	expiredtime := make([][]int64, 0, 2000)
	lastalarmtime := make([]int64, 0, 1)
	rulemap  := make([]*Rule,0,2000)
	//Ruleslice := Rules()
	Ruleslice := Rules(ApiUrl + "elk/")
	fmt.Println(Ruleslice)
	for _, r := range Ruleslice {
		/*
		k := "als" + strconv.Itoa(i)
		expiredtime[k] = s1
		lastalarmtime[k] = 0
		rulemap[k] = r
		*/
		expiredtime = append(expiredtime,make([]int64, 0, 2000))
		lastalarmtime = append(lastalarmtime,0)
		rulemap = append(rulemap,r)

	}

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init consumer
	brokers := []string{cfg.Kafka}
	topics := []string{cfg.Topic}
	consumer, err := cluster.NewConsumer(brokers, cfg.Group, topics, config)
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
				//fmt.Println(string(msg.Value))
				log := JsontoStr(msg.Value)
				//fmt.Println(log)
				go func(log *Log) {
					for k, v := range rulemap {
						if v.Reg(log) {
							//fmt.Println(log)
							//fmt.Println(v.LogPattern,v.FilePattern)
							//	go func() {
							if v.CheckTime(time.Now().Hour()) {
								//fmt.Println(k)
								nowtime := time.Now().Unix()
								expiredtime[k] = append(expiredtime[k], nowtime+v.Expired)
								expiredtime[k] = Expire(expiredtime[k])
								ncount := len(expiredtime[k])
								if v.CheckCount(ncount) && v.CheckLastTime(lastalarmtime[k], nowtime) {
									lastalarmtime[k] = nowtime

									fmt.Println("alarm", len(expiredtime[k]), k, expiredtime[k], v.FilePattern)
									//fmt.Println(len(expiredtime["als0"]), len(expiredtime["als1"]))
									users := v.User
									us := strings.Split(users,",")
									//fmt.Println(us)
									userslice := Users(ApiUrl + "users")
									//fmt.Println(us,userslice)
									for _, u := range us {
										for _, u2 := range userslice {
											//fmt.Println(u2.Username)
											if u2.Username == u {
												//fmt.Println(u)
												emsg := v.Msg + "\n" + v.LogPattern + "\n" + log.Source
												SendMail(cfg.Mailurl, u2.Email, emsg)
												Sendwechat(cfg.Wechaturl, u2.Wechat, emsg)
												Callback(v.Callback)
											}

										}

									}

								}
								if ncount > MaxCount {
									expiredtime[k] = expiredtime[k][:0]
								}
							}
							//	}()
						}
					}
				}(&log)

				//	consumer.MarkOffset(msg, "")	// mark message as processed
				consumer.MarkOffset(msg, "") // mark message as processed
			//	time.Sleep(1e9)
				// fmt.Println("------------------------------")
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
