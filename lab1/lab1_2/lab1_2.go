package main

import (
	filedriver "gitea.com/goftp/file-driver"
	"goftp.io/server"
)

func main() {
	var (
		user = "admin"
		pass = "123456"
	)

	factory := &filedriver.FileDriverFactory{
		RootPath: "ftproot",
		Perm:     server.NewSimplePerm("root", "root")}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     21,
		Hostname: "185.20.227.83",
		Auth:     &server.SimpleAuth{Name: user, Password: pass}}
	server := server.NewServer(opts)
	server.ListenAndServe()
}
