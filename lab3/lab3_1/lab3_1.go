package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
)

// Password represents encrypted password
type Password struct {
	Secret     []byte `json:"secret"`
	Passphrase string `json:"passphrase"`
}

// Config represents settings of connection to SMTP server
type Config struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Sender   string   `json:"sender"`
	Password Password `json:"password"`
}

func getConfig(path string) Config {
	// открытие JSON файла конфигурации
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// считывание информации из файла в массив байтов
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	// десериализация JSON строки в экземпляр структуры
	json.Unmarshal(byteValue, &config)

	return config
}

func decryptPassword(p Password) string {
	// дешифрация пароля
	key := []byte(p.Passphrase)
	ciphertext := p.Secret

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	password, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}

	return string(password)
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "path to config file")
	flag.Parse()

	config := getConfig(configPath)

	// авторизация
	auth := smtp.PlainAuth("", config.Sender, decryptPassword(config.Password), config.Host)

	// формирование сообщения
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите To:")
	scanner.Scan()
	to := scanner.Text()

	fmt.Println("Введите Subject:")
	scanner.Scan()
	subject := scanner.Text()

	fmt.Println("Введите MessageBody:")
	scanner.Scan()
	messageBody := scanner.Text()

	// отправка сообщения
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		auth,
		config.Sender,
		[]string{
			to,
		},
		[]byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, messageBody)),
	)
	if err != nil {
		log.Fatal(err)
	}
}
