package main

import (
	"fmt"
	"github.com/coreswitch/cmd"
	"log"
	"strings"
)

type LagoBr struct {
	Name       string
	Dpid       string
	Controller []string
	Port       map[string]*LagoPort
	FailMode   string
	Create     bool
	Change     bool
}

var Bridges = make(map[string]*LagoBr)

func GetBridge(name string) *LagoBr {
	_, ok := Bridges[name]
	if !ok {
		Bridges[name] = &LagoBr{name, "", []string{}, make(map[string]*LagoPort), "", true, false}
	}
	return Bridges[name]
}

func ConfigBrDpid(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	dpid := Args[1].(string)

	br := GetBridge(name)
	switch Cmd {
	case cmd.Set:
		if br.Dpid != dpid {
			br.Change = true
		}
		br.Dpid = dpid
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func ConfigBrFailMode(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	mode := Args[1].(string)

	br := GetBridge(name)
	switch Cmd {
	case cmd.Set:
		if br.FailMode != mode {
			br.Change = true
		}
		br.FailMode = mode
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func ConfigBrCon(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	cons := Args[1].(string)

	br := GetBridge(name)
	switch Cmd {
	case cmd.Set:
		br.Change = true
		br.Controller = strings.Split(cons, " ")
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func GetBrPort(name string, br *LagoBr) *LagoPort {
	_, ok := br.Port[name]
	if !ok {
		br.Port[name] = &LagoPort{name, "", "", true, false}
	}
	return br.Port[name]
}

func ConfigBrPort(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	port := Args[1].(string)
	id := Args[2].(string)

	br := GetBridge(name)
	switch Cmd {
	case cmd.Set:
		br.Change = true
		br.Port[port] = GetBrPort(port, br)
		br.Port[port].OFPort = id
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func MakeBrDSL(br *LagoBr) string {
	dsl := []string{"bridge", br.Name, "", "-dpid", br.Dpid}
	if br.Create {
		dsl[2] = "create"
	} else if br.Change {
		dsl[2] = "config"
	} else {
		dsl = []string{}
	}
	br.Create = false
	br.Change = false

	for _, con := range br.Controller {
		dsl = append(dsl, fmt.Sprint("-controller ", con))
	}
	for _, port := range br.Port {
		dsl = append(dsl, fmt.Sprint("-port ", port.Name, " ", port.OFPort))
	}
	if br.FailMode != "" {
		dsl = append(dsl, fmt.Sprint("-fail-mode ", br.FailMode))
	}
	return strings.Join(dsl, " ")
}

func SetBrDSL() ([]string, error) {
	dsl := make([]string, len(Bridges))
	i := 0
	for _, br := range Bridges {
		dsl[i] = MakeBrDSL(br)
		i++
	}
	return dsl, nil
}
