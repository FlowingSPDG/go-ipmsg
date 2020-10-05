package ipmsg

import "fmt"

// EvFunc is Event handler type definition.
type EvFunc func(cd *ClientData, ipmsg *IPMSG) error

// EventHandler is handler struct
type EventHandler struct {
	String   string
	Handlers map[Command]EvFunc
	Debug    bool
}

// NewEventHandler generates new pointer of EventHandler{}
func NewEventHandler() *EventHandler {
	return &EventHandler{
		Handlers: make(map[Command]EvFunc),
	}
}

// Regist Regist new EventFunc handler with specified Command(cmd).
func (ev *EventHandler) Regist(cmd Command, evfunc EvFunc) {
	handlers := ev.Handlers
	handlers[cmd] = evfunc
}

// Run Run specified event handler
func (ev *EventHandler) Run(cd *ClientData, ipmsg *IPMSG) error {
	if ev.Debug {
		ev.RunDebug(cd)
	}
	cmd := cd.Command.Mode()
	evfunc := ev.Handlers[cmd]
	if evfunc == nil {
		// just do nothing when handler is undefined
		return nil
	}
	return (evfunc(cd, ipmsg))
}

// RunDebug Run event handler with Debug output.
func (ev *EventHandler) RunDebug(cd *ClientData) {
	cmdstr := cd.Command.Mode().String()
	fmt.Println("EventHandler.RunDebug cmdstr=", cmdstr)
	fmt.Println("EventHandler.RunDebug key=", cd.Key())
}

//func (ev EventHandler) Run(cd *ClientData, ipmsg *IPMSG) error {
//	cmdstr := cd.Command.String()
//	v := reflect.ValueOf(&ev)
//	method := v.MethodByName(cmdstr)
//	if !method.IsValid() {
//		err := fmt.Errorf("method for Command(%v) not defined", cmdstr)
//		return err
//	}
//	in := []reflect.Value{reflect.ValueOf(cd), reflect.ValueOf(ipmsg)}
//	err := method.Call(in)[0].Interface()
//	// XXX only works if you sure about the return value is always type(error)
//	if err == nil {
//		return nil
//	}
//	return err.(error)
//	//reflect.ValueOf(&ev).MethodByName(cmdstr).Call(in)
//}
//
//func (ev *EventHandler) BR_ENTRY(cd *ClientData, ipmsg *IPMSG) error {
//	ipmsg.SendMSG(cd.Addr, ipmsg.Myinfo(), ANSENTRY)
//	return nil
//}
