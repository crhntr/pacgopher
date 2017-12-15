package pacmound

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	trueStr = "true"
)

const (
	MaxLoops = 300
)

func NewGameMux(getGopher, getPython AgentGetter) http.Handler {
	checkFiles()
	mux := http.NewServeMux()
	serveFile(mux, "/", "../../src/index.html", "")
	serveFile(mux, "/src/gopher.png", "../../src/assets/img/gopher.png", "")
	serveFile(mux, "/src/python.png", "../../src/assets/img/python.png", "")
	serveFile(mux, "/src/dirt.png", "../../src/assets/img/dirt.png", "")
	serveFile(mux, "/src/stone.png", "../../src/assets/img/stone.png", "")
	serveFile(mux, "/src/carrot.png", "../../src/assets/img/carrot.png", "")

	mux.HandleFunc("/api/level/00", LevelHandler(level00, getGopher, getPython))
	mux.HandleFunc("/api/level/01", LevelHandler(level01, getGopher, getPython))
	mux.HandleFunc("/api/level/02", LevelHandler(level02, getGopher, getPython))
	mux.HandleFunc("/api/level/03", LevelHandler(level03, getGopher, getPython))
	mux.HandleFunc("/api/level/04", LevelHandler(level04, getGopher, getPython))
	return mux
}

func checkFiles() {
	paths := []string{
		"../../src/assets/img/gopher.png",
		"../../src/assets/img/python.png",
		"../../src/assets/img/dirt.png",
		"../../src/assets/img/stone.png",
		"../../src/assets/img/carrot.png"}

	pathNotFound := false
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			pathNotFound = true
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}
	if pathNotFound {
		panic("some path not found")
	}
}

type LevelData struct {
	MaxSteps int                `json:"maxSteps"`
	Scores   []float64          `json:"scores"`
	States   [][][]EncodedBlock `json:"states"`
	Agent    Agent              `json:"agent,omitempty"`
}

type LevelFunc func(getGopher, getPython AgentGetter, loop func(m *Maze, agentData *AgentData) bool)

func LevelHandler(levelFunc LevelFunc, getGopher, getPython AgentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loopCount := 0

		data := LevelData{}
		data.MaxSteps = MaxLoops
		var agent Agent
		levelFunc(getGopher, getPython, func(m *Maze, agentData *AgentData) bool {
			data.States = append(data.States, m.encodable())
			data.Scores = append(data.Scores, agentData.score)

			if agent == nil {
				agent = agentData.a
			}

			remReward := m.RemainingReward()

			if !m.loop() || remReward <= 0 {
				data.Scores = append(data.Scores, agentData.score)
				return false
			}
			loopCount++
			return true
		})

		data.Agent = agent
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Print(data)
			log.Print(err)
		}
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
