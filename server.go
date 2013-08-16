package tigertonic

import (
	"net"
	"net/http"
	"sync"
)

// Server is an http.Server with better defaults.
type Server struct {
	http.Server
}

// NewServer returns an http.Server with better defaults.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{http.Server{
		Addr: addr,
		Handler: &serverHandler{
			handler:   handler,
			waitGroup: &sync.WaitGroup{},
		},
		MaxHeaderBytes: 4096,
		ReadTimeout:    1e9,
		WriteTimeout:   1e9,
	}}
}

func (s *Server) ListenAndServe() error {
	return s.Server.ListenAndServe()
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) Serve(l net.Listener) error {
	return s.Server.Serve(l)
}

type serverHandler struct {
	handler   http.Handler
	waitGroup *sync.WaitGroup
}

func (sh *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO sh.waitGroup.Add(1) // FIXME Wrong, needs to happen before go c.serve() in net/http/server.go.
	// TODO defer sh.waitGroup.Done()
	// r.Header.Set("Host", r.Host) // Should I?
	r.URL.Host = r.Host
	if nil != r.TLS {
		r.URL.Scheme = "https"
	} else {
		r.URL.Scheme = "http"
	}
	sh.handler.ServeHTTP(w, r)
}
