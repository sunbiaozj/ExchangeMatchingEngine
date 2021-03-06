package main

import (
	"fmt"
	"os"

	"github.com/farice/EME/redis"
	log "github.com/sirupsen/logrus"
)

type StdOutHook struct{}

func (h *StdOutHook) Levels() []log.Level {
	return []log.Level{
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}

var fmter = new(log.TextFormatter)

func (h *StdOutHook) Fire(entry *log.Entry) (err error) {
	line, err := fmter.Format(entry)
	if err == nil {
		fmt.Fprintf(os.Stderr, string(line))
	}
	return
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("/var/log/erss/exchange.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stdout")
	}

	log.AddHook(&StdOutHook{})

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)

	// MARK: - Global Transaction ID Counter

	ex, err := redis.Exists("TransactionCounter")
	if !ex {
		redis.Set("TransactionCounter", 0)
	}

}

func main() {
	var clientCount = 0

	// Set up logging for performance benchmarking
	CreateBenchmarkingLog()

	// addr: exchange, port: 12345
	server := NewTCPServer("exchange:12345")

	// MARK: - Implement new client, message, and closed connection callbacks

	server.OnNewConnection(func(c *Connection) {
		// New Client Connected
		log.WithFields(log.Fields{
			"client count": clientCount,
		}).Info("New client connection")
		clientCount += 1

		//c.Send("Hello\n")

	})

	server.OnNewMessage(func(c *Connection, message []byte) {
		c.handleRequest(message)

	})

	server.OnClientConnectionClosed(func(c *Connection, err error) {
		// Lost connection with Client
		log.WithFields(log.Fields{
			"error": err,
		}).Info("Connection closed")

	})

	server.Listen()
}
