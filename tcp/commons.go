package tcp

import "net"

func AppendPacketLine(slice []byte, elements []byte) []byte {
	n := len(slice)
	total := len(slice) + len(elements)
	if total > cap(slice) {
		// Reallocate. Grow to 1.5 times the new size, so we can still grow.
		newSize := total*3/2 + 1
		newSlice := make([]byte, total, newSize)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[:total]
	copy(slice[n:], elements)
	return slice
}

type TCPSocketListener interface {
	//
	OnPacketReceived(client *Client, packet []byte)
	OnPacketSended(client *Client, packet []byte)
	OnError(client *Client, err error)
	OnConnected(client *Client)
	OnDisconnected(client *Client)
}

type TCPSocketInterface interface {
	Start(async bool) error
	Stop() error
	Send(line []byte) error
	GetClientList() map[net.Addr]*Client
	BroadCast(line []byte)
}
