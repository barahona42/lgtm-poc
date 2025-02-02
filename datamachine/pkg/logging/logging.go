package logging

import (
	"io"
	"log/slog"
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
