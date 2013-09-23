package main

import (
	"io/ioutil"
	"fmt"
	"os"
   	"encoding/json"
)

type ServersConfig struct {
    Servers[] Server
}

type Server struct {
	Domain string
	Timeout int
	Interval int
}


func main() {
    configFile, err := ioutil.ReadFile("servers.json")
	if err != nil {
        fmt.Printf("File error: %v\n", err) 
        os.Exit(1) 
	}
	
	ServerList := new(ServersConfig)
   	json.Unmarshal(configFile, &ServerList)
	
	fmt.Println("Going to check the following servers:")
	
	fmt.Printf("%3s  %20s  %7s  %8s  %10s\n", "#", "Domain", "Timeout", "Interval", "Last state")
	for index, value := range ServerList.Servers {
		showServer(index, value, lastState)
	}
	
}

func showServer(index int, value Server, lastState) {
	fmt.Printf("%3d  %20s  %7d  %8d %10s\n", index+1, value.Domain, value.Timeout, value.Interval, lastState)	
}
