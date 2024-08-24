package internalhttp

import "net/http"

func (s *Server) handleHello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, world\n"))
}
