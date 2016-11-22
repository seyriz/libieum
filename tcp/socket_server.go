package tcp

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

var clients map[net.Addr]*Client
var socketListener TCPSocketListener
var serverInstance server

type Client struct {
	income     []byte
	outgoing   chan []byte
	reader     *bufio.Reader
	writer     *bufio.Writer
	connection net.Conn
	Remote     net.Addr
}

type server struct {
	port        int
	delim       byte
	async       bool
	listener    TCPSocketListener
	controller  *TCPSocketInterface
	netListener net.Listener
}

func (t *server) Start(async bool) error {

}

func (t *server) Stop() error {

}

func (t *server) Send(line []byte) error {

}

func (t *server) GetClientList() map[net.Addr]*Client {

}

func (t *server) BroadCast(line []byte) {

}

func (client *Client) read() {
	for {
		line, err := client.reader.ReadBytes(serverInstance.delim)
		if client == nil {
			break
		}
		if err != nil {
			if err.Error() == "EOF" {
				client.disconnect()
				return
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			socketListener.OnError(client, err)
		} else {
			socketListener.OnPacketReceived(client, line)
			client.income = line
		}
	}
}

func (client *Client) write() {
	for {
		if client == nil {
			break
		}
		line := <-client.outgoing
		socketListener.OnPacketSended(client, line)
		_, err := client.writer.Write(line)
		if err != nil {
			socketListener.OnError(client, err)
		}
		client.writer.Flush()
	}
}

func (client *Client) sendMessage(line []byte) {
	client.outgoing <- line
}

func (client *Client) listen() {
	go client.read()
	go client.write()
}

func (client *Client) disconnect() {
	socketListener.OnDisconnected(client)
	client.connection.Close()
	delete(clients, client.connection.LocalAddr())
}

func broadCast(line []byte) {
	for _, v := range clients {
		v.outgoing <- line
	}
}

func getClientList() map[net.Addr]*Client {
	return clients
}

func newClient(connection net.Conn) *Client {
	// log.Println("Client : ", connection.RemoteAddr())
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)
	client := &Client{
		income:     make([]byte, 0),
		outgoing:   make(chan []byte, 0),
		reader:     reader,
		writer:     writer,
		connection: connection,
		Remote:     connection.RemoteAddr(),
	}
	return client
}

func listenSocket(svr server) error {

	for {
		conn, err := netListener.Accept()
		if err != nil {
			socketListener.OnError(nil, err)
			break
		} else {
			client := newClient(conn, delim)
			client.listen()
			defer client.disconnect()

			clients[client.connection.LocalAddr()] = client
			socketListener.OnConnected(client)
		}
	}
	return nil
}

func listenSocketAsync(svr server) error {
	return nil
}

func closeSocket() {
	if netListener != nil {
		netListener.Close()
	}
}

func Init(portP int, delimP byte, listenerP TCPSocketListener, controllerP *TCPSocketInterface) {
	controllerP := &TCPSocketInterface{}
	serverConf := server{
		port:       portP,
		delim:      delimP,
		listener:   listenerP,
		controller: controllerP,
	}

	clients = make(map[net.Addr]*Client)
	socketListener = listener
	var err error
	netListener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		socketListener.OnError(nil, err)
		if netListener != nil {
			netListener.Close()
		}
		return err
	}
	defer netListener.Close()
}
