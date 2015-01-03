/*
Structs and functions for work with http://www.diary.ru
Skills: JSON, http, Cookies, Auth
Written by Yuri (Yuriy) Astrov
P.S.
JSON library is true, server is bad, because integers must be without quotes!
*/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/disintegration/charmap"
	"io" //For one version JSON parsing and MyRequest
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
)

type Unit struct {
	Count int `json:"count,string"`
	/*This must we doing because int represented as string*/
	Obj *json.RawMessage
}

type UmailZeroInfo struct {
	Sender string `json:"from_username"`
	Title  string `json:"title"`
}

type UmailInfo struct {
	Count int            `json:"count,string"`
	First *UmailZeroInfo `json:"0"`
}

type UserInfo struct {
	Userid    string
	Username  string
	Shortname string
}

type DiaryInfoRuJson struct {
	Newcomments Unit
	Discuss     Unit
	Umails      UmailInfo `json:"umails"`
	Userinfo    UserInfo
	Error       string `json:"error"`
}

func (o *DiaryInfoRuJson) Print() {
	if o.Error == "" {
		fmt.Printf("Shortname: %s\n", o.Userinfo.Shortname)
		fmt.Printf("Username: %s\n", o.Userinfo.Username)
		fmt.Printf("Comments count: %d\n", o.Newcomments.Count)
		fmt.Printf("Discuss count: %d\n", o.Discuss.Count)
		fmt.Printf("Umails count: %d\n", o.Umails.Count)
		if o.Umails.First != nil {
			fmt.Println(&(o.Umails.First).Sender, ": ", o.Umails.First.Title)
		}
	} else {
		fmt.Printf("Error: %s\n", o.Error)
	}
}

const (
	DiaryMainUrl        = "http://www.diary.ru/"
	DiaryJSONRequestUrl = "http://pay.diary.ru/yandex/online.php"
	UserAgent           = "goDiaryRyClient"
)

type DiaryRuClient struct {
	HttpClient  *http.Client
	LastMessage DiaryInfoRuJson
	Re          *regexp.Regexp
}

func MyRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, urlStr, body)
	r.Header.Add("User-agent", UserAgent)
	r.Header.Add("Referer", DiaryMainUrl)
	return r, err
}

func (o *DiaryRuClient) Init() {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	o.HttpClient = &http.Client{
		Jar: cookieJar,
	}
	o.Re = regexp.MustCompile("name=\"signature\" value=\"(.+?)\"")
}

func (o *DiaryRuClient) Auth(user, password string) {
	strcp1251, _ := charmap.Encode(user, "cp-1251")
	//re := regexp.MustCompile("name=\"signature\" value=\"(.+?)\"")
	// Prepare values
	data := url.Values{}
	data.Add("user_login", strcp1251)
	data.Add("user_pass", password)
	// Get signature
	resp, _ := o.HttpClient.Get(DiaryMainUrl)
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	s := o.Re.FindStringSubmatch(string(body))
	data.Add("signature", s[1])
	// Prepare Url
	apiUrl := "http://pda.diary.ru"
	resource := "/login.php"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)
	// Send POST with auth data
	r, _ := MyRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, _ = o.HttpClient.Do(r)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
}

func (o *DiaryRuClient) RestJsonInfo() DiaryInfoRuJson {
	r, _ := MyRequest("GET", DiaryJSONRequestUrl, nil)
	resp, _ := o.HttpClient.Do(r)

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	// First JSON parsing version (For array in root)
	/*dec := json.NewDecoder(resp.Body)
	    var m DiaryInfoRuJson
		for {
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
				fmt.Println(err)
			}
			fmt.Println(m)
		}*/
	// Second version
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &o.LastMessage)
	return o.LastMessage
}

func main() {
	userPtr := flag.String("user", "", "Username (Login)")
	passPtr := flag.String("pass", "", "Password")
	flag.Parse()
	var diary DiaryRuClient
	diary.Init()
	diary.Auth(*userPtr, *passPtr)
	m := diary.RestJsonInfo()
	m.Print()
}
