package ssdp

import (
	"os/exec"

	ssdp "github.com/koron/go-ssdp"
	"github.com/rs/xid"
)

// Start Spawns a new goroutine SSDP server with the specified ST and server name
// It returns any errors that have occurred
func Start(st string, serverName string) {
	id := xid.New()
	go ssdp.Advertise(st, id.String(), "http://"+currentIP()+":3001/description.xml", serverName, 1800)
}

func currentIP() string {
	ip, err := exec.Command("hostname", "-I").Output()

	if err != nil {
		panic(err)
	}

	return string(ip)
}
