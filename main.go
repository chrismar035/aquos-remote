package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"bytes"
)

func main() {
	var login string
	var password string
	var ip string
	var port string
	flag.StringVar(&login, "login", "tv", "Login for your TV")
	flag.StringVar(&password, "password", "teevee", "Password for your TV")
	flag.StringVar(&ip, "ip", "10.0.1.220", "IP address of your TV")
	flag.StringVar(&port, "port", "10002", "Port number of your TV")
	flag.Parse()

	tvAddr, err := net.ResolveTCPAddr("tcp", ip + ":" + port)
	if err != nil {
		fmt.Println("Could not resolve ", ip, "on", port)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tvAddr)
	if err != nil {
		fmt.Println("Could not connect to ", ip, "on", port)
		os.Exit(1)
	}
	defer conn.Close()

	conn.SetReadBuffer(1024)
	conn.SetWriteBuffer(1024)

	var recv []byte
	conn.Read(recv)
	conn.Write([]byte(login + "\r"))
	conn.Read(recv)
	conn.Write([]byte(password + "\r"))
	conn.Read(recv)
	if bytes.Equal(recv, []byte("")) {
		fmt.Println("Logged in")
	} else {
		fmt.Println("Couldn't log in")
		os.Exit(1)
	}

	conn.Write([]byte("POWR0   \r"))
}
