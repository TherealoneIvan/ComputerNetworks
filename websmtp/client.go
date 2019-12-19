package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
)

type authData struct {
	login, pass, mail string
}

var userData = make(map[string]*authData, 1024)

func getAuthData(r *http.Request) (*authData, error) {
	token, err := r.Cookie("id_token")
	if err != nil {
		md := md5.Sum([]byte(r.RemoteAddr))
		token := hex.EncodeToString(md[:])
		auth := userData[token]
		if auth == nil {
			return nil, err
		}
		return auth, nil

	}
	return userData[token.Value], nil
}

func auth(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "auth.html")
}

func initiation(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	md := md5.Sum([]byte(r.RemoteAddr))
	token := hex.EncodeToString(md[:])
	r.AddCookie(&http.Cookie{
		Name:  "id_token",
		Value: token,
	})
	userData[token] = &authData{
		login: r.FormValue("login"),
		pass:  r.FormValue("pass"),
		mail:  r.FormValue("mail"),
	}
	http.Redirect(w, r, "/index", http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func send(w http.ResponseWriter, r *http.Request) {
	user, err := getAuthData(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
	_ = r.ParseForm()
	var login = user.login
	var pass = user.pass
	from := mail.Address{
		Name:    "",
		Address: login,
	}
	to := mail.Address{
		Name:    "",
		Address: r.FormValue("to"),
	}
	subject := r.FormValue("subject")
	msgBody := r.FormValue("msg")

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + msgBody
	servername := "smtp." + user.mail + ":587"
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", login, pass, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	c, err := smtp.Dial(servername)
	if err != nil {
		log.Panic(err)
	}

	err = c.StartTLS(tlsConfig)
	if err != nil {
		log.Panic(err)
	}

	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	writer, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = writer.Close()
	if err != nil {
		log.Panic(err)
	}

	_ = c.Quit()

	http.Redirect(w, r, "/index", http.StatusFound)
}

func main() {
	http.HandleFunc("/", auth)
	http.HandleFunc("/init", initiation)
	http.HandleFunc("/index", index)
	http.HandleFunc("/send", send)
	if err := http.ListenAndServe(":2508", nil); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
