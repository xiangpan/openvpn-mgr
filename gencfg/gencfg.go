package gencfg

import (
	"html/template"
	"os"

	"bufio"
	"fmt"
	"strings"

	"net"

	"github.com/siddontang/ledisdb/Godeps/_workspace/src/github.com/siddontang/go/log"
    "regexp"
)

type VpnCfg struct {
	Proto       string
	Port        string
	VpnSubnet   string
	PushRoute   []string
	DuplicateCN string
	MaxClients  string
}

func Input(msg, value string, p interface{}) error {
	fmt.Printf("%s\n", msg)
	fmt.Printf("default value is %s :", value)
	r := bufio.NewReader(os.Stdin)
	t, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	t = strings.TrimSpace(t)
	if t == "" {
		t = value
	}
	switch p1 := p.(type) {
	case *string:
		*p1 = t
	case *[]string:
		*p1 = strings.Split(t, ",")
	}
	return nil
}

//取网卡网段
func GetInterFace() (ret string) {
	i, _ := net.InterfaceAddrs()
	for _, v := range i {
		tmp := v.String()
		if ok, _ := regexp.MatchString(`^((192\.168\.)|(10\.)|(^172\.((1[6-9])|(2[0-9])|(3[01]))\.))`, tmp); ! ok {
			continue
		}
		reg := regexp.MustCompile(`[0-9]+/.*$`)
		tmp = reg.ReplaceAllString(tmp, "0")
		if ret == "" {
			ret = tmp + " 255.255.255.0"
		} else {
			ret = ret + "," + tmp + " 255.255.255.0"
		}
	}
	return ret
}

func GenCfg() {
	t := &VpnCfg{
		Proto:       "tcp",
		Port:        "11941",
		VpnSubnet:   "10.8.0.0",
		PushRoute:   []string{},
		DuplicateCN: "duplicate-cn",
		MaxClients:  "50",
	}
	tmpl, err := template.ParseFiles("template/server.cfg")
	if err != nil {
		log.Fatal(err)
	}
	fp, err:= os.Create("/etc/openvpn/server.conf")
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	Input("Which local IP address should OpenVPN listen on", "11941", &t.Port)
	Input("TCP or UDP server", "tcp", &t.Proto)
	Input("Configure server mode and supply a VPN subnet", "10.8.0.0 255.255.255.0", &t.VpnSubnet)
	Input("Push routes to the client to allow it to reach other private subnets behind the server",
		GetInterFace(),
        &t.PushRoute)
	Input("multiple clients might connect with the same certificate/key", "duplicate-cn", &t.DuplicateCN)
	Input("The maximum number of concurrently connected clients", "50", &t.MaxClients)
	tmpl.Execute(fp, t)
}