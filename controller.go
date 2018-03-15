package main

import (
	"github.com/coreswitch/cmd"
	"log"
	"strings"
)

type LagoCon struct {
	Name           string
	Channel        string
	Role           string
	ConnectionType string
	Create         bool
	Change         bool
}

var Controllers = make(map[string]*LagoCon)

func GetController(name string) *LagoCon {
	_, ok := Controllers[name]
	if !ok {
		Controllers[name] = &LagoCon{name, "", "", "", true, false}
	}
	return Controllers[name]
}

func ConfigConChannel(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	channel := Args[1].(string)

	con := GetController(name)
	switch Cmd {
	case cmd.Set:
		if con.Channel != channel {
			con.Change = true
		}
		con.Channel = channel
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func ConfigConRole(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	role := Args[1].(string)

	con := GetController(name)
	switch Cmd {
	case cmd.Set:
		if con.Role != role {
			con.Change = true
		}
		con.Role = role
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func ConfigConConnectType(Cmd int, Args cmd.Args) int {
	name := Args[0].(string)
	contype := Args[1].(string)

	con := GetController(name)
	switch Cmd {
	case cmd.Set:
		if con.ConnectionType != contype {
			con.Change = true
		}
		con.ConnectionType = contype
	case cmd.Delete:
		log.Println("Not Support: Delete")
	}
	return cmd.Success
}

func MakeConDSL(con *LagoCon) string {
	dsl := []string{"controller", con.Name, "", "-role", con.Role, "-channel", con.Channel, "-connection-type", con.ConnectionType}
	if con.Create {
		dsl[2] = "create"
	} else if con.Change {
		dsl[2] = "config"
	} else {
		dsl = []string{}
	}
	con.Create = false
	con.Change = false
	return strings.Join(dsl, " ")
}

func SetConDSL() ([]string, error) {
	dsl := make([]string, len(Controllers))
	i := 0
	for _, con := range Controllers {
		dsl[i] = MakeConDSL(con)
		i++
	}
	return dsl, nil
}
