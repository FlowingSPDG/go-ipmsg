package ipmsg

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

// IPMSG IPMessager main struct.
type IPMSG struct {
	ClientData ClientData
	Conn       *net.UDPConn
	Conf       *Config
	Handlers   []*EventHandler
	PacketNum  int
}

// Config IPMessager configration struct.
type Config struct {
	NickName  string
	GroupName string
	UserName  string
	HostName  string
	Port      int
	Local     string
}

const (
	// DefaultPort default UDP Port number(2425).
	DefaultPort int = 2425
	// Buflen net.Conn data receive buffer size(65535).
	Buflen int = 65535
)

// NewIPMSGConf Generate New pointer of Config struct.
func NewIPMSGConf() *Config {
	return &Config{
		Port: DefaultPort,
	}
}

// NewIPMSG Generate new instance pointer of IPMSG. Return error if it failed to solve/listen UDP Address.
func NewIPMSG(conf *Config) (*IPMSG, error) {
	ipmsg := &IPMSG{
		PacketNum: 0,
	}
	ipmsg.Conf = conf
	// UDP server
	service := fmt.Sprintf("%v:%d", conf.Local, conf.Port)
	//fmt.Println("service =", service)
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		return ipmsg, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return ipmsg, err
	}
	ipmsg.Conn = conn
	return ipmsg, err
}

// Close Close UDP connection
func (ipmsg *IPMSG) Close() error {
	conn := ipmsg.Conn
	if conn == nil {
		err := errors.New("Conn is not defined")
		return err
	}
	return conn.Close()
}

// BuildData Generates new pointer of ClientData struct.
func (ipmsg *IPMSG) BuildData(addr *net.UDPAddr, msg string, cmd Command) *ClientData {
	conf := ipmsg.Conf
	clientdata := NewClientData("", addr)
	clientdata.Version = 1
	clientdata.PacketNum = ipmsg.GetNewPacketNum()
	clientdata.User = conf.UserName
	clientdata.Host = conf.HostName
	clientdata.Command = cmd
	clientdata.Option = msg
	return clientdata
}

// SendMSG Send UDP message/cmd to specified addr.
func (ipmsg *IPMSG) SendMSG(addr *net.UDPAddr, msg string, cmd Command) error {
	clientdata := ipmsg.BuildData(addr, msg, cmd)
	conn := ipmsg.Conn
	_, err := conn.WriteToUDP([]byte(clientdata.String()), addr)
	if err != nil {
		return err
	}
	return nil
}

// RecvMSG Receive message from UDP.
func (ipmsg *IPMSG) RecvMSG() (*ClientData, error) {
	var buf [Buflen]byte
	conn := ipmsg.Conn
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return nil, err
	}
	trimmed := bytes.Trim(buf[:], "\x00")
	clientdata := NewClientData(string(trimmed[:]), addr)

	handlers := ipmsg.Handlers
	for _, v := range handlers {
		err := v.Run(clientdata, ipmsg)
		if err != nil {
			return clientdata, err
		}
	}
	return clientdata, nil
}

// UDPAddr convert net.Addr to net.UDPAddr
func (ipmsg *IPMSG) UDPAddr() (*net.UDPAddr, error) {
	conn := ipmsg.Conn
	if conn == nil {
		err := errors.New("Conn is not defined")
		return nil, err
	}
	addr := conn.LocalAddr()
	network := addr.Network()
	str := addr.String()
	//fmt.Println("str =", str)
	udpAddr, err := net.ResolveUDPAddr(network, str)
	return udpAddr, err
}

// AddEventHandler Add new event handler
func (ipmsg *IPMSG) AddEventHandler(ev *EventHandler) {
	ipmsg.Handlers = append(ipmsg.Handlers, ev)
}

// GetNewPacketNum Get new packet number with Unix Time.
func (ipmsg *IPMSG) GetNewPacketNum() int {
	ipmsg.PacketNum++
	return int(time.Now().Unix()) + ipmsg.PacketNum
}

// Myinfo Get Nickname[x00]groupname[x00] string.
func (ipmsg *IPMSG) Myinfo() string {
	conf := ipmsg.Conf
	return fmt.Sprintf("%s\x00%s\x00", conf.NickName, conf.GroupName)
}
