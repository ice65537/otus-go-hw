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
	address  string
	timeout  time.Duration
	in       io.ReadCloser
	out      io.Writer
	session  net.Conn
	inBuffer []byte
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	instance := Cli{
		address:  address,
		timeout:  timeout,
		in:       in,
		out:      out,
		session:  nil,
		inBuffer: make([]byte, 4096),
	}
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
	outBuffer, err := io.ReadAll(cli.in)
	if err != nil {
		return err
	}
	if _, err := cli.session.Write(outBuffer); err != nil {
		return err
	}
	return nil
}

func (cli *Cli) Receive() error {
	var err error
	var n int
	n, err = cli.session.Read(cli.inBuffer)
	if err != nil {
		return err
	}
	_, err = cli.out.Write(cli.inBuffer[:n])
	if err != nil {
		return err
	}
	return nil
}

func (cli *Cli) Close() error {
	if cli.session == nil {
		return nil
	}
	if err := cli.session.Close(); err != nil {
		return err
	}
	return nil
}
