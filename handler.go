package pacmound

import (
	"net/http"
)

var (
	trueStr = "true"
)

const (
	MaxLoops = 300
)

func NewGameMux(getGopher, getPython AgentGetter) http.Handler {
	mux := http.NewServeMux()
	serveFile(mux, "/", "../../src/index.html", "")
	serveFile(mux, "/src/gopher.png", "../../src/gopher.png", "")
	serveFile(mux, "/src/python.png", "../../src/python.png", "")
	serveFile(mux, "/src/dirt.png", "../../src/dirt.png", "")
	serveFile(mux, "/src/stone.png", "../../src/stone.png", "")
	serveFile(mux, "/src/carrot.png", "../../src/carrot.png", "")

	mux.HandleFunc("/api/level/00", Level00Handler(getGopher()))
	mux.HandleFunc("/api/level/01", Level01Handler(getGopher, getPython))
	mux.HandleFunc("/api/level/02", Level02Handler(getGopher, getPython))
	mux.HandleFunc("/api/level/03", Level03Handler(getGopher, getPython))
	mux.HandleFunc("/api/level/04", Level04Handler(getGopher, getPython))
	return mux
}

type LevelData struct {
	MaxSteps int                `json:"maxSteps"`
	Scores   []float64          `json:"scores"`
	States   [][][]EncodedBlock `json:"states"`
	Agent    Agent              `json:"agent,omitempty"`
}

func serveFile(mux *http.ServeMux, url, path, contentType string) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		http.ServeFile(w, r, path)
	})
}
