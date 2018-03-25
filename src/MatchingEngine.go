package main

import (
	"./tcp_server"
	"github.com/golang/glog"
	"os"
	"flag"
	"fmt"
	)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n", )
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line 
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()
}

func main() {
	server := tcp_server.New("127.0.0.1:12345")

	server.OnNewClient(func(c *tcp_server.Client) {
		// New Client Connected
		glog.Info("New client connection")
		c.Send("Hello")
		glog.Flush()
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// New Message Received
		glog.Info("New message received")
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// Lost connection with Client
		glog.Info("New client connection")
	})

	server.Listen()
}