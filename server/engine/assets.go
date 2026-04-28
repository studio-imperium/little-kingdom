package engine

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type ProjectileSpawnData struct {
	ID    uint8   `json:"id"`
	Angle uint16  `json:"angle"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
}

type BombSpawnData struct {
	ID uint8   `json:"id"`
	X  float32 `json:"x"`
	Y  float32 `json:"y"`
}

type SummonData struct {
	ID uint8   `json:"id"`
	X  float32 `json:"x"`
	Y  float32 `json:"y"`
}

type AttackData struct {
	Animation   uint8                 `json:"animation"`
	Reload      float32               `json:"reload"`
	Wait        float32               `json:"wait,omitempty"`
	Projectiles []ProjectileSpawnData `json:"projectiles,omitempty"`
	Bombs       []BombSpawnData       `json:"bombs,omitempty"`
	Summons     []SummonData          `json:"summons,omitempty"`
}

type NpcModeData struct {
	Duration  float32      `json:"duration"`
	MaxHealth uint16       `json:"max_health,omitempty"`
	MinHealth uint16       `json:"min_health,omitempty"`
	SingleUse bool         `json:"single_use,omitempty"`
	OnSpawn   bool         `json:"on_spawn,omitempty"`
	Movement  string       `json:"movement"`
	Attacks   []AttackData `json:"attacks,omitempty"`
}

type NpcData struct {
	ID     uint8         `json:"id"`
	Name   string        `json:"display"`
	Health uint16        `json:"health"`
	Speed  float32       `json:"speed"`
	Range  float32       `json:"range"`
	Hitbox uint8         `json:"hitbox"`
	Modes  []NpcModeData `json:"modes,omitempty"`
}

type ItemData struct {
	ID      uint8        `json:"id"`
	Slot    string       `json:"type"`
	Speed   float32      `json:"speed"`
	Attacks []AttackData `json:"attacks,omitempty"`
}

type ProjectileData struct {
	ID       uint8   `json:"id"`
	Speed    float32 `json:"speed"`
	Range    float32 `json:"range"`
	Damage   uint16  `json:"damage"`
	Piercing bool    `json:"piercing"`
	Hitbox   uint8   `json:"hitbox"`
}

type BombData struct {
	ID      uint8   `json:"id"`
	Airtime float32 `json:"airtime"`
	Damage  uint16  `json:"damage"`
	Radius  uint8   `json:"radius"`
}

//go:embed assets/items.json
var itemsJSON []byte

//go:embed assets/npcs.json
var npcsJSON []byte

//go:embed assets/projectiles.json
var projectilesJSON []byte

//go:embed assets/bombs.json
var bombsJSON []byte

//go:embed assets/tiles.json
var tilesJSON []byte

var npcData []NpcData
var itemData []ItemData
var projectileData []ProjectileData
var bombData []BombData

func GetNpcData() []NpcData {
	return npcData
}
func GetItemData() []ItemData {
	return itemData
}
func GetProjectileData() []ProjectileData {
	return projectileData
}
func GetBombData() []BombData {
	return bombData
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
	if err := json.Unmarshal(bombsJSON, &bombData); err != nil {
		fmt.Println("Error parsing Bomb JSON")
	}

	fmt.Println("Initialized", len(npcData), "NPCs")
	fmt.Println("Initialized", len(itemData), "Items")
	fmt.Println("Initialized", len(projectileData), "Projectiles")
	fmt.Println("Initialized", len(bombData), "Bombs")
}
