package main

import (
	"datamachine/pkg/configuration"
	"datamachine/pkg/logging"
	"datamachine/pkg/messages"
	"fmt"
	"log/slog"
	"net"
	"sync"
)

func main() {
	config, err := configuration.Load()
	if err != nil {
		panic(err)
	}
	if err := logging.ConfigureSlog(fmt.Sprintf("%s/%s", config.LogDir, "collector.log")); err != nil {
		panic(err)
	}
	slog.With("logdir", config.LogDir).Info("logger configured")
	udpaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", config.Collector.Addr, config.Collector.Port))
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			var (
				b   []byte = make([]byte, messages.MessageSize)
				oob []byte = make([]byte, messages.MessageSize)
			)
			_, _, _, _, err := conn.ReadMsgUDP(b, oob)
			if err != nil {
				slog.With("error", err).Error("failed to read udp message")
				continue
			}
			msg := messages.Message{}
			if err := msg.UDPDecode(b); err != nil {
				slog.With("error", err).Error("failed to decode message")
			}
			slog.Info(fmt.Sprintf("%v", msg))
		}
	}()
	// wg.Add(1)

	wg.Wait()
}
