package main

import (
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.Reader
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.Reader, out io.Writer) TelnetClient {
	return &Client{
		in:      in,
		out:     out,
		address: address,
		timeout: timeout,
	}
}

func (c *Client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return err
	}

	byteLine := []byte("...Connection was closed by peer\n")
	n, err := os.Stderr.Write(byteLine)
	if err != nil {
		return err
	}
	if n != len(byteLine) {
		return errors.New("the app could not send ...EOF to Stderr")
	}

	return nil
}

func (c *Client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return err
	}

	byteLine := []byte("...EOF\n")
	n, err := os.Stderr.Write(byteLine)
	if err != nil {
		return err
	}
	if n != len(byteLine) {
		return errors.New("the app could not send ...EOF to Stderr")
	}

	return nil
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.conn = conn

	byteLine := []byte(concat("...Connection to ", c.address, "\n"))
	n, err := os.Stderr.Write(byteLine)
	if err != nil {
		return err
	}
	if n != len(byteLine) {
		return errors.New("the app could not send ...Connected to Stderr")
	}

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func concat(s ...string) string {
	var builder strings.Builder
	var lenght int

	for _, v := range s {
		lenght += len(v)
	}

	builder.Grow(lenght)

	for _, v := range s {
		builder.WriteString(v)
	}

	return builder.String()
}
