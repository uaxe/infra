package telnet

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	ctx       context.Context
	conn      net.Conn
	writer    chan []byte
	writerRet chan result
	reader    chan any
	readerRet chan result
	opt       option
}

type result struct {
	n   int
	err error
	raw []byte
}

type option struct {
	ip       string
	port     string
	auth     bool
	userName string
	passWord string
}

type OptionFunc func(opt *option)

func WithIP(ip string) OptionFunc {
	return func(opt *option) {
		opt.ip = ip
	}
}

func WithPort(port string) OptionFunc {
	return func(opt *option) {
		opt.port = port
	}
}

func WithAuth(auth bool) OptionFunc {
	return func(opt *option) {
		opt.auth = auth
	}
}

func WithUserName(userName string) OptionFunc {
	return func(opt *option) {
		opt.userName = userName
	}
}

func WithPassWord(passWord string) OptionFunc {
	return func(opt *option) {
		opt.passWord = passWord
	}
}

func defaultOption() option {
	return option{
		ip:   "127.0.0.1",
		port: "23",
		auth: false,
	}
}

func NewClient(ctx context.Context, opts ...OptionFunc) *Client {
	c := Client{
		ctx:       ctx,
		writer:    make(chan []byte, 1),
		writerRet: make(chan result, 1),
		reader:    make(chan any, 1),
		readerRet: make(chan result, 1),
		opt:       defaultOption(),
	}
	for _, opt := range opts {
		opt(&c.opt)
	}
	return &c
}

func (c *Client) Write(b []byte) error {
	c.writer <- b
	w := <-c.writerRet
	return w.err
}

func (c *Client) Read() ([]byte, error) {
	c.reader <- nil
	r := <-c.readerRet
	return r.raw, r.err
}

func (c *Client) _read() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-c.reader:
				b := make([]byte, 1024)
				n, err := c.conn.Read(b)
				c.readerRet <- result{n, err, b}
			}
		}
	}()
}

func (c *Client) _write() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case cmd := <-c.writer:
			cmd = append(cmd, '\n')
			n, err := c.conn.Write(cmd)
			time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
			c.writerRet <- result{n, err, cmd}
		}
	}
}

func (c *Client) handshake() error {
	var buf [4096]byte
	n, err := c.conn.Read(buf[0:])
	if err != nil {
		return err
	}

	buf[1] = 252
	buf[4] = 252
	buf[7] = 252
	buf[10] = 252

	n, err = c.conn.Write(buf[0:n])
	if err != nil {
		return err
	}

	n, err = c.conn.Read(buf[0:])
	if err != nil {
		return err
	}

	if !c.opt.auth {
		return nil
	}

	n, err = c.conn.Write([]byte(c.opt.userName + "\n"))
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
	n, err = c.conn.Read(buf[0:])
	if err != nil {
		return err
	}

	n, err = c.conn.Write([]byte(c.opt.passWord + "\n"))
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
	n, err = c.conn.Read(buf[0:])
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Run() error {
	addr := fmt.Sprintf("%s:%s", c.opt.ip, c.opt.port)
	conn, err := net.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		return err
	}
	c.conn = conn

	if err = c.handshake(); err != nil {
		return err
	}

	c._read()

	go func() { c._write() }()

	return nil
}

func (c *Client) Stdin() error {
	addr := fmt.Sprintf("%s:%s", c.opt.ip, c.opt.port)
	conn, err := net.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		return err
	}
	c.conn = conn

	if err = c.handshake(); err != nil {
		return err
	}

	c._read()

	go func() { c._write() }()

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			if err = c.Write([]byte(readLine())); err != nil {
				fmt.Println(err)
				break
			}
			ret, err := c.Read()
			if err != nil {
				fmt.Println(err)
				break
			} else {
				fmt.Println(string(ret))
			}
		}
	}
}

func (c *Client) Close() error {
	close(c.writer)
	return c.conn.Close()
}

func readLine() string {
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return strings.TrimSpace(line)
}

const (
	TIME_DELAY_AFTER_WRITE = 500
)
