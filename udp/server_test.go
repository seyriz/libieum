package udp

import "testing"

type SockListener struct{}

var tests *testing.T

func TestTCP(t *testing.T) {
	tests = t
	listener := SockListener{}
	var serverInstance *ServerInstance
	var svI chan *ServerInstance
	svI = make(chan *ServerInstance)
	go StartTCP(11234, 0x0A, listener, serverInstance)
	svI <- serverInstance
	sv, err := ConnectTCP("localhost", 11234, 0x0A, listener)
	t.Log("asd1")
	if err != nil {
		t.Log(err)
	}
	t.Log("asd2")
	sv.SendMessageToServer([]byte("ASDFE"))

	t.Log("asd3")
	err = serverInstance.Stop()
	t.Log("asd4")
	if err != nil {
		t.Log(err)
	}
}

func (l SockListener) OnPacketReceived(addr string, packet []byte) {
	tests.Log("OnPacketReceived", addr, string(packet))
}
func (l SockListener) OnPacketSended(addr string, packet []byte) {
	tests.Log("OnPacketSended", addr, string(packet))
}
func (l SockListener) OnError(addr string, err error) {
	tests.Log("OnError", addr, err)
	tests.Fail()
}
func (l SockListener) OnConnected(addr string) {
	tests.Log("OnConnected", addr)
}
func (l SockListener) OnDisconnected(addr string) {
	tests.Log("OnDisconnected", addr)
}
