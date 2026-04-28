package engine

import (
	"bytes"
	"encoding/binary"
	"math/rand/v2"
	"time"
)

type Npc struct {
	id       uint8
	entityID uint32
	x        float32
	y        float32
	health   uint16
	origin   Object
	instance *Engine

	target      Object
	movement    string
	mode        uint8
	usedModes   []bool
	modeTimer   float32
	attack      uint8
	attackTimer float32

	nearby map[uint32]*Character
	Dead   bool
}

func (npc Npc) GetX() float32      { return npc.x }
func (npc Npc) GetY() float32      { return npc.y }
func (npc Npc) GetHitbox() float32 { return float32(npcData[npc.id].Hitbox) }
func (npc *Npc) Damage(amount uint16) {
	npc.health -= amount

	if len(npc.nearby) == 0 {
		return
	}

	data := new(bytes.Buffer)
	data.WriteByte(byte(5))
	binary.Write(data, binary.LittleEndian, npc.entityID)
	packet := data.Bytes()

	for _, character := range npc.nearby {
		*character.send <- packet
	}

	if npc.health <= 0 {
		npc.Dead = true
	}
}

func (npc *Npc) Pack() []byte {
	data := new(bytes.Buffer)

	binary.Write(data, binary.LittleEndian, npc.id)
	binary.Write(data, binary.LittleEndian, npc.x)
	binary.Write(data, binary.LittleEndian, npc.y)
	binary.Write(data, binary.LittleEndian, npc.health)

	return data.Bytes()
}

func (npc *Npc) UpdateTarget() {
	min_dist := float64(npcData[npc.id].Range)
	found_character := false

	for _, character := range npc.nearby {
		dist := Distance(npc, character)
		if dist < min_dist {
			npc.target = character
			min_dist = dist
			found_character = true
		}
	}

	if !found_character {
		if _, ok := npc.target.(*Character); ok {
			npc.target = nil
		}
	}
}

func (npc *Npc) EnterView(id uint32, character *Character) {
	npc.nearby[id] = character
}
func (npc *Npc) ExitView(id uint32, character *Character) {
	delete(npc.nearby, id)
	if target, ok := npc.target.(*Character); ok && target == character {
		npc.target = nil
	}
}

func (npc *Npc) Data() NpcData {
	return GetNpcData()[npc.id]
}
func (npc *Npc) ValidMode(idx uint8) bool {
	data := GetNpcData()[npc.id]
	mode := data.Modes[idx]

	if mode.SingleUse && npc.usedModes[idx] {
		return false
	}
	if mode.MaxHealth < npc.health {
		return false
	}
	if mode.MinHealth > npc.health {
		return false
	}
	return true
}

func (npc *Npc) NewMode() {
	data := GetNpcData()[npc.id]
	pool := make([]uint8, 0)

	for idx := range data.Modes {
		if npc.ValidMode(uint8(idx)) {
			pool = append(pool, uint8(idx))
		}
	}

	if len(pool) > 0 {
		mode := pool[rand.IntN(len(pool))]
		npc.mode = mode
		npc.usedModes[mode] = true
		npc.modeTimer = data.Modes[mode].Duration
		npc.movement = data.Modes[mode].Movement
		npc.attackTimer = 0
	}
}
func (npc *Npc) Tick(delta time.Duration) {
	deltaMs := float32(delta) / float32(time.Millisecond)
	npc.modeTimer -= deltaMs / 1000.0
	npc.attackTimer -= deltaMs / 1000.0

	if npc.health <= 0 {
		npc.Dead = true
	}

	if !npc.ValidMode(npc.mode) || npc.modeTimer < 0 {
		npc.NewMode()
	}

	if npc.attackTimer < 0 {
		npc.NewAttack()
	}
}

func (npc *Npc) InCombat() bool {
	_, ok := npc.target.(*Character)
	return ok
}

func (npc *Npc) Move(delta time.Duration) {
	if !npc.InCombat() {
		npc.movement = "wander"
	}
	switch npc.movement {
	case "wander":
		npc.Wander(delta)
	case "chase":
		npc.Chase(delta)
	case "run":
		npc.Run(delta)
	case "overshoot":
		npc.Overshoot(delta)
	default:
		return
	}
}

func DefaultNpc(id uint8, x float32, y float32) *Npc {
	health := npcData[id].Health

	return &Npc{
		x:      x,
		y:      y,
		health: health,
		origin: Point{x, y},

		target:    nil,
		movement:  "wander",
		mode:      0,
		usedModes: make([]bool, len(GetNpcData()[id].Modes)),
		modeTimer: 0,

		attack:      0,
		attackTimer: 0,
		nearby:      make(map[uint32]*Character),
		Dead:        false,
	}
}
