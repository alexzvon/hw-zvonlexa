package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "таймаут на подключение к серверу")

	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		log.Fatalln("Неверный вызов\n\nПримеры вызовов:\n$ go-telnet --timeout=10s host port\n$ go-telnet mysite.ru 8080\n$ go-telnet --timeout=3s 1.1.1.1 123")
	}

	var address string

	address = net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		defer cancel()

		if err := client.Send(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer cancel()

		if err := client.Receive(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-ctx.Done()
}
