/*
Structs and functions for work with http://www.diary.ru API
Skills: JSON, http, Auth, MD5Sum and hex.EncodeToString, url.URL, Unknown JSON
response decoding
Written by Yuri (Yuriy) Astrov
P.S.
JSON library is true, server is bad, because integers must be without quotes!
*/
package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/disintegration/charmap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	appkey       = "" // ok also
	key          = "" // pk also
	DiaryMainUrl = "http://www.diary.ru/api/"
	UserAgent    = "goDiaryRyClient"
)

type DiaryAPIAuthResponse struct {
	Result int    `json:"result,string"`
	SID    string `json:"sid"`
	Error  string `json:"error"`
}

type DiaryAPIClient struct {
	HttpClient *http.Client
	SID        string
	Timestamp  time.Time
	URL        *url.URL
}

func (this *DiaryAPIClient) Init() {
	this.HttpClient = &http.Client{}
}

func (this *DiaryAPIClient) Auth(user, password string) {
	strcp1251, _ := charmap.Encode(user, "cp-1251")

	values := url.Values{}
	values.Add("username", strcp1251)
	hash := md5.Sum([]byte(key + password))
	values.Add("password", hex.EncodeToString(hash[:]))
	values.Add("method", "user.auth")
	values.Add("appkey", appkey)

	this.URL, _ = url.Parse(DiaryMainUrl)
	r, _ := this.dorequest(values, nil)
	resp, _ := this.HttpClient.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	var message DiaryAPIAuthResponse
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&message)
	if err != nil {
		log.Fatal(err)
	}
	this.SID = message.SID
	this.Timestamp = time.Now()
}

func (this *DiaryAPIClient) testAPITime() bool {
	since := time.Since(this.Timestamp)
	if int64(since/time.Minute) > 24 {
		return false
	}
	return true
}

func (this *DiaryAPIClient) dorequest(values url.Values, data []byte) (r *http.Request, err error) {
	if values != nil {
		if this.SID != "" && values.Get("sid") == "" {
			values.Add("sid", this.SID)
		}
		this.URL.RawQuery = values.Encode()
	} else {
		values = url.Values{}
		this.URL.RawQuery = values.Encode()
	}
	if data == nil {
		r, err = http.NewRequest("GET", this.URL.String(), nil)

	} else {
		reader := bytes.NewBuffer(data)
		r, err = http.NewRequest("POST", this.URL.String(), reader)
		r.Header.Add("Content-Length", strconv.Itoa(reader.Len()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	r.Header.Add("User-agent", UserAgent)
	return r, err
}

type PostStruct struct {
	postid              string            `json:"postid"`
	juserid             string            `json:"juserid"`
	shortname           string            `json:"shortname"`
	journal_name        string            `json:"journal_name"`
	message_src         string            `json:"message_src"`
	message_html        string            `json:"message_html"`
	author_userid       string            `json:"author_userid"`
	author_shortname    string            `json:"author_shortname"`
	author_username     string            `json:"author_shortname"`
	author_title        string            `json:"author_title"`
	title               string            `json:"title"`
	no_comments         string            `json:"no_comments"`         // Flag for no comments
	comments_count_data string            `json:"comments_count_data"` //Count of comments
	tags_data           map[string]string `json:"tags_data"`
}

func (this *DiaryAPIClient) post_get(shortname string) {
	values := url.Values{}
	values.Add("sid", this.SID)
	values.Add("type", "diary")
	values.Add("method", "post.get")
	values.Add("shortname", shortname)

	r, _ := this.dorequest(values, nil)
	resp, _ := this.HttpClient.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	var message map[string]json.RawMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Fatal(err)
	}
	//Also may try: var posts map[string]*PostStruct
	var posts map[string]PostStruct //It's comfortable, but don't work
	//var posts map[string]interface{} //It's work fine, but noncomfortable
	err = json.Unmarshal(message["posts"], &posts)
	fmt.Println(posts)
	if err != nil {
		log.Fatal(err)
	}
	for id, post_unit := range posts {
		fmt.Println(id, post_unit)
	}
}

func (this *DiaryAPIClient) post_create(title, message string) {
	values := url.Values{}
	values.Add("sid", this.SID)
	values.Add("message", message)
	values.Add("message_src", message)
	values.Add("method", "post.create")
	values.Add("title", title)
	values.Add("close_access_mode", "0")

	this.URL.RawQuery = ""
	resp, err := http.PostForm(this.URL.String(), values)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}*/
	type APIResponse struct {
		Result  int    `json:"result,string"`
		Error   string `json:"error,string"`
		Message string `json:"message"`
		PostID  string `json:"postid"`
	}
	var message APIResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&message)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	userPtr := flag.String("user", "", "Username (Login)")
	passPtr := flag.String("pass", "", "Password")
	flag.Parse()
	var diary DiaryAPIClient
	diary.Init()
	diary.Auth(*userPtr, *passPtr)
	fmt.Println(diary.SID)
	diary.post_create("Test Title", "Test message")
}
