package main

import (
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

type client struct{
	DSL []string
}

func (c client) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	for _, dsl := range c.DSL {
		oi.LongWrite(w, []byte(dsl))
	}
}
