package main

import (
	"crypto/rand"
	"datamachine/pkg/configuration"
	"datamachine/pkg/logging"
	"datamachine/pkg/messages"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
	"time"
)

var (
	TickDuration time.Duration = 5 * time.Second
	RunDuration  time.Duration = 30 * time.Second
)

func GenerateMessage(t time.Time) (*messages.Message, error) {
	m := messages.Message{Time: t.UnixMilli(), Data: [16]byte{}}
	// generate random data
	if _, err := rand.Read(m.Data[:]); err != nil {
		return nil, err
	}
	// generate a random value
	n, err := rand.Int(rand.Reader, big.NewInt(1e6))
	if err != nil {
		return nil, err
	}
	m.Value = n.Int64()
	return &m, nil
}

func main() {
	config, err := configuration.Load()
	if err != nil {
		panic(err)
	}
	if err := logging.ConfigureSlog(fmt.Sprintf("%s/%s", config.LogDir, "generator.log")); err != nil {
		panic(err)
	}
	slog.With("logdir", config.LogDir).Info("logger configured")

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", config.Collector.Addr, config.Collector.Port))
	if err != nil {
		slog.With("error", err).Error("failed to set up udp connection")
		os.Exit(1)
	}
	defer conn.Close()

	slog.Info("starting ticker")
	tick := time.NewTicker(TickDuration)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				slog.Info("received done signal")
				return
			case t := <-tick.C:
				slog.With("time", t.Format(time.RFC3339)).Info("received a tick")
				msg, err := GenerateMessage(t)
				if err != nil {
					slog.With("error", err).Error("failed to generate message")
					continue
				}
				b, err := msg.UDPEncode()
				if err != nil {
					slog.With("error", err).Error("failed to encode message")
					continue
				}
				n, err := conn.Write(b)
				if err != nil {
					slog.With("error", err).Error("failed to write message to connection")
					continue
				}
				slog.With("conn_write_size", n).Info("wrote message to connection")
			}
		}
	}()
	time.Sleep(RunDuration)
	tick.Stop()
	done <- true
	slog.Info("ticker stopped")
}
