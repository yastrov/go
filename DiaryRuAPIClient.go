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
	//"github.com/disintegration/charmap"
	"errors"
	"github.com/yastrov/charmap"
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

func (this *DiaryAPIClient) Auth(user, password string) error {
	strcp1251, _ := charmap.Encode(user, "cp-1251")

	values := url.Values{}
	values.Add("username", strcp1251)
	hash := md5.Sum([]byte(key + password))
	values.Add("password", hex.EncodeToString(hash[:]))
	values.Add("method", "user.auth")
	values.Add("appkey", appkey)

	this.URL, _ = url.Parse(DiaryMainUrl)
	r, err := this.dorequest(values, nil)
	resp, err := this.HttpClient.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	var message DiaryAPIAuthResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&message)
	if err != nil {
		return err
	}
	this.SID = message.SID
	this.Timestamp = time.Now()
	return nil
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
	Avatar_path         string            `json:"avatar_path"`
	Postid              string            `json:"postid"`
	Juserid             string            `json:"juserid"`
	Shortname           string            `json:"shortname"`
	Journal_name        string            `json:"journal_name"`
	Message_src         string            `json:"message_src"`
	Message_html        string            `json:"message_html"`
	Author_userid       string            `json:"author_userid"`
	Author_shortname    string            `json:"author_shortname"`
	Author_username     string            `json:"author_shortname"`
	Author_title        string            `json:"author_title"`
	Title               string            `json:"title"`
	No_comments         string            `json:"no_comments"`         // Flag for no comments
	Comments_count_data string            `json:"comments_count_data"` //Count of comments
	Tags_data           map[string]string `json:"tags_data"`
	Subscribed          string            `json:"subscribed"`
	Can_edit            string            `json:"can_edit"`
	Avatarid            string            `json:"avatarid "`
	No_smile            string            `json:"no_smile"`
	Jaccess             string            `json:"jaccess"`
	Dateline_cdate      string            `json:"dateline_cdate"`
	Close_access_mode2  string            `json:"close_access_mode2"`
	Close_access_mode   string            `json:"close_access_mode"`
	Dateline_date       string            `json:"dateline_date"`
	Access              string            `json:"access"`
}

type DiaryAPIPostGet struct {
	Result int                    `json:"result,string"`
	Posts  map[string]*PostStruct `json:"posts"`
	Error  string                 `json:"error"`
}

func (this *DiaryAPIClient) post_get(shortname, type_, from string) ([]*PostStruct, error) {
	values := url.Values{}
	values.Add("sid", this.SID)
	values.Add("type", type_)
	values.Add("method", "post.get")
	if shortname != "" {
		values.Add("shortname", shortname)
	}
	if from != "" {
		values.Add("from", from)
	}
	r, err := this.dorequest(values, nil)
	resp, err := this.HttpClient.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	var message DiaryAPIPostGet
	err = json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}
	result := make([]*PostStruct, len(message.Posts))
	for id, post_unit := range message.Posts {
		fmt.Println(id, "-", post_unit.Postid)
		result = append(result, post_unit)
	}
	return result, nil
}

func (this *DiaryAPIClient) post_create(title, message string) (string, error) {
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
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	type APIResponse struct {
		Result  int    `json:"result,string"`
		Error   string `json:"error,string"`
		Message string `json:"message"`
		PostID  string `json:"postid"`
	}
	var mess APIResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&mess)
	if err != nil {
		return "", err
	}
	return mess.PostID, nil
}

type DiaryAPICommentGet struct {
	Result   int                       `json:"result,string"`
	Comments map[string]*CommentStruct `json:"comments"`
	Error    string                    `json:"error"`
}

type CommentStruct struct {
	Avatar_path      string `json:"avatar_path"`
	Postid           string `json:"postid"` // For manually control
	Commentid        string `json:"commentid"`
	Shortname        string `json:"shortname"`
	Message_html     string `json:"message_html"`
	Author_userid    string `json:"author_userid"`
	Author_shortname string `json:"author_shortname"`
	Author_avatar    string `json:"author_avatar"`
	Author_username  string `json:"author_shortname"`
	Author_title     string `json:"author_title"`
	Can_edit         string `json:"can_edit"`
	Can_delete       string `json:"can_delete"`
	Dateline         string `json:"dateline"`
}

func (this *DiaryAPIClient) comment_get(postid string) ([]*CommentStruct, error) {
	values := url.Values{}
	values.Add("sid", this.SID)
	values.Add("method", "comment.get")
	values.Add("postid", postid)
	r, err := this.dorequest(values, nil)
	resp, err := this.HttpClient.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var message DiaryAPICommentGet
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&message)
	if err != nil {
		return nil, err
	}
	result := make([]*CommentStruct, len(message.Comments))
	for _, comment_unit := range message.Comments {
		comment_unit.Postid = postid
		result = append(result, comment_unit)
	}
	return result, nil
}

func (this *DiaryAPIClient) comment_get_for_post(post PostStruct) ([]*CommentStruct, error) {
	values := url.Values{}
	values.Add("sid", this.SID)
	values.Add("method", "comment.get")
	values.Add("postid", post.Postid)
	r, err := this.dorequest(values, nil)
	resp, err := this.HttpClient.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var message DiaryAPICommentGet
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&message)
	if err != nil {
		return nil, err
	}
	result := make([]*CommentStruct, len(message.Comments))
	for _, comment_unit := range message.Comments {
		comment_unit.Postid = post.Postid
		result = append(result, comment_unit)
	}
	return result, nil
}

func (this *DiaryAPIClient) comment_get_for_posts(posts []*PostStruct) (result []*CommentStruct, err error) {
	//var result []*PostStruct
	result = make([]*PostStruct, 20)
	for post := range posts {
		if post.Comments_count_data != "" {
			comments, err := comment_get_for_post(*post)
			result = append(result, comments)
		}
	}
	return result, err
}

func main() {
	userPtr := flag.String("user", "", "Username (Login)")
	passPtr := flag.String("pass", "", "Password")
	flag.Parse()
	var diary DiaryAPIClient
	diary.Init()
	err := diary.Auth(*userPtr, *passPtr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(diary.SID)
	//diary.post_create("Test Title", "Test message")
	posts, err := diary.post_get("", "diary", "0")
	if err != nil {
		log.Fatal(err)
	}
	for post := range posts {
		fmt.Println((*post).Avatar_path)
	}
	comments, err := diary.comment_get()
	if err != nil {
		log.Fatal(err)
	}
	for post := range comments {
		fmt.Println(post.PostID)
	}
}
