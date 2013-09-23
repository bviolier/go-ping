package main

import "fmt"
import "time"
import "net/http"
import "net"
import "io/ioutil"
import "errors"

// Define Stringer as an interface type with one method, String.
type Pingu struct {
    url string
    interval uint8
    lastCheck time.Time
}

type CheckResult struct {
	target Pingu
    result string
}


func pingtime(c chan time.Time) {
	ticker := time.Tick(1 * time.Second)
	for now := range ticker {
		c <- now
	}
}

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, 2 * time.Second)
}

func doRequest(url string) (error) {
 	transport := http.Transport{
        Dial: dialTimeout,
    }

    client := http.Client{
        Transport: &transport,
    }

    resp,err := client.Get(url)
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

func CheckTarget (c chan CheckResult, now time.Time, target *Pingu) {
	
	if (target.lastCheck.Unix() + int64(target.interval) < now.Unix()) {
		target.lastCheck = now
		
		err := doRequest(target.url)

		if err != nil {
			c <- CheckResult{target: *target, result: "NOT OK:" + err.Error()}
		} else {
			c <- CheckResult{target: *target, result: "OK"}
		}
	} 
}


func main() {
	fmt.Println("Go Ping Pong made in China")

	var targets []Pingu

	targets = append(targets, Pingu{url:"http://127.0.0.1", interval:5})
	targets = append(targets, Pingu{url:"http://127.0.0.1/jemoeder", interval:7})
	c := make(chan time.Time)
	results := make(chan CheckResult)

	go pingtime(c)

	for { 
		select {
			case q := <- c:
				fmt.Println("Checking all the things! ", q)
				
				for i,_ := range targets {
					go CheckTarget(results, q, &targets[i])
				}


			case result := <- results:
				fmt.Println(result.target.url + ":" + result.result)
		}
	}


	// channel with timer	

	// ->

	// 	foreach json 

	// 		go ping

	// 			: timer? 10000000000000

	// 		result <- 

	// 		output
}