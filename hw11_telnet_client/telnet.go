package main

import (
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
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	session net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	instance := Cli{address: address, timeout: timeout, in: in, out: out, session: nil}
	return &instance
}

func (cli *Cli) Connect() error {
	var err error
	cli.session, err = net.DialTimeout("tcp", cli.address, cli.timeout)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Cli) Send() error {
	data, err := io.ReadAll(cli.in)
	if err != nil {
		return err
	}
	if _, err = cli.session.Write(data); err != nil {
		return err
	}
	return nil
}

func (cli *Cli) Receive() error {
	var err error
	var n int
	inBuffer := make([]byte, 0, 4096)
	tmp := make([]byte, 256)
	for err != io.EOF {
		n, err = cli.session.Read(tmp)
		if err != nil && err != io.EOF {
			return err
		}
		inBuffer = append(inBuffer, tmp[:n]...)
	}
	_, err = cli.out.Write(inBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Cli) Close() error {
	if err := cli.session.Close(); err != nil {
		return err
	}
	return nil
}
