package udp

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

// Client ...
type Client struct {
	income     []byte
	outgoing   chan []byte
	reader     *bufio.Reader
	writer     *bufio.Writer
	connection net.Conn
	Remote     net.Addr
	server     *ServerInstance
}

// ServerInstance ...
type ServerInstance struct {
	netListener    net.Listener
	socketListener SocketListener
	clients        map[string]*Client
	delim          byte
}

func (s *ServerInstance) start() error {
	for {
		conn, err := s.netListener.Accept()
		if err != nil {
			s.socketListener.OnError("nil", err)
			break
		} else {
			client := newClient(s, conn)
			client.listen()
			defer client.disconnect()

			s.clients[client.Remote.String()] = client
			s.socketListener.OnConnected("C:" + client.Remote.String())
		}
	}
	return nil
}

// Stop ...
func (s *ServerInstance) Stop() error {
	for _, v := range s.clients {
		v.disconnect()
	}
	err := s.netListener.Close()
	return err
}

// SendToClient ...
func (s *ServerInstance) SendToClient(client *Client, line []byte) {
	client.sendMessage(line)
}

// Broadcast ...
func (s *ServerInstance) Broadcast(line []byte) {
	for _, v := range s.clients {
		v.sendMessage(line)
	}
}

func (client *Client) read() {
	for {
		line, err := client.reader.ReadBytes(client.server.delim)
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
			client.server.socketListener.OnError("C:"+client.Remote.String(), err)
		} else {
			client.server.socketListener.OnPacketReceived("C:"+client.Remote.String(), line)
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
		client.server.socketListener.OnPacketSended("C:"+client.Remote.String(), line)
		_, err := client.writer.Write(line)
		if err != nil {
			client.server.socketListener.OnError("C:"+client.Remote.String(), err)
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
	client.server.socketListener.OnDisconnected("C:" + client.Remote.String())
	client.connection.Close()
	delete(client.server.clients, client.connection.LocalAddr().String())
}

func newClient(serverInstance *ServerInstance, connection net.Conn) *Client {
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
		server:     serverInstance,
	}
	return client
}

// StartTCP start the TCP server.
// it will return ServerInstance instance when no error.
func StartTCP(port int, delimiter byte, listener SocketListener, serverInstance *ServerInstance) {
	netlistener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		listener.OnError("nil", err)
		if netlistener != nil {
			netlistener.Close()
		}
		return
	}
	defer netlistener.Close()
	serverInstance = &ServerInstance{
		netListener:    netlistener,
		socketListener: listener,
		clients:        make(map[string]*Client),
		delim:          delimiter,
	}
	serverInstance.start()
}

// StartUDP start the TCP server.
func StartUDP(port int, delimiter byte, listener SocketListener, serverInstance *ServerInstance) {
	netlistener, err := net.Listen("udp", ":"+strconv.Itoa(port))
	if err != nil {
		listener.OnError("nil", err)
		if netlistener != nil {
			netlistener.Close()
		}
		return
	}
	defer netlistener.Close()
	serverInstance = &ServerInstance{
		netListener:    netlistener,
		socketListener: listener,
		clients:        make(map[string]*Client),
		delim:          delimiter,
	}
	serverInstance.start()
}
