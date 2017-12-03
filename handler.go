package pacmound

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

var (
	trueStr = "true"
)

const (
	MaxLoops = 300
)

func NewGameMux(getGopher, getPython AgentGetter) http.Handler {
	mut := &sync.Mutex{}
	mux := http.NewServeMux()
	serveFile(mux, "/", "../../src/index.html", "")
	serveFile(mux, "/src/gopher.png", "../../src/gopher.png", "")
	serveFile(mux, "/src/python.png", "../../src/python.png", "")
	serveFile(mux, "/src/dirt.png", "../../src/dirt.png", "")
	serveFile(mux, "/src/stone.png", "../../src/stone.png", "")
	serveFile(mux, "/src/carrot.png", "../../src/carrot.png", "")

	mux.HandleFunc("/api/level/00", LevelHandler(level00, getGopher, getPython, mut))
	mux.HandleFunc("/api/level/01", LevelHandler(level01, getGopher, getPython, mut))
	mux.HandleFunc("/api/level/02", LevelHandler(level02, getGopher, getPython, mut))
	mux.HandleFunc("/api/level/03", LevelHandler(level03, getGopher, getPython, mut))
	mux.HandleFunc("/api/level/04", LevelHandler(level04, getGopher, getPython, mut))
	return mux
}

type LevelData struct {
	MaxSteps int                `json:"maxSteps"`
	Scores   []float64          `json:"scores"`
	States   [][][]EncodedBlock `json:"states"`
	Agent    Agent              `json:"agent,omitempty"`
}

type LevelFunc func(getGopher, getPython AgentGetter, loop func(m *Maze, agentData *AgentData) bool)

func LevelHandler(levelFunc LevelFunc, getGopher, getPython AgentGetter, mut *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mut.Lock()
		defer mut.Unlock()

		maxLoops := MaxLoops
		loopLimit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			loopLimit = maxLoops
		}
		loopCount := 0

		data := LevelData{}
		data.MaxSteps = loopLimit

		levelFunc(getGopher, getPython, func(m *Maze, agentData *AgentData) bool {
			data.States = append(data.States, m.encodable())
			data.Scores = append(data.Scores, agentData.score)

			remReward := m.RemainingReward()

			if !m.loop() || remReward <= 0 || loopCount > loopLimit || agentData.dead {
				data.Scores = append(data.Scores, agentData.score)
				return false
			}
			loopCount++
			return true
		})

		data.Agent = getGopher()
		json.NewEncoder(w).Encode(data)
	}
}

func serveFile(mux *http.ServeMux, url, path, contentType string) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		http.ServeFile(w, r, path)
	})
}
