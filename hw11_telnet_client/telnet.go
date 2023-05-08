package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Cli struct {
	address   string
	timeout   time.Duration
	in        io.ReadCloser
	out       io.Writer
	session   net.Conn
	inBuffer  []byte
	outBuffer []byte
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	instance := Cli{address: address, timeout: timeout, in: in, out: out, session: nil}
	if instance.timeout == 0 {
		instance.timeout = 10 * time.Second
	}
	instance.inBuffer = make([]byte, 1024, 2048)
	instance.outBuffer = make([]byte, 1024, 2048)
	return instance
}

func (self Cli) Connect() error {
	var err error
	if self.session, err = net.DialTimeout("TCP", self.address, self.timeout); err != nil {
		return err
	}
	return nil
}

func (self Cli) Send() error {
	n, err := self.in.Read(self.outBuffer)
	if err != nil {
		return err
	}
	if _, err = self.session.Write(self.outBuffer[:n]); err != nil {
		return err
	}
	return nil
}

func (self Cli) Receive() error {
	n, err := self.session.Read(self.inBuffer)
	if err != nil {
		return err
	}
	fmt.Print(self.inBuffer[:n])
	return nil
}

func (self Cli) Close() error {
	if err := self.session.Close(); err != nil {
		return err
	}
	return nil
}
