package filter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
)

type User struct {
	Username string
	Phone    string
	Wechat   string
	Email    string
}

func Users(user_url string) []*User {
	u, _ := url.Parse(user_url)
	q := u.Query()
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var r []*User
	json.Unmarshal(result, &r)
	return r
}

func GrouptoUser(groups []string, user_url string) []string {
	var rs []*User
	for _, g := range groups {
		user_url = user_url + "?groupname=" + g
		u, _ := url.Parse(user_url)
		q := u.Query()
		u.RawQuery = q.Encode()
		res, err := http.Get(u.String())
		if err != nil {
			log.Fatal(err)
		}
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var r []*User
		json.Unmarshal(result, &r)
		rs = append(rs, r...)
	}
	us := make([]string, 0, 20)
	for _, u := range rs {
		us = append(us, u.Username)
	}
	return us
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	sort.Strings(a)
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func Sendwechat(chat, user, msg, zuser, mslogrul string) {
	data := make(url.Values)
	data["content"] = []string{msg}
	data["tos"] = []string{user}
	res, err := http.PostForm(chat, data)
	if err != nil {
		log.Println(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		Sendmslog("0", mslogrul, "wechat", user, msg)
		log.Println(err)
		log.Printf("send wechat to  %s error.\n%s", zuser, msg)
	} else {
		log.Printf("send wechat to  %s success.\n%s", zuser, msg)
		Sendmslog("1", mslogrul, "wechat", zuser, msg)
	}
	log.Printf(string(result))

}

func SendMail(emailurl, user, msg, zuser, mslogrul string) {

	data := make(url.Values)
	data["subject"] = []string{"Elk monitor"}
	data["content"] = []string{msg}
	data["tos"] = []string{user}
	res, err := http.PostForm(emailurl, data)
	if err != nil {
		log.Println(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		Sendmslog("0", mslogrul, "Email", zuser, msg)
		log.Println(err)
		log.Printf("send email to  %s error.\n%s", zuser, msg)
	} else {
		Sendmslog("1", mslogrul, "Email", user, msg)
		log.Printf("send email to  %s success.\n%s", zuser, msg)
	}
	log.Println(string(result))

}

func Sendmslog(status, logurl, channel, user, msg string) {
	data := make(url.Values)
	data["channel"] = []string{channel}
	data["user"] = []string{user}
	data["msg"] = []string{msg}
	data["status"] = []string{status}
	res, err := http.PostForm(logurl, data)
	if err != nil {
		log.Println(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
	}

	log.Println(string(result))
}

func Callback(callurl string) {
	log.Println(callurl)
}
