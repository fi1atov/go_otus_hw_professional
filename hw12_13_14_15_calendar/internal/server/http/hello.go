package internalhttp

import "net/http"

func (s *server) handleHello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, world\n"))
}
