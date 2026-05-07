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
	MaxHealth float32      `json:"max_health,omitempty"`
	MinHealth float32      `json:"min_health,omitempty"`
	SingleUse bool         `json:"single_use,omitempty"`
	OnSpawn   bool         `json:"on_spawn,omitempty"`
	Movement  string       `json:"movement"`
	Attacks   []AttackData `json:"attacks,omitempty"`
}

type NpcData struct {
	ID     uint8         `json:"id"`
	Name   string        `json:"display"`
	Health float32       `json:"health"`
	Speed  float32       `json:"speed"`
	Range  float32       `json:"range"`
	Hitbox uint8         `json:"hitbox"`
	Modes  []NpcModeData `json:"modes,omitempty"`
}

type SpawnData struct {
	Chance float32      `json:"chance"`
	Npcs   []SummonData `json:"npcs"`
}

type LootData struct {
	Loot   uint8   `json:"loot"`
	Chance float32 `json:"chance"`
	SB     bool    `json:"soulbound"`
}

type Stats struct {
	Health float32 `json:"health,omitempty"`
	Regen  float32 `json:"regen,omitempty"`
	Speed  float32 `json:"speed,omitempty"`
	Damage float32 `json:"damage,omitempty"`
	Reload float32 `json:"reload,omitempty"`
}

type ItemData struct {
	ID      uint8        `json:"id"`
	Slot    string       `json:"type"`
	Stats   Stats        `json:"stats"`
	Attacks []AttackData `json:"attacks,omitempty"`
}

type ProjectileData struct {
	ID       uint8   `json:"id"`
	Speed    float32 `json:"speed"`
	Range    float32 `json:"range"`
	Damage   float32 `json:"damage"`
	Piercing bool    `json:"piercing"`
	Hitbox   uint8   `json:"hitbox"`
}

type BombData struct {
	ID      uint8   `json:"id"`
	Airtime float32 `json:"airtime"`
	Damage  float32 `json:"damage"`
	Radius  uint8   `json:"radius"`
}

//go:embed assets/items.json
var itemsJSON []byte

//go:embed assets/npcs.json
var npcsJSON []byte

//go:embed assets/spawns.json
var spawnsJSON []byte

//go:embed assets/loot.json
var lootJSON []byte

//go:embed assets/projectiles.json
var projectilesJSON []byte

//go:embed assets/bombs.json
var bombsJSON []byte

//go:embed assets/tiles.json
var tilesJSON []byte

var npcData []NpcData
var spawnsData [][]SpawnData
var lootData [][]LootData
var itemData []ItemData
var projectileData []ProjectileData
var bombData []BombData

var biomeSpawns map[uint8][]SpawnData = map[uint8][]SpawnData{}

func GetNpcData(id uint8) NpcData {
	return npcData[id]
}
func GetSpawnsData(id uint8) []SpawnData {
	return spawnsData[id]
}
func GetLootData(id uint8) []LootData {
	return lootData[id]
}
func GetItemData(id uint8) ItemData {
	return itemData[id]
}
func GetProjectileData(id uint8) ProjectileData {
	return projectileData[id]
}
func GetBombData(id uint8) BombData {
	return bombData[id]
}

func InitAssets() {
	if err := json.Unmarshal(npcsJSON, &npcData); err != nil {
		fmt.Println("Error parsing Npc JSON")
	}
	if err := json.Unmarshal(spawnsJSON, &spawnsData); err != nil {
		fmt.Println("Error parsing Spawns JSON")
	}
	if err := json.Unmarshal(lootJSON, &lootData); err != nil {
		fmt.Println("Error parsing Loot JSON")
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

	biomeSpawns[20] = spawnsData[0]
	biomeSpawns[21] = spawnsData[0]
	biomeSpawns[22] = spawnsData[0]

	fmt.Println("Initialized", len(npcData), "NPCs")
	fmt.Println("Initialized", len(spawnsData), "Spawns")
	fmt.Println("Initialized", len(lootData), "Loot Tables")
	fmt.Println("Initialized", len(itemData), "Items")
	fmt.Println("Initialized", len(projectileData), "Projectiles")
	fmt.Println("Initialized", len(bombData), "Bombs")
}
