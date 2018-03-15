package main

import (
	"github.com/coreswitch/cmd"
	"log"
	"strings"
)

type LagoPort struct {
	Name      string
	Interface string
	OFPort    string
	Create    bool
	Change    bool
}

var Ports = make(map[string]*LagoPort)

func GetPort(name string) *LagoPort {
	_, ok := Ports[name]
	if !ok {
		Ports[name] = &LagoPort{name, "", "", true, false}
	}
	return Ports[name]
}

func ConfigPortInterface(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	iface := Args[1].(string)

	port := GetPort(name)
	switch Cmd {
	case cmd.Set:
		if port.Interface != iface {
			port.Change = true
		}
		port.Interface = iface
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func MakePtDSL(port *LagoPort) string {
	dsl := []string{"port", port.Name, "", "-interface", port.Interface}
	if port.Create {
		dsl[2] = "create"
	} else if port.Change {
		dsl[2] = "config"
	} else {
		dsl = []string{}
	}
	port.Create = false
	port.Change = false
	return strings.Join(dsl, " ")
}

func SetPtDSL()([]string, error) {
	dsl := make([]string, len(Ports))
	i := 0
	for _, port := range Ports {
		dsl[i] = MakePtDSL(port)
		i++
	}
	return dsl, nil
}
