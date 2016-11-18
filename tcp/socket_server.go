package tcp

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

var clients map[net.Addr]*Client
var tcpListener TCPSocketServerListener
var socketListener net.Listener

type Client struct {
	income     []byte
	outgoing   chan []byte
	reader     *bufio.Reader
	writer     *bufio.Writer
	connection net.Conn
	Remote     net.Addr
	delim      byte
}

type TCPSocketServerListener interface {
	OnPacketReceived(client *Client, packet []byte)
	OnPacketSended(client *Client, packet []byte)
	OnError(client *Client, err error)
	OnConnected(client *Client)
	OnDisconnected(client *Client)
}

func (client *Client) Read() {
	for {
		line, err := client.reader.ReadBytes(client.delim)
		if client == nil {
			break
		}
		if err != nil {
			if err.Error() == "EOF" {
				client.Disconnect()
				return
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			tcpListener.OnError(client, err)
		} else {
			tcpListener.OnPacketReceived(client, line)
			client.income = line
		}
	}
}

func (client *Client) Write() {
	for {
		if client == nil {
			break
		}
		line := <-client.outgoing
		tcpListener.OnPacketSended(client, line)
		_, err := client.writer.Write(line)
		if err != nil {
			tcpListener.OnError(client, err)
		}
		client.writer.Flush()
	}
}

func (client *Client) SendMessage(line []byte) {
	client.outgoing <- line
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func (client *Client) Disconnect() {
	tcpListener.OnDisconnected(client)
	client.connection.Close()
	delete(clients, client.connection.LocalAddr())
}

func BroadCast(line []byte) {
	for _, v := range clients {
		v.outgoing <- line
	}
}

func GetClientList() map[net.Addr]*Client {
	return clients
}

func NewClient(connection net.Conn, delim byte) *Client {
	// log.Println("Client : ", connection.RemoteAddr())
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)
	client := &Client{
		income:     make([]byte, 0),
		outgoing:   make(chan []byte, 0),
		reader:     reader,
		writer:     writer,
		connection: connection,
		delim:      delim,
		Remote:     connection.RemoteAddr(),
	}
	return client
}

func ListenSocket(port int, delim byte, listener TCPSocketServerListener) {
	clients = make(map[net.Addr]*Client)
	tcpListener = listener
	var err error
	socketListener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		tcpListener.OnError(nil, err)
		return
	}
	// log.Println("Listen in : ", port)
	defer socketListener.Close()

	for {
		conn, err := socketListener.Accept()
		if err != nil {
			tcpListener.OnError(nil, err)
			break
		} else {
			client := NewClient(conn, delim)
			client.Listen()
			defer client.Disconnect()

			clients[client.connection.LocalAddr()] = client
			tcpListener.OnConnected(client)
		}
	}
}

func CloseSocket() {
	if socketListener != nil {
		socketListener.Close()
	}
}
