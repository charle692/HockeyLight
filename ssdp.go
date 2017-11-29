package main

import (
	ssdp "github.com/koron/go-ssdp"
	"github.com/rs/xid"
	"os/exec"
)

func startSSDPServer() {
	id := xid.New()

	_, err := ssdp.Advertise("my:hockey-light", id.String(), "http://" + currentIP() + ":3001/description.xml", "Hockey Light SSDP", 1800)

	if err != nil {
		panic(err)
	}

	// run Advertiser infinitely.
	quit := make(chan bool)
	<-quit
}

func currentIP() string {
	ip, err := exec.Command("hostname", "-I").Output()

	if err != nil {
		panic(err)
	}

	return string(ip)
}
