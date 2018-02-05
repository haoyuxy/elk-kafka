package filter

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	goconf "github.com/akrennmair/goconf"
)

type Cfg struct {
	Kafka     string
	Group     string
	Topic     string
	Apiurl    string
	Mailurl   string
	Wechaturl string
	Phone     string
}

type Rule struct {
	FilePattern   string
	LogPattern    string
	Expired       int64
	Count         int
	Compare       string
	StartTime     int
	EndTime       int
	User          string
	Rulel         string
	Callback      string
	Nextalarmtime int64
	Msg           string
}

/*
func (r *Rule) SetRule(FilePattern string, LogPattern string, Compare string, User string, Rulel string, Callback string, Expired int64,Nextalarmtime int64, Count, StartTime, EndTime int) {
	r.FilePattern = FilePattern
	r.LogPattern = LogPattern
	r.Compare = Compare
	r.User = User
	r.Rulel = Rulel
	r.Callback = Callback
	r.Expired, r.Count, r.StartTime, r.EndTime, r.Nextalarmtime = Expired, Count, StartTime, EndTime, Nextalarmtime

}
*/

func (r *Rule) Reg(log *Log) bool {
	b, _ := regexp.MatchString(r.FilePattern, log.Source)
	//fmt.Println(log.Source,r.FilePattern,b)
	if !b {
		return false
	}
	f, _ := regexp.MatchString(r.LogPattern, log.Message)
	if f {
		//fmt.Println(log.Message)
		return true
	} else {
		return false
	}
}

func (r *Rule) CheckTime(t int) bool {
	if r.EndTime > r.StartTime {
		if t <= r.EndTime && t >= r.StartTime {
			return true
		} else {
			return false
		}
	} else {
		if t >= r.EndTime && t >= r.StartTime {
			return true
		} else {
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

func (r *Rule) CheckLastTime(last, now int64) bool {
	if last == 0 || r.Nextalarmtime == 0 {
		return true
	} else {
		if last+r.Nextalarmtime*60 < now {
			return true
		} else {
			return false
		}
	}
}

/*
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
*/

func Rules(rule_url string) []*Rule {
	u, _ := url.Parse(rule_url)
	q := u.Query()
	//q.Set("username", "user")
	//q.Set("password", "passwd")
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	var r []*Rule
	json.Unmarshal(result, &r)
	return r
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
	
	conf.Mailurl, err = c.GetString("default", "mail")
	if err != nil {
		return err
	}
	
	conf.Wechaturl, err = c.GetString("default", "wechat")
	if err != nil {
		return err
	}
	

	conf.Apiurl, err = c.GetString("default", "apiurl")
	return err
}
