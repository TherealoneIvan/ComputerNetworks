# Лабораторная работа № 3.1

## Разработка smtp клиента


Рассматривается задача разработки smtp-клиента на языке GO. Используя пакеты net/smtp, crypto/tls, а также в зависимости от реализации возможно понадобится strings.  Необходимо реализовать задачи приведенные ниже.
 
**Задача 1:** Реализовать smtp-клиент на GO и запустить его на сервере 185.20.227.83 или 185.20.226.174. Для данного приложения необходимо создать тестовую учетную запись на почтовом сервере например на [yandex.ru](https://yandex.ru) или [gmail.com](https://gmail.com), после чего используя параметры SMTP соединения с сервером соответствующего почтового сервиса выполнить отправку тестового сообщения. 

**Задача 2:** Протестировать работу приложения путем отправки тестового сообщения на ящик danila@posevin.com в зависимости от группы и преподавателя, в котором указать фамилию студента и группу. 

**Задача 3:** Реализовать следующие функции: 

1. ввод значения поля *To* из командной строки;
2. ввод значения поля *Subject* из командной строки;
3. ввод сообщения в поле *Message body* из командной строки;

Дополнительные требования:

1. параметры соединения должны храниться во внешнем файле;
2. пароль соединения не должен храниться в явном виде, а должен быть зашифрован любым известным алгоритмом шифрования.

**Замечание 1:** Для корректной работы GO на серверах 185.20.227.83, 185.20.226.174 необходимо задать переменную окружения 
```bash
export GOPATH=~/go.
```

# Лабораторная работа № 3.2

## Разработка ICMP приложений

Рассматривается задача разработки приложения посылающего ICMP пакеты на языке GO. Необходимо реализовать ping приложение. Для реализации возможно использовать один из следующих пакетов [go-fastping](https://github.com/tatsushid/go-fastping) или [go-ping](https://github.com/sparrc/go-ping). Дополнительно необходимо рассмотреть пакет [trace](https://golang.org/x/net/trace) позволяющий выполнить трассировку до указанного хоста.
 
**Пример 1:** Посылаем хосту три ICMP пакета используя go-ping:
```go
pinger, err := ping.NewPinger("www.google.com")
if err != nil {
	panic(err)
}

pinger.Count = 3
pinger.Run()                 // blocks until finished
stats := pinger.Statistics() // get send/receive/rtt stats
```
**Пример 2:** Моделирование Unix ping команды на основе go-ping:
```go
pinger, err := ping.NewPinger("www.google.com")
if err != nil {
	fmt.Printf("ERROR: %s\n", err.Error())
	return
}

pinger.OnRecv = func(pkt *ping.Packet) {
	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
}
pinger.OnFinish = func(stats *ping.Statistics) {
	fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
	fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
}

fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
pinger.Run()
```
**Пример 3.** Моделирование Unix ping команды на основе go-fastping:
```go
p := fastping.NewPinger()
ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
p.AddIPAddr(ra)
p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
	fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
}
p.OnIdle = func() {
	fmt.Println("finish")
}
err = p.Run()
if err != nil {
	fmt.Println(err)
}
```
**Задача 1:** Реализовать приложение посылающее ICMP пакеты к заданному хосту и выводящее результаты ответа.

**Задача 2:** Реализовать вид DDoS-атаки типа ICMP-флуд используя горутины http://golang-book.ru/chapter-10-concurrency.html

**Задача 3.** Реализовать трассировку до заданного хоста.

**Замечание 1:** Для корректной работы GO на серверах 185.20.227.83, 185.20.226.174 необходимо задать переменную окружения 
```bash
export GOPATH=~/go.
```

**Источники информации**

1. https://godoc.org/golang.org/x/net/trace
2. https://godoc.org/github.com/sparrc/go-ping
3. https://github.com/tatsushid/go-fastping