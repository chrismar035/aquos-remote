package main

import (
	"net"
	"os"
	"bytes"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Aquos Remote"
	app.Usage = "Control your Aquos TV like magic"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "login",
			Value: "tv",
			Usage: "login for your TV",
		},
		cli.StringFlag{
			Name: "password",
			Value: "teevee",
			Usage: "password for your TV",
		},
		cli.StringFlag{
			Name: "ip",
			Value: "10.0.1.220",
			Usage: "IP address of your TV",
		},
		cli.StringFlag{
			Name: "port",
			Value: "10002",
			Usage: "port number of your TV",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "power-on",
			Usage: "turn the TV on",
			Action: func(c *cli.Context) {
				conn := login(c)
				defer conn.Close()
				conn.Write([]byte("POWR1   \r"))
			},
		},
		{
			Name: "power-off",
			Usage: "turn the TV off",
			Action: func(c *cli.Context) {
				conn := login(c)
				defer conn.Close()
				conn.Write([]byte("POWR0   \r"))
			},
		},
		{
			Name: "input",
			Usage: "select an input\n1 - HDMI 1\n2 - HDMI 2\n3 - HDMI 3\n4 - HDMI 4\n5 - Component\n6 - Video 1\n7 - Video 2",
			Action: func(c *cli.Context) {
				conn := login(c)
				defer conn.Close()
				conn.Write([]byte("IAVD" + c.Args()[0] + "   \r"))
			},
		},
		{
			Name:  "volume",
			Usage: "select a volume level (between 0 and 60)",
			Action: func(c *cli.Context) {
				conn := login(c)
				defer conn.Close()
				conn.Write([]byte("VOLM" + c.Args()[0] + "  \r"))
			},
		},
	}

	app.Run(os.Args)
}

func login(c *cli.Context) (conn *net.TCPConn) {
	ip, port := c.GlobalString("ip"), c.GlobalString("port")
	tvAddr, err := net.ResolveTCPAddr("tcp", ip + ":" + port)
	if err != nil {
		println("Could not resolve", ip, "on", port)
		os.Exit(1)
	}

	conn, err = net.DialTCP("tcp", nil, tvAddr)
	if err != nil {
		println("Could not connect to", ip, "on", port)
		os.Exit(1)
	}

	conn.SetReadBuffer(1024)
	conn.SetWriteBuffer(1024)

	login, password := c.GlobalString("login"), c.GlobalString("password")

	var recv []byte
	conn.Read(recv)
	conn.Write([]byte(login + "\r"))
	conn.Read(recv)
	conn.Write([]byte(password + "\r"))
	conn.Read(recv)
	if !bytes.Equal(recv, []byte("")) {
		println("Couldn't log in")
		os.Exit(1)
	}

	return conn
}
