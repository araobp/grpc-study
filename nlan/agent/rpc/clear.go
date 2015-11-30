package rpc

import (
	"github.com/araobp/go-nlan/nlan/agent/context"
	"github.com/araobp/go-nlan/nlan/util"
)

func Clear(con *context.Context) {
	logger := con.Logger
	logger.Print("Clearing...")
	cmd, _ := con.GetCmd()
	l := util.GetBridgeNames()
	for e := l.Front(); e != nil; e = e.Next() {
		v := e.Value
		bridge := v.(string)
		cmd("ovs-vsctl", "del-br", bridge)
	}
}
