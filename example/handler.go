package main

import (
	"strconv"

	goipmsg "github.com/FlowingSPDG/go-ipmsg"
)

var (
	users    = make(map[string]*goipmsg.ClientData)
	messages = []*goipmsg.ClientData{}
)

func RECEIVE_BR_ENTRY(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	users[cd.Key()] = cd
	ipmsg.SendMSG(cd.Addr, ipmsg.Myinfo(), goipmsg.ANSENTRY)
	return nil
}

func RECEIVE_ANSENTRY(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	users[cd.Key()] = cd
	return nil
}

func RECEIVE_SENDMSG(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	messages = append(messages, cd)

	cmd := cd.Command
	if cmd.Get(goipmsg.SENDCHECK) {
		num := cd.PacketNum
		err := ipmsg.SendMSG(cd.Addr, strconv.Itoa(num), goipmsg.RECVMSG)
		return err
	}
	return nil
}
