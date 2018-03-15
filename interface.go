package main

import (
	"log"
	"strings"
	"github.com/coreswitch/cmd"
)


type LagoIf struct {
	Name   string
	Type   string
	Device string
	Create bool
	Change bool
}

var Interfaces = make(map[string]*LagoIf)


func GetIface(name string) *LagoIf {
	_, ok := Interfaces[name]
	if !ok {
		Interfaces[name] = &LagoIf{name, "", "", true, false}
	}
	return Interfaces[name]
}

func ConfigInterfaceType(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	iftype := Args[1].(string)

	iface := GetIface(name)
	switch Cmd {
	case cmd.Set:
		if iface.Type != iftype {
			iface.Change = true
		}
		iface.Type = iftype
	case cmd.Delete:
		log.Println("Not Support: Interface Delete")
	}
	return cmd.Success
}

func ConfigInterfaceDevice(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	ifdevice := Args[1].(string)

	iface := GetIface(name)
	switch Cmd {
	case cmd.Set:
		if Interfaces[name].Device != ifdevice {
			iface.Change = true
		}
		iface.Device = ifdevice
	case cmd.Delete:
		log.Println("Not Support: Interface Delete")
	}
	return cmd.Success
}

func MakeIfDSL(lagoif *LagoIf) string{
	dsl := []string{"interface", lagoif.Name, "", "-type", lagoif.Type, "-device", lagoif.Device}
	if lagoif.Create {
		dsl[2] = "create"
	} else if lagoif.Change {
		dsl[2] = "config"
	} else {
		dsl = []string{}
	}
	lagoif.Create = false
	lagoif.Change = false
	return strings.Join(dsl, " ")
}

func SetIfDSL()([]string, error) {
	dsl := make([]string, len(Interfaces))
	i := 0
	for _, lagoif := range Interfaces {
		dsl[i] = MakeIfDSL(lagoif)
		i++
	}
	return dsl, nil
}
