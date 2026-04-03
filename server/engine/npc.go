package engine

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Npc struct {
	id     uint8
	x      float32
	y      float32
	health uint16
	target Object
	origin Object
	nearby map[uint32]*Character
	Dead   bool
}

func (npc Npc) GetX() float32 { return npc.x }
func (npc Npc) GetY() float32 { return npc.y }
func (npc *Npc) Damage(amount uint16) {
	npc.health -= amount
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

func (npc *Npc) Tick() {
	if npc.health <= 0 {
		npc.Dead = true
	}
}

func (npc *Npc) Move(delta time.Duration) {
	if _, ok := npc.target.(*Character); ok {
		switch npcData[npc.id].Movement.Combat {
		case "chase":
			npc.Chase(delta)
		case "run":
			npc.Run(delta)
		case "overshoot":
			npc.Overshoot(delta)
		default:
			return
		}
	} else {
		switch npcData[npc.id].Movement.Idle {
		case "wander":
			npc.Wander(delta)
		case "travel":
			npc.Travel(delta)
		default:
			return
		}
	}
}

func DefaultNpc(id uint8, x float32, y float32) *Npc {
	health := npcData[id].Health

	return &Npc{
		id:     id,
		x:      x,
		y:      y,
		health: health,
		target: nil,
		origin: Point{x, y},
		nearby: make(map[uint32]*Character),
		Dead:   false,
	}
}
