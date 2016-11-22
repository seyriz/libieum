package udp

import (
	"bufio"
	"errors"
	"net"
	"strconv"
	"strings"
)

// ConnectedServer ...
type ConnectedServer struct {
	income         []byte
	reader         *bufio.Reader
	writer         *bufio.Writer
	connection     net.Conn
	Remote         net.Addr
	delim          byte
	socketListener SocketListener
}

func (s *ConnectedServer) connect() {
	go s.read()
}

func (s *ConnectedServer) disconnect() {
	s.socketListener.OnDisconnected("S:" + s.Remote.String())
	s.connection.Close()
}

func (s *ConnectedServer) read() {
	for {
		if s == nil {
			break
		}
		line, err := s.reader.ReadBytes(s.delim)
		if err != nil {
			if err.Error() == "EOF" {
				s.disconnect()
				return
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			s.socketListener.OnError("S:"+s.Remote.String(), err)
		} else {
			s.socketListener.OnPacketReceived("S:"+s.Remote.String(), line)
			s.income = line
		}
	}
}

// SendMessageToServer ...
func (s *ConnectedServer) SendMessageToServer(line []byte) error {
	if s == nil {
		return errors.New("Closed connection")
	}
	_, err := s.writer.Write(line)
	if err != nil {
		return err
	}
	s.writer.Flush()
	return nil
}

// ConnectTCP ...
func ConnectTCP(raddr string, port int, delim byte, listner SocketListener) (*ConnectedServer, error) {
	connection, err := net.Dial("tcp", raddr+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)
	connectedServer := &ConnectedServer{
		income:         make([]byte, 0),
		reader:         reader,
		writer:         writer,
		connection:     connection,
		Remote:         connection.RemoteAddr(),
		delim:          delim,
		socketListener: listner,
	}
	return connectedServer, nil
}
