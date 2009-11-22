package irc

import "io"
import "os"
import "net"
import "log"
import "bufio"
import "strings"

type RecvFunc	func(*Message)
type ErrorFunc	func(os.Error)

type Conn struct {
	conn		net.Conn;
	recvFunc	RecvFunc;
	errorFunc	ErrorFunc;
	sendq		chan string;
	recvq		chan string;
	error		chan os.Error;
}

func NewConn(netconn net.Conn, recvFunc RecvFunc, errorFunc ErrorFunc) *Conn {
	if netconn == nil {
		return nil
	}
	conn := &Conn{netconn, recvFunc, errorFunc, make(chan string, 64), make(chan string, 64), make(chan os.Error, 0)};
	go conn.run();
	return conn;
}

func (conn *Conn) quit(err os.Error) {
	if err != nil && err != os.EOF {
		log.Stderrf("connection from %s error: %s\n", conn.conn.RemoteAddr(), err);
	}
	if conn.errorFunc != nil {
		conn.errorFunc(err)
	}

	// This will cause writer() to return
	close(conn.sendq);

	// Sends should fail silently
	close(conn.recvq);
	close(conn.error);

	conn.conn.Close();
}

func (conn *Conn) reader() {
	reader := bufio.NewReader(conn.conn);
	for {
		line, err := reader.ReadString('\n');
		if err != nil {
			conn.error <- err;
			return;
		}
		if strings.HasSuffix(line, "\r\n") {
			conn.recvq <- strings.Split(line, "\r\n", 2)[0]
		}
	}
}

func (conn *Conn) writer() {
	for {
		line := <-conn.sendq;

		// Stop if the send queue is closed
		if line == "" && closed(conn.sendq) {
			return
		}

		log.Stderrf("-> %s: %s\n", conn.conn.RemoteAddr(), line);
		_, err := io.WriteString(conn.conn, line + "\r\n");
		if err != nil {
			conn.error <- err;
			return;
		}
	}
}

func (conn *Conn) run() {
	go conn.reader();
	go conn.writer();
	for {
		select {
		case err := <-conn.error:
			conn.quit(err);
			return;

		case line := <-conn.recvq:
			if conn.recvFunc != nil {
				msg := Parse(line);
				log.Stderrf("<- %s: %s\n", conn.conn.RemoteAddr(), msg);
				conn.recvFunc(msg);
			}
		}
	}
}

func (conn *Conn) Quit() {
	conn.error <- nil
}

func (conn *Conn) Send(msg *Message) {
	conn.sendq <- msg.String()
}

func (conn *Conn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}
