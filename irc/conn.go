package irc

import "io"
import "net"
import "bufio"
import "strings"

type Handler interface {
	ProcessMessage(*Conn, *Message);
}

type Conn struct {
	conn	net.Conn;
	handler	Handler;
	sendq	chan string;
	recvq	chan string;
	quit	chan bool;
}

func NewConn(netconn net.Conn, handler Handler) *Conn {
	if netconn == nil {
		return nil
	}
	return &Conn{netconn, handler, make(chan string, 64), make(chan string, 64), make(chan bool, 1)}
}

func (conn *Conn) close() {
	if conn.conn != nil {
		conn.conn.Close();
		conn.conn = nil;
	}
}

func (conn *Conn) reader() {
	reader := bufio.NewReader(conn.conn);
	for {
		line, err := reader.ReadString('\n');
		if strings.HasSuffix(line, "\r\n") {
			conn.recvq <- strings.Split(line, "\r\n", 2)[0];
		}
		if err != nil {
			// TODO handle error
			return
		}
	}
}

func (conn *Conn) Run() {
	go conn.reader();
	for {
		select {
		case <-conn.quit:
			conn.close();
			return;

		case line := <-conn.sendq:
			_, err := io.WriteString(conn.conn, line + "\r\n");
			if err != nil {
				// TODO handle error
				return
			}

		case line := <-conn.recvq:
			msg := Parse(line);
			if msg != nil && conn.handler != nil {
				conn.handler.ProcessMessage(conn, msg)
			}
		}
	}
}

func (conn *Conn) Quit() {
	conn.quit <- true
}

func (conn *Conn) Send(msg *Message) {
	conn.sendq <- msg.String()
}
