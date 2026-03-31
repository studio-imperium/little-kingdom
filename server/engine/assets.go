package engine

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type NpcData struct {
	ID       uint8   `json:"id"`
	Name     string  `json:"display"`
	Health   int16   `json:"health"`
	Speed    float32 `json:"speed"`
	Range    float32 `json:"range"`
	Movement struct {
		Combat string `json:"combat"`
		Idle   string `json:"idle"`
	} `json:"movement"`
}

//go:embed assets/items.json
var itemsJSON []byte

//go:embed assets/npcs.json
var npcsJSON []byte

//go:embed assets/tiles.json
var tilesJSON []byte

var npcData []NpcData

func InitAssets() {
	if err := json.Unmarshal(npcsJSON, &npcData); err != nil {
		fmt.Println("Error parsing Npc JSON")
	}

	fmt.Println("Initialized", len(npcData), "NPCs")
}
