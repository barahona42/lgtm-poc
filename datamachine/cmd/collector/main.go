package main

import (
	"datamachine/pkg/configuration"
	"datamachine/pkg/logging"
	"datamachine/pkg/messages"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
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
	go func() { // udp server
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
	wg.Add(1)
	go func() { // http server
		defer wg.Done()
		srvmux := http.NewServeMux()
		// attach healthcheck endpoint
		srvmux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			log := logging.HTTPRequestLogger(r)
			switch r.Method {
			case http.MethodGet:
				log.Info("received healthcheck")
				w.WriteHeader(http.StatusOK)
			default:
				log.Info("invalid healthcheck received")
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})
		slog.Info(fmt.Sprintf("starting http server at %s:%s", config.Collector.Addr, config.Collector.HealthcheckPort))
		srv := &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.Collector.Addr, config.Collector.HealthcheckPort),
			Handler: srvmux}
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			slog.Info("server closed")
		} else if err != nil {
			slog.With("error", err).Error("error with http server")
		}
	}()
	wg.Wait()
}
