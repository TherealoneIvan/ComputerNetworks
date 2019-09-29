package main

// Server gives info about server to connect
type Server struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func main() {

}
