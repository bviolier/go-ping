package main

import (
	"fmt"
	"time"
	"net/http"
	"net"
	"io/ioutil"
	"errors"
	"encoding/json"
	"os"
	"os/exec"
)


type ServersConfig struct {
    Servers[] Server
}

type Server struct {
	Domain string
	Timeout int
	Interval int
	LastCheck time.Time
	LastState string
}

type CheckResult struct {
	target *Server
    result string
}


func pingtime(c chan time.Time) {
	ticker := time.Tick(1 * time.Second)
	for now := range ticker {
		c <- now
	}
}


func doRequest(server Server) (error) {
	mydialout := func(network, addr string)(net.Conn, error) { 
			return net.DialTimeout(network, addr, time.Duration(server.Timeout) * time.Second)
		}
 	transport := http.Transport{
        Dial: mydialout,
    }

    client := http.Client{
        Transport: &transport,
    }

    resp,err := client.Get(server.Domain)
    if err != nil {
    	return err
    }

    defer resp.Body.Close()

	_,err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if (resp.StatusCode != 200) {
		return errors.New("HTTP STATUS != 200")
	}

	return nil
}

func CheckTarget (c chan CheckResult, now time.Time, target *Server) {
	
	if (target.LastCheck.Unix() + int64(target.Interval) <= now.Unix()) {
		target.LastCheck = now
		
		err := doRequest(*target)

		if err != nil {
			c <- CheckResult{target: target, result: "NOT OK: " + err.Error()}
		} else {
			c <- CheckResult{target: target, result: "OK"}
		}
	} 
}


func main() {
	fmt.Println("Go Ping Pong made in China")
	
    configFile, _ := ioutil.ReadFile("servers.json")
	ServerList := new(ServersConfig)
   	json.Unmarshal(configFile, &ServerList)
	
	fmt.Println("Going to check the following servers:")
	
	showServers(ServerList);
	
	c := make(chan time.Time)
	results := make(chan CheckResult)

	go pingtime(c)

	for { 
		select {
			case q := <- c:				
				for i,_ := range ServerList.Servers {
					go CheckTarget(results, q, &ServerList.Servers[i])
				}


			case result := <- results:
				result.target.LastState = result.result
				showServers(ServerList)
		}
	}
}

func showServers(ServerList *ServersConfig) {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	fmt.Printf("%3s  %30s  %7s  %8s  %13s %13s\n", "#", "Domain", "Timeout", "Interval", "Last check", "Last state")
	for index, value := range ServerList.Servers {
		showServer(index, value)
	}
	
}

func showServer(index int, value Server) {
	const layout = "15:04:05"
	fmt.Printf("%3d  %30s  %7d  %8d  %13s    %s\n", index+1, value.Domain, value.Timeout, value.Interval, value.LastCheck.Format(layout), value.LastState)	
}
