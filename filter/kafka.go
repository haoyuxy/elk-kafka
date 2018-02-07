package filter

import (
	"bytes"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"github.com/influxdata/influxdb/client/v2"
)

func KafkaOut(MaxCount int, cfg *Cfg) {

	ApiUrl := cfg.Apiurl                    //获取规则和用户的url
	expiredtime := make([][]int64, 0, 2000) //存储匹配到的指标队列
	lastalarmtime := make([]int64, 0, 1)    //每一条规则的上次告警时间
	//ruleslice := make([]*Rule, 0, 2000)       //规则队列
	lastmap := make(map[string]int64)
	//Ruleslice := Rules()
	Ruleslice := Rules(ApiUrl + "elk/")
	fmt.Println(Ruleslice)
	//for _, _ := range Ruleslice { //初始化各种队列
	for i:=0 ; i < len(Ruleslice); i++ {
		expiredtime = append(expiredtime, make([]int64, 0, 2000))
		lastalarmtime = append(lastalarmtime, 0)
		//ruleslice = append(ruleslice, r)

	}
	
	influxconfig := client.UDPConfig{Addr: cfg.Influxdburl}
	ic, err := client.NewUDPClient(influxconfig)
	if err != nil {
		panic(err.Error())
	}
	
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
	})

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
				//fmt.Println(string(msg.Value))
				log := JsontoStr(msg.Value)
				go func(log *Log) { //用每条规则检查日志
					fmt.Println(log.Source)
					for k, v := range Ruleslice {
						if v.Reg(log) { //关键字检查
							tags := map[string]string{"rule": v.Rulename,"logfile":log.Source,"ip":log.Beat.Name}
							fields := map[string]interface{}{"value":1,}
							pt, _ := client.NewPoint("elkmonitor", tags, fields, time.Now())
							bp.AddPoint(pt)
							ic.Write(bp)
							if v.CheckTime(time.Now().Hour()) { //检查是否在告警时间段
								nowtime := time.Now().Unix()
								expiredtime[k] = append(expiredtime[k], nowtime+v.Expired)
								expiredtime[k] = Expire(expiredtime[k])                                 //清除过期数据
								ncount := len(expiredtime[k])                                           //当前队列长度
								if v.CheckCount(ncount) && v.CheckLastTime(lastalarmtime[k], nowtime) { //超过阈值并且上次告警时间超过定义时间发送告警
									lastkey := log.Beat.Name + v.LogPattern + log.Source
									if CheckLastmsg(lastmap, lastkey, nowtime, v.Nextalarmtime) {
										lastalarmtime[k] = nowtime
										lastkey := log.Beat.Name + v.LogPattern + log.Source
										lastmap[lastkey] = nowtime
										users := v.User //获取报警用户
										us := strings.Split(users, ",")
										userslice := Users(ApiUrl + "users")
										for _, u := range us {
											for _, u2 := range userslice {
												if u2.Username == u {
													currentTime := time.Now().Format("2006-01-02 15:04:05")
													var emsg bytes.Buffer
													emsg.WriteString(currentTime)
													emsg.WriteString("\n")
													emsg.WriteString("规则名称: ")
													emsg.WriteString(v.Rulename)
													emsg.WriteString("\n")
													emsg.WriteString("告警信息: ")
													emsg.WriteString(v.Msg)
													emsg.WriteString("\n")
													emsg.WriteString("ip: ")
													emsg.WriteString(log.Beat.Name)
													emsg.WriteString("\n")
													emsg.WriteString("匹配规则: ")
													emsg.WriteString(v.LogPattern)
													emsg.WriteString("\n")
													emsg.WriteString("日志路径: ")
													emsg.WriteString(log.Source)
													//emsg := v.Msg + "\n" + v.LogPattern + "\n" + log.Source + "\n" + log.Beat.Name //告警信息
													semsg := emsg.String()
													if u2.Email != "" && cfg.Mailurl != "" {
														SendMail(cfg.Mailurl, u2.Email, semsg) //发送邮件
													}
													if u2.Wechat != "" && cfg.Wechaturl != "" {
														Sendwechat(cfg.Wechaturl, u2.Wechat, semsg) //发送微信
													}
													Callback(v.Callback) //回调函数
												}

											}

										}
									}

								}
								if ncount > MaxCount { //超过队列最大存储长度清空队列
									expiredtime[k] = expiredtime[k][:0]
								}
							}
						}
					}
				}(&log)

				consumer.MarkOffset(msg, "") // mark message as processed
				time.Sleep(1e9)              //测试时使用，上线关掉
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
