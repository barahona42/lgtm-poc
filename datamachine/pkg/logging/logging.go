package logging

import (
	"io"
	"log/slog"
	"net/http"
	"os"
)

var (
	buildtime string = "NOTSET"
)

func ConfigureSlog(file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	l := slog.
		New(slog.NewJSONHandler(io.MultiWriter(os.Stderr, fd), &slog.HandlerOptions{AddSource: true})).
		With("version", buildtime)
	slog.SetDefault(l)
	return nil
}

func HTTPRequestLogger(r *http.Request) *slog.Logger {
	return slog.Default().
		WithGroup("request").
		With("uri", r.RequestURI, "remote", r.RemoteAddr, "method", r.Method)
}
