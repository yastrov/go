package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"time"
)

const (
	urlStr = "http://some.service.com" //Service we check
	HourDelay = 25 // How many hours file is fresh
	/*Email part*/
	SubjectConst = "My Subject"
	BodyConst = "My\nBody"
	MyEmailConst = "me@service.com"
	MyPassworldConst = "pass"
	ToEmailConst = "he@service.com"
	SMTPHostConst = "smtp.service.com:465"
)

// SSL/TLS Email Example
/*Thanks for https://github.com/chrisgillis
https://gist.github.com/chrisgillis/10888032
*/
func SendEmail(From, To, subj, body, Login, Password, HostPort string) {
	from := mail.Address{"", From}
	to := mail.Address{"", To}

	// Setup headers
	headers := make(map[string]string, 3)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := HostPort
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", Login, Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}
	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}
	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}
	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}
	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}
	err = w.Close()
	if err != nil {
		log.Panic(err)
	}
	c.Quit()
}

func EmailMe(LastMod string) {
	subj := SubjectConst + LastMod
	body := BodyConst
	SendEmail(MyEmailConst, ToEmailConst, subj, body, MyEmailConst, MyPassworldConst, SMTPHostConst)
}

func ParseTime(t string) (d time.Time, err error) {
	d, err = time.Parse(time.RFC1123, t)
	if err != nil {
		return d, nil
	}
	return
}

func main() {
	resp, _ := http.Head(urlStr)
	defer resp.Body.Close()
	LastMod := resp.Header["Last-Modified"][0]
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	LastModDay, _ := ParseTime(LastMod)
	since := time.Since(LastModDay)
	hoursSince := int64(since / time.Hour)
	if hoursSince < HourDelay {
		log.Println("File Updated")
		EmailMe(LastMod)
	} else {
		log.Println("Old File")
	}
}
