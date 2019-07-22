package http

import (
	"log"
	"net"
	"net/http"

	"github.com/ko1eda/backupmanager/wasabi"
)

// Server represents an http server
type Server struct {
	S3Service  *wasabi.S3Service
	IAMService *wasabi.IAMService
	Router     *http.ServeMux
	listener   net.Listener
	Address    string
	// Hasrer
}

// NewServer returns a new sever instance
func NewServer(s3 *wasabi.S3Service, iam *wasabi.IAMService, opts ...func(*Server)) *Server {
	s := &Server{
		S3Service:  s3,
		IAMService: iam,
		Address:    ":8080",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithAddress sets the listening address and port for the server
func WithAddress(address string) func(*Server) {
	return func(s *Server) {
		s.Address = ":" + address
	}
}

// Open opens the server and listens at the specifed address
func (s *Server) Open() error {
	// Open socket.
	ln, err := net.Listen("tcp", s.Address)

	if err != nil {
		return err
	}

	s.listener = ln

	// Start HTTP server. Note this is non-blocking so
	// we must block in the calling code
	go func() { http.Serve(s.listener, s.router()) }()

	log.Println("Server started listening on port " + s.Address[1:] + "....")

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

	wh := newWasabiHandler(s.S3Service, s.IAMService)

	s.Router.Handle("/backups/cloud/infrastructure/create", wh.handleCreateBackupInfrastructure())

	return s.Router
}
