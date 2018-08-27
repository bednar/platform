package nats

import (
	stand "github.com/nats-io/nats-streaming-server/server"
)

const ServerName = "platform"

// Server wraps a connection to a NATS streaming server
type Server struct {
	Server *stand.StanServer
	config Config
}

// Open starts a NATS streaming server
func (s *Server) Open() error {
	opts := stand.GetDefaultOptions()
	opts.StoreType = stores.TypeFile
	opts.ID = ServerName
	opts.FilestoreDir = s.config.FilestoreDir
	server, err := stand.RunServerWithOpts(opts, nil)
	if err != nil {
		return err
	}

	s.Server = server

	return nil
}

// Config is the configuration for the NATS streaming server
type Config struct {
	// The directory where nats persists message information
	FilestoreDir string
}

// NewServer creates and returns a new server struct from the provided config
func NewServer(c Config) *Server {
	return &Server{
		config: Config
	}
}
