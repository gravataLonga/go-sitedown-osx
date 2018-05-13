/*
Author: Jonathan Fontes
Site: https://jonathan.pt

To Do:
 - Create configuration based on ini (maybe?)
 - Site monitoring retrive by ini configuration
 - Change time monitoring
 - Change Icon and Sound.
 - Create cli App and put running in background

*/
package main

import (
	"github.com/deckarep/gosx-notifier"
	"net/http"
	"time"
)

type SiteStatus struct {
	Site   string
	Status bool
}

//a slice of string sites that you are interested in watching
var sites []string = []string{
	"https://jonathan.pt"
}

func main() {
	ch := make(chan SiteStatus)

	for _, s := range sites {
		go pinger(ch, s)
	}

	for {
		select {
		case result := <-ch:
			send(&result)
		}
	}
}

func send(result *SiteStatus) {
	var notify = gosxnotifier.NewNotification("The site is down: " + result.Site)
	if result.Status == true {
		notify = gosxnotifier.NewNotification("The site is up: " + result.Site)
	}
	notify.Sound = "Glass"
	notify.Link = result.Site
	notify.Push()
}

func pinger(ch chan SiteStatus, site string) {
	for {
		res, err := http.Get(site)
		var status = false
		if err == nil && res.StatusCode == 200 {
			status = true
		}
		ch <- SiteStatus{
			Site:   site,
			Status: status,
		}
		res.Body.Close()
		time.Sleep(20 * time.Second)
	}
}
