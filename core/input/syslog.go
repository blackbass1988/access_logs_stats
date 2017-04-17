package input

import (
	"bufio"
	"errors"
	"github.com/blackbass1988/access_logs_stats/core/re"
	"io"
	"log"
	"net"
	"sync"
)

const udpSafePackSize = 2048

var (
	//ErrorIncorrectDSN says that given incorrect input DSN
	ErrorIncorrectDSN = errors.New("Incorrect DSN")
	//ErrorUnknownProtocol says that given unknwon protocol. Not TCP or UDP
	ErrorUnknownProtocol = errors.New("Unknown protocol")
)

//SyslogInputReader implements BufferedReader for reading data from syslog.
//It's starts syslog server and receive input data
type SyslogInputReader struct {
	BufferedReader

	m      sync.Mutex
	buffer []byte

	protocol    string
	listen      string
	application string

	acceptor Acceptor

	parser *syslogParser
}

//CreateSyslogInputReader created BufferedReader
//dsn examples::
// syslog:udp:binding_ip:binding_port/application
// syslog:tcp:binding_ip:binding_port/application
// syslog:tcp:binding_ip:binding_port/application
func CreateSyslogInputReader(dsn string) (r *SyslogInputReader, err error) {

	r = &SyslogInputReader{}
	r.buffer = []byte{}
	r.parser, err = newSyslogParser()
	check(err)
	//read dsn
	err = r.parseDsn(dsn)
	if err != nil {
		return
	}

	//create udp or tcp server
	err = r.startServer()

	return
}

func (r *SyslogInputReader) parseDsn(dsn string) (err error) {
	r.protocol, r.listen, r.application, err = parseSyslogDsn(dsn)
	return
}

//ReadToBuffer implements BufferedReader ReadToBuffer method for SyslogInputReader
func (r *SyslogInputReader) ReadToBuffer() {

	if r.protocol == "udp" {
		r.readToBufferUDP()
	} else {
		r.readToBufferTCP()
	}

}

//FlushBuffer implements BufferedReader FlushBuffer method for SyslogInputReader
func (r *SyslogInputReader) FlushBuffer() []byte {
	r.m.Lock()
	buffer := r.buffer
	r.buffer = []byte{}
	r.m.Unlock()
	return buffer
}

//Close implements BufferedReader Close method for SyslogInputReader
func (r *SyslogInputReader) Close() {
	r.acceptor.Close()
}

func (r *SyslogInputReader) readToBufferTCP() {
	// accept new connect and send it to handler
	for {
		conn, err := r.acceptor.Accept()
		check(err)
		go r.handleConnectionTCP(conn)
	}
}

func (r *SyslogInputReader) readToBufferUDP() {
	//create one listener and handle messages in one loop
	conn, err := r.acceptor.Accept()
	check(err)
	r.handleConnectionUDP(conn)
}

func (r *SyslogInputReader) startServer() (err error) {

	var (
		acceptor Acceptor
	)

	if r.protocol == "udp" {
		acceptor, err = r.getUDPAcceptor()
		if err != nil {
			return err
		}
	} else {
		acceptor, err = r.getTCPAcceptor()
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	r.acceptor = acceptor

	return
}

func (r *SyslogInputReader) getUDPAcceptor() (acceptor Acceptor, err error) {

	addr, err := net.ResolveUDPAddr(r.protocol, r.listen)
	if err != nil {
		return acceptor, err
	}
	udpconn, err := net.ListenUDP(r.protocol, addr)
	acceptor = &udpAcceptor{acceptor: udpconn}

	return acceptor, err
}

func (r *SyslogInputReader) getTCPAcceptor() (acceptor Acceptor, err error) {

	addr, err := net.ResolveTCPAddr(r.protocol, r.listen)
	if err != nil {
		return acceptor, err
	}
	tcpl, err := net.ListenTCP(r.protocol, addr)
	acceptor = &tcpAcceptor{acceptor: tcpl}
	return acceptor, err
}

func (r *SyslogInputReader) handleConnectionUDP(conn net.Conn) {
	var (
		read int
		err  error
		b    []byte
	)
	b = make([]byte, udpSafePackSize)
	defer conn.Close()

	for {
		read, err = conn.Read(b)
		if err != nil {
			log.Println(err)
		}
		bytesBuf := b[0:read]
		r.m.Lock()
		if r.appendToBuffer(bytesBuf) {
			r.buffer = append(r.buffer, '\n')
		}
		r.m.Unlock()
		//log.Println(string(r.buffer))
	}
}

func (r *SyslogInputReader) handleConnectionTCP(conn net.Conn) {
	buffer := bufio.NewReader(conn)
	defer conn.Close()

	for {
		bytesBuf, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			r.m.Lock()
			if r.appendToBuffer(bytesBuf) {
				r.buffer = append(r.buffer, '\n')
			}
			r.m.Unlock()
			break
		} else if err != nil {
			check(err)
		}
		r.m.Lock()
		r.appendToBuffer(bytesBuf)
		r.m.Unlock()
	}
	//log.Println(string(r.buffer))
}

func (r *SyslogInputReader) appendToBuffer(byteBuf []byte) bool {
	if len(byteBuf) == 0 {
		return false
	}

	//parse message.
	m, err := r.parser.parseSyslogMsg(string(byteBuf))

	if err == ErrorUnknownInputStringFormat {
		log.Println(ErrorUnknownInputStringFormat, string(byteBuf))
	}
	//Filter by application
	if m.Application != r.application {
		return false
	}

	r.buffer = append(r.buffer, m.Message...)
	return true
}

func parseSyslogDsn(dsn string) (protocol string, listen string, application string, err error) {
	r, err := re.Compile(`(syslog):([a-zA-Z0-9]+):([^/]+)/(\S+)`)
	if err != nil {
		return
	}

	matches := r.FindStringSubmatch(dsn)
	if len(matches) != 5 {
		err = ErrorIncorrectDSN
		return
	}
	protocol = matches[2]
	listen = matches[3]
	application = matches[4]

	if protocol != "udp" && protocol != "tcp" {
		err = ErrorUnknownProtocol
	}

	return
}

//Acceptor is an interface than Can Accept new connection and close it
type Acceptor interface {
	Accept() (net.Conn, error)
	Close() error
}

type tcpAcceptor struct {
	acceptor *net.TCPListener
}

//Accept accepts new tcp connection
func (l *tcpAcceptor) Accept() (net.Conn, error) {
	return l.acceptor.Accept()
}

//Close closes tcp connection
func (l *tcpAcceptor) Close() error {
	return l.acceptor.Close()
}

type udpAcceptor struct {
	acceptor *net.UDPConn
}

//Accept accepts new udp packet
func (l *udpAcceptor) Accept() (net.Conn, error) {
	return l.acceptor, nil
}

//Close closes udp file descriptor
func (l *udpAcceptor) Close() error {
	return l.acceptor.Close()
}
