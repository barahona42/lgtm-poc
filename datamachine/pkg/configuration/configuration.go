package configuration

import (
	"os"
)

type Server struct {
	Addr string
	Port string
}

// Configuration: global configuration struct
type Configuration struct {
	LogDir    string // where the logfile should be created
	Collector Server
	Generator Server
}

var (
	c *Configuration = new(Configuration)
)

func Load() (Configuration, error) {
	c.LogDir = os.Getenv("LOG_DIR")
	c.Collector = Server{Addr: os.Getenv("COLLECTOR_ADDR"), Port: os.Getenv("COLLECTOR_PORT")}
	c.Generator = Server{Addr: os.Getenv("GENERATOR_ADDR"), Port: os.Getenv("GENERATOR_PORT")}
	return *c, nil
}
