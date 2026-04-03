package engine

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type NpcData struct {
	ID       uint8   `json:"id"`
	Name     string  `json:"display"`
	Health   uint16  `json:"health"`
	Speed    float32 `json:"speed"`
	Range    float32 `json:"range"`
	Hitbox   uint8   `json:"hitbox"`
	Movement struct {
		Combat string `json:"combat"`
		Idle   string `json:"idle"`
	} `json:"movement"`
}

type ItemData struct {
	ID      uint8   `json:"id"`
	Slot    string  `json:"type"`
	Speed   float32 `json:"speed"`
	Attacks []struct {
		Animation   uint8   `json:"animation"`
		Reload      float32 `json:"reload"`
		Projectiles []struct {
			ID    uint8   `json:"id"`
			Angle uint16  `json:"angle"`
			X     float32 `json:"x"`
			Y     float32 `json:"y"`
		} `json:"projectiles,omitempty"`
	} `json:"attacks,omitempty"`
}

type ProjectileData struct {
	ID       uint8   `json:"id"`
	Speed    float32 `json:"speed"`
	Range    float32 `json:"range"`
	Damage   uint16  `json:"damage"`
	Piercing bool    `json:"piercing"`
	Hitbox   uint8   `json:"hitbox"`
}

//go:embed assets/items.json
var itemsJSON []byte

//go:embed assets/npcs.json
var npcsJSON []byte

//go:embed assets/projectiles.json
var projectilesJSON []byte

//go:embed assets/tiles.json
var tilesJSON []byte

var npcData []NpcData
var itemData []ItemData
var projectileData []ProjectileData

func GetNpcData() []NpcData {
	return npcData
}
func GetItemData() []ItemData {
	return itemData
}
func GetProjectileData() []ProjectileData {
	return projectileData
}

func InitAssets() {
	if err := json.Unmarshal(npcsJSON, &npcData); err != nil {
		fmt.Println("Error parsing Npc JSON")
	}
	if err := json.Unmarshal(itemsJSON, &itemData); err != nil {
		fmt.Println("Error parsing Item JSON")
	}
	if err := json.Unmarshal(projectilesJSON, &projectileData); err != nil {
		fmt.Println("Error parsing Projectile JSON")
	}

	fmt.Println("Initialized", len(npcData), "NPCs")
	fmt.Println("Initialized", len(itemData), "Items")
	fmt.Println("Initialized", len(projectileData), "Projectiles")
}
