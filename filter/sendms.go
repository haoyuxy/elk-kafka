package filter

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
)

type User struct {
	Username   string
	Phone  string
	Wechat string
	Email  string
}

func Users(user_url string) []*User {
	u, _ := url.Parse(user_url)
	q := u.Query()
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

	var r []*User
	json.Unmarshal(result, &r)
	return r
}

func Sendwechat(chat, user, msg string) {
	data := make(url.Values)
	data["content"] = []string{msg}
	data["tos"] = []string{user}
	res, err := http.PostForm(chat, data)
	if err != nil {

		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(result))

}

func SendMail(emailurl, user, msg string) {

	data := make(url.Values)
	data["subject"] = []string{"Elk monitor"}
	data["content"] = []string{msg}
	data["tos"] = []string{user}
	res, err := http.PostForm(emailurl, data)
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(result))

}

func Callback(callurl string) {
	log.Printf(callurl)
}
