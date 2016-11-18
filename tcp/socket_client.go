package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	write      *bufio.Writer
	read       *bufio.Reader
	income     []byte
	connection net.Conn
}

func (server *Server) Read() {
	for {
		line, isPrefix, err := server.read.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed")
				server.connection.Close()
			}
		} else {
			fmt.Println("Line : ", line)
			if isPrefix {
				server.income = AppendPacketLine(server.income, line)
			} else {
				server.income = line
			}
		}
	}
}

func (server *Server) Write(output []byte) {
	nn, err := server.write.Write(output)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Write: ", nn)
	}
}

func Connect(server string, port int) Server {
	connection, err := net.Dial("tcp", server+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	return Server{
		write:      bufio.NewWriter(connection),
		read:       bufio.NewReader(connection),
		income:     make([]byte, 0),
		connection: connection,
	}
}
