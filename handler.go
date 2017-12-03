package pacmound

import (
	"net/http"
)

func NewGameMux(gopher Agent, python1, python2, python3 Agent) http.Handler {
	mux := http.NewServeMux()
	serveFile(mux, "/", "../../src/index.html", "")
	serveFile(mux, "/src/gopher.png", "../../src/gopher.png", "")
	serveFile(mux, "/src/python.png", "../../src/python.png", "")
	serveFile(mux, "/src/dirt.png", "../../src/dirt.png", "")
	serveFile(mux, "/src/stone.png", "../../src/stone.png", "")
	serveFile(mux, "/src/carrot.png", "../../src/carrot.png", "")

	mux.HandleFunc("/api/level/00", Level00Handler(gopher))
	mux.HandleFunc("/api/level/01", Level01Handler(gopher, python1))
	mux.HandleFunc("/api/level/02", Level02Handler(gopher, python1))
	return mux
}

func serveFile(mux *http.ServeMux, url, path, contentType string) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		http.ServeFile(w, r, path)
	})
}
