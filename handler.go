package pacmound

import (
	"net/http"
)

type GameServer struct {
	mux http.Handler
}

func NewGameMux(agent Agent) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../pacman/src/index.html")
	})
	mux.HandleFunc("/src/vue.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../pacman/src/vue.js")
	})
	mux.HandleFunc("/api/level/00", Level00Handler(agent))
	mux.HandleFunc("/api/level/01", Level01Handler(agent))

	return mux
}
