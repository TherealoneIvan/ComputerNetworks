package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHIdentity represents data for auth
type SSHIdentity struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Output is a wrapper for response from this server
type Output struct {
	List []string `json:"list"`
}

func parseCommand(s string) []string {
	s = strings.TrimRight(s, "\n")
	re := regexp.MustCompile(`\s+`)
	re.ReplaceAllString(s, " ")
	data := strings.Split(s, " ")
	return data
}

// AuthRouterHandler : Запись данных для авторизации на SSH сервере
func AuthRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	i, err := strconv.Atoi(r.Form["port"][0])
	if err != nil {
		fmt.Println(err)
		return
	}

	identity := SSHIdentity{
		Host:     r.Form["host"][0],
		Port:     i,
		Login:    r.Form["login"][0],
		Password: r.Form["password"][0],
	}

	jsonBytes, _ := json.Marshal(identity)

	http.SetCookie(w, &http.Cookie{
		Name:    "identity",
		Value:   base64.StdEncoding.EncodeToString(jsonBytes),
		Expires: time.Now().Add(30 * time.Minute),
	})

	http.Redirect(w, r, "/console", 301)
}

// ExecuteRouterHandler : Выполнение команды на сервере на SSH сервере
func ExecuteRouterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	var identity SSHIdentity
	response := Output{
		List: []string{},
	}
	cmd := parseCommand(r.Form["command"][0])
	divider := fmt.Sprintf("cmd:%s", strings.Join(cmd, " "))

	cookie, err := r.Cookie("identity")
	if err != nil {
		result, _ := json.Marshal(response)
		w.Write(result)
		return
	}

	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		result, _ := json.Marshal(response)
		w.Write(result)
		return
	}

	err = json.Unmarshal(data, &identity)
	if err != nil {
		result, _ := json.Marshal(response)
		w.Write(result)
		return
	}

	// конфигурация подключения
	config := &ssh.ClientConfig{
		User: identity.Login,
		Auth: []ssh.AuthMethod{
			ssh.Password(identity.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}}
	// подключение
	conn, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%d", identity.Host, identity.Port), config)
	session, _ := conn.NewSession()
	defer session.Close()
	// получение управления вводом и выводом
	stdin, _ := session.StdinPipe()
	stdout, _ := session.StdoutPipe()
	// старт оболочки на удалённой машине
	err = session.Shell()
	if err != nil {
		result, _ := json.Marshal(response)
		w.Write(result)
		return
	}
	// отправка команд
	fmt.Fprintf(stdin, "echo %s\n", divider)
	fmt.Fprintln(stdin, strings.Join(cmd, " "))
	fmt.Fprintln(stdin, "exit")
	// чтение вывода
	out, _ := ioutil.ReadAll(stdout)
	outstr := string(out)
	fmt.Println(outstr)
	response.List = strings.Split(outstr[strings.Index(outstr, divider+"\n")+len(divider)+1:], "\n")

	result, _ := json.Marshal(response)
	w.Write(result)
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css")))) // доступ к стилям
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))    // доступ к стилям
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/auth.html")
	}) // возврат страница авторизации по запросу домашней страницы

	// установка обработчиков запросов
	http.HandleFunc("/auth", AuthRouterHandler)
	http.HandleFunc("/console", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("identity")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/", 301)
			}
		}
		http.ServeFile(w, r, "static/console.html")
	})
	http.HandleFunc("/execute", ExecuteRouterHandler)

	err := http.ListenAndServe(":8017", nil) // задаём слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
