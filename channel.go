package main

import (
	"fmt"
	"github.com/coreswitch/cmd"
	"log"
	"strings"
)

type LagoCh struct {
	Name     string
	DstAddr  string
	DstPort  string
	Protocol string
	Create   bool
	Change   bool
}

var Channels = make(map[string]*LagoCh)

func GetChannel(name string) *LagoCh {
	_, ok := Channels[name]
	if !ok {
		Channels[name] = &LagoCh{name, "", "", "", true, false}
	}
	return Channels[name]
}

func ConfigChannelDSTAddr(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	addr := Args[1].(string)

	channel := GetChannel(name)
	switch Cmd {
	case cmd.Set:
		if channel.DstAddr != addr {
			channel.Change = true
		}
		channel.DstAddr = addr
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func ConfigChannelDSTProtocol(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	protocol := Args[1].(string)

	channel := GetChannel(name)
	switch Cmd {
	case cmd.Set:
		if channel.Protocol != protocol {
			channel.Change = true
		}
		channel.Protocol = protocol
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func ConfigChannelDSTPort(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	port := Args[1].(string)

	channel := GetChannel(name)
	switch Cmd {
	case cmd.Set:
		if channel.DstPort != port {
			channel.Change = true
		}
		channel.DstPort = port
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func MakeChDSL(channel *LagoCh) string {
	dsl := []string{"channel", channel.Name, "", "-dst-addr", channel.DstAddr, "-protocol", channel.Protocol}
	if channel.Create {
		dsl[2] = "create"
	} else if channel.Change {
		dsl[2] = "config"
	} else {
		dsl = []string{}
	}
	if channel.DstPort != "" {
		dsl = append(dsl, fmt.Sprint("-dst-port ", channel.DstPort))
	}
	channel.Create = false
	channel.Change = false
	return strings.Join(dsl, " ")
}

func SetChDSL()([]string, error) {
	dsl := make([]string, len(Channels))
	i := 0
	for _, channel := range Channels {
		dsl[i] = MakeChDSL(channel)
		i++
	}
	return dsl, nil
}
