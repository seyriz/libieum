package udp

import "net"

// SocketListener ...
type SocketListener interface {
	OnPacketReceived(addr string, packet []byte)
	OnPacketSended(addr string, packet []byte)
	OnError(addr string, err error)
	OnConnected(addr string)
	OnDisconnected(addr string)
}

// SocketInterface ...
type SocketInterface interface {
	Start(async bool) error
	Stop() error
	Send(line []byte) error
	BroadCast(line []byte)
	GetClientList() map[net.Addr]*Client
}
