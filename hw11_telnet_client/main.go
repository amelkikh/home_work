package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "Connect timeout duration, eg. --timeout=10s")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Please provide string with server and port arguments, eg. 127.0.0.1 8080")
		return
	}

	if _, err := strconv.Atoi(args[1]); err != nil {
		fmt.Println("Invalid port number:", args[1])
		return
	}

	addr := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Println("Connect error:", addr, err)
		return
	}

	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		for {
			err := client.Receive()
			switch err { //nolint:errorlint
			case nil:
			default:
				stop()
				return
			}
		}
	}()

	go func() {
		for {
			err := client.Send()
			switch err { //nolint:errorlint
			case nil:
			default:
				stop()
				return
			}
		}
	}()

	<-ctx.Done()
}
