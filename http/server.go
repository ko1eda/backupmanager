package http

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ko1eda/backupmanager/log"
	"github.com/ko1eda/backupmanager/wasabi"
)

// Server represents an http server
type Server struct {
	S3Service  *wasabi.S3Service
	IAMService *wasabi.IAMService
	Logger     *log.Logger
	Router     *http.ServeMux
	Address    string
	listener   net.Listener
	// Hasrer
}

// NewServer returns a new sever instance
func NewServer(addr string, s3s *wasabi.S3Service, iam *wasabi.IAMService, logger *log.Logger) *Server {
	return &Server{
		S3Service:  s3s,
		IAMService: iam,
		Logger:     logger,
		Address:    addr,
	}
}

// Open opens the server and listens at the specifed address
func (s *Server) Open() error {
	// Open socket.
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Address))

	if err != nil {
		return err
	}

	s.listener = ln

	// Start HTTP server.
	go func() { http.Serve(s.listener, s.router()) }()

	s.Logger.Log("Server started listening on port " + s.Address + "....")

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.listener != nil {
		s.listener.Close()
	}

	return nil
}

// Router creates a new servermux and registers all routes to it
func (s *Server) router() http.Handler {
	s.Router = http.NewServeMux()

	wh := newWasabiHandler(s.S3Service, s.IAMService, s.Logger)

	s.Router.Handle("/cloud/backups/infrastructure/create", wh.handleCreateBackupInfrastructure())

	return s.Router
}
