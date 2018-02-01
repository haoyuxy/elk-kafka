package filter

import (
	"flag"
	"fmt"
	goconf "github.com/akrennmair/goconf"
	"os"
	"regexp"
)

type Cfg struct {
	Kafka      string
	Group      string
	Topic      string
	Rule_url   string
	User_media string
}

type Rule struct {
	FilePattern string
	LogPattern  string
	Expired     int64
	Count       int
	Compare     string
	StartTime   int
	EndTime     int
	User        string
	Rulel       string
}


func (r *Rule) SetRule(FilePattern string, LogPattern string, Compare string, User string, Rulel string, Expired int64, Count, StartTime, EndTime int) {
	r.FilePattern = FilePattern
	r.LogPattern = LogPattern
	r.Compare = Compare
	r.User = User
	r.Rulel = Rulel
	r.Expired, r.Count, r.StartTime, r.EndTime = Expired, Count, StartTime, EndTime

}

func (r *Rule) Reg(log Log) bool {
	b, _ := regexp.MatchString(r.LogPattern, log.Message)
	if !b {
		return false
	}
	f, _ := regexp.MatchString(r.FilePattern, log.Source)
	if  f {
		//fmt.Println(log.Message)
		return true
	}else{
		return false
	}
}

func (r *Rule) CheckTime(t int) bool {
	if r.EndTime > r.StartTime {
		if t <= r.EndTime && t >= r.StartTime {
			return true
		}else {
			return false
		}
	}else {
		if t >= r.EndTime && t >= r.StartTime {
			return true
		}else {
			return false
		}
	}
	
}

func (r *Rule) CheckCount(c int) bool {
	switch r.Compare {
		   case ">":
                return c > r.Count
	     
        case "<":
                return c < r.Count
        case ">=":
                return c >= r.Count
        case "<=":
                return c <= r.Count
	}
	return false
}

func Rules() []*Rule {
	var r1, r2, r3 Rule
	Ruleslice := make([]*Rule, 0, 300)
	r1.SetRule("^/log/ruby", "INFO", ">", "hao.yu", "rule1", 6, 5, 10, 20)
	r2.SetRule("^/log/ruby", "ERROR", ">", "hao.yu", "rule2", 60, 5, 10, 20)
	r3.SetRule("^/log/ruby", "16", ">", "hao.yu", "rule3", 10, 2, 10, 20)
	Ruleslice = append(Ruleslice, &r1, &r2, &r3)
	//fmt.Println(Ruleslice)
	//for _, r := range Ruleslice {
	//	fmt.Println(*r)
	//}
		return Ruleslice 
}

func Config() Cfg {
	var cfgFile string
	flag.StringVar(&cfgFile, "c", "al.cfg", "go alarm config")
	flag.Parse()
	//fmt.Println("cfg:", cfgFile)

	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("config does no exists ", err)
		}
	}
	var cfg Cfg
	if err := cfg.readconf(cfgFile); err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		return cfg
	}

}

func (conf *Cfg) readconf(file string) error {
	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		return err
	}

	conf.Kafka, err = c.GetString("default", "kafka")
	if err != nil {
		return err
	}

	conf.Group, err = c.GetString("default", "group")
	if err != nil {
		return err
	}
	conf.Topic, err = c.GetString("default", "topic")
	if err != nil {
		return err
	}
	conf.Rule_url, err = c.GetString("default", "rule_url")
	if err != nil {
		return err
	}
	conf.User_media, err = c.GetString("default", "user_media")
	return err
}


