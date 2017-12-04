package pacmound

import "encoding/json"

type Block struct {
	obsticle bool
	reward   float64
	agent    *AgentData
}

func (b Block) IsObstructed() bool         { return b.obsticle }
func (b Block) IsOccupied() bool           { return b.agent != nil }
func (b Block) IsOccupiedWithPython() bool { return b.agent != nil && b.agent.t < 0 }

func (b Block) Reward() float64 {
	if b.IsObstructed() {
		return -float64(DamageColision)
	}
	return b.reward
}

type EncodedBlock struct {
	Obsticle bool      `json:"obsticle"`
	Reward   float64   `json:"reward"`
	Agent    AgentType `json:"agent"`
}

func (b Block) encodable() EncodedBlock {
	eb := EncodedBlock{
		Obsticle: b.obsticle,
		Reward:   b.reward,
	}
	if b.agent != nil {
		eb.Agent = b.agent.t
	}
	return eb
}

func (b Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.encodable())
}
