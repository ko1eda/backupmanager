package http

import (
	"log"
	"net"
	"net/http"

	"github.com/ko1eda/backupmanager/wasabi"
)

// Server represents an http server
type Server struct {
	listener   net.Listener
	S3Service  *wasabi.S3Service
	IAMService *wasabi.IAMService
	router     *http.ServeMux
	address    string
	// Hasrer
}

// NewServer returns a new sever instance
func NewServer(opts ...func(*Server)) *Server {
	s := &Server{
		address: ":8080",
	}

	s.router = http.NewServeMux()

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithAddress sets the listening address and port for the server
func WithAddress(address string) func(*Server) {
	return func(s *Server) {
		s.address = ":" + address
	}
}

// Open opens the server and listens at the specifed address
func (s *Server) Open() error {
	// create the routes for the server
	s.routes()

	// Open socket.
	ln, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}

	s.listener = ln

	// Start HTTP server. Note this is non-blocking so
	// we must block in the calling code
	go func() { http.Serve(s.listener, s.router) }()

	log.Println("Server started listening on port " + s.address[1:] + "....")

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.listener != nil {
		s.listener.Close()
	}

	return nil
}

// routes maps all route handlers to their respoective paths
func (s *Server) routes() {
	s.router.Handle("/cloud/infrastructure/create", s.handleCreateBackupInfrastructure())
}

// Router creates a new servermux and registers all routes to it
// func (s *Server) router() http.Handler {
// 	r := http.NewServeMux()

// 	wh := newWasabiHandler(s.S3Service, s.IAMService)

// 	r.Handle("/cloud/infrastructure/create", wh.handleCreateBackupInfrastructure())

// 	return r
// }
