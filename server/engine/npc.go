package engine

import (
	"bytes"
	"encoding/binary"
	"math/rand/v2"
	"sync"
	"time"
)

type Npc struct {
	id       uint8
	entityID uint32
	x        float32
	y        float32
	health   float32
	origin   Object
	instance *Engine

	target      Object
	looking     Entity
	movement    string
	mode        uint8
	usedModes   []bool
	modeTimer   float32
	attack      uint8
	attackTimer float32

	// nearby is written by every client's StartSimulation goroutine
	// (EnterView/ExitView) and read during the global engine tick. Those run
	// under different locks, so the map needs its own mutex or Go fatals on a
	// concurrent map write.
	nearby   map[uint32]*Character
	nearbyMu sync.Mutex
	damage   map[uint32]float32
	Dead     bool
}

// NearbyChars returns a snapshot of the characters currently in view of this
// npc, so callers can range/send without holding nearbyMu.
func (npc *Npc) NearbyChars() []*Character {
	npc.nearbyMu.Lock()
	defer npc.nearbyMu.Unlock()
	out := make([]*Character, 0, len(npc.nearby))
	for _, character := range npc.nearby {
		out = append(out, character)
	}
	return out
}

// NearbyMap returns a snapshot copy of the nearby map (id -> character).
func (npc *Npc) NearbyMap() map[uint32]*Character {
	npc.nearbyMu.Lock()
	defer npc.nearbyMu.Unlock()
	out := make(map[uint32]*Character, len(npc.nearby))
	for id, character := range npc.nearby {
		out[id] = character
	}
	return out
}

func (npc *Npc) NearbyCount() int {
	npc.nearbyMu.Lock()
	defer npc.nearbyMu.Unlock()
	return len(npc.nearby)
}

func (npc *Npc) GetX() float32      { return npc.x }
func (npc *Npc) GetY() float32      { return npc.y }
func (npc *Npc) GetId() uint32      { return npc.entityID }
func (npc *Npc) GetHitbox() float32 { return float32(npcData[npc.id].Hitbox) }
func (npc *Npc) Damage(amount float32) {
	npc.health -= amount

	if npc.health <= 0 {
		// Killing blow: broadcast the death (CHARACTER_DEAD) instead of a hit
		// flash so clients remove it instantly rather than waiting out the
		// no-frames timeout.
		npc.Die()
		return
	}

	nearby := npc.NearbyChars()
	if len(nearby) == 0 {
		return
	}

	data := new(bytes.Buffer)
	data.WriteByte(byte(5))
	binary.Write(data, binary.LittleEndian, npc.entityID)
	packet := data.Bytes()

	for _, character := range nearby {
		trySend(character.send, packet)
	}
}

// Die marks the npc dead, runs its death/loot logic, and tells every nearby
// client to remove it immediately via a CHARACTER_DEAD packet. Idempotent.
func (npc *Npc) Die() {
	if npc.Dead {
		return
	}
	npc.Dead = true
	npc.Death()

	data := new(bytes.Buffer)
	data.WriteByte(byte(11)) // CHARACTER_DEAD
	binary.Write(data, binary.LittleEndian, npc.entityID)
	packet := data.Bytes()

	for _, character := range npc.NearbyChars() {
		trySend(character.send, packet)
	}
}
func (npc *Npc) Death() {
	// we would use enemies loot pool id
	data := GetNpcData(npc.id)
	lootPool := GetLootData(data.Loot)
	SBThreshold := min(float32(200), data.Health/10.0)

	nearby := npc.NearbyMap()
	for id, char := range nearby {
		damage, ok := npc.damage[id]
		for _, loot := range lootPool {
			odds := rand.Float32() <= loot.Chance
			if odds && loot.SB && ok && damage >= SBThreshold {
				l := CreateLoot(loot.Loot, npc.x, npc.y)
				char.Simulation.AddLoot(l)
			}
		}
	}
	for _, loot := range lootPool {
		odds := rand.Float32() <= loot.Chance
		if odds && !loot.SB {
			l := CreateLoot(loot.Loot, npc.x, npc.y)
			for _, char := range nearby {
				char.Simulation.AddLoot(l)
			}
		}
	}
}

func (npc *Npc) Look(obj Entity) {
	npc.looking = obj
}

func (npc *Npc) Pack() []byte {
	data := new(bytes.Buffer)

	binary.Write(data, binary.LittleEndian, npc.id)
	binary.Write(data, binary.LittleEndian, npc.x)
	binary.Write(data, binary.LittleEndian, npc.y)
	binary.Write(data, binary.LittleEndian, npc.health)

	if npc.looking != nil {
		data.WriteByte(1)
		binary.Write(data, binary.LittleEndian, npc.looking.GetId())
	} else {
		data.WriteByte(0)
	}
	return data.Bytes()
}

func (npc *Npc) UpdateTarget() {
	min_dist := GetNpcData(npc.id).Range
	found_character := false

	for _, character := range npc.NearbyChars() {
		dist := float32(Distance(npc, character))
		if dist < min_dist {
			npc.target = character
			min_dist = dist
			found_character = true
		}
	}

	if !found_character {
		if _, ok := npc.target.(*Character); ok {
			npc.target = nil
			npc.Look(nil)
		}
	}
}

func (npc *Npc) EnterView(id uint32, character *Character) {
	npc.nearbyMu.Lock()
	npc.nearby[id] = character
	npc.nearbyMu.Unlock()
}
func (npc *Npc) ExitView(id uint32, character *Character) {
	npc.nearbyMu.Lock()
	delete(npc.nearby, id)
	npc.nearbyMu.Unlock()
	if target, ok := npc.target.(*Character); ok && target == character {
		npc.target = nil
	}
}

func (npc *Npc) Data() NpcData {
	return GetNpcData(npc.id)
}

func (npc *Npc) ValidMode(idx uint8) bool {
	data := GetNpcData(npc.id)
	mode := data.Modes[idx]

	if mode.SingleUse && npc.usedModes[idx] {
		return false
	}
	if mode.MaxHealth < npc.health {
		return false
	}
	if mode.MinHealth >= npc.health {
		return false
	}
	return true
}

func (npc *Npc) NewMode() {
	npc.Look(nil)
	data := GetNpcData(npc.id)
	pool := make([]uint8, 0)

	for idx, mode := range data.Modes {
		if npc.ValidMode(uint8(idx)) {
			pool = append(pool, uint8(idx))

			if mode.Priority {
				pool = []uint8{uint8(idx)}
				break
			}
		}
	}

	if len(pool) > 0 {
		mode := pool[rand.IntN(len(pool))]
		npc.mode = mode
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
		npc.Die()
	}

	if npc.InCombat() {
		if !npc.ValidMode(npc.mode) || npc.modeTimer <= 0 {
			npc.usedModes[npc.mode] = true
			npc.NewMode()
		}

		if npc.movement == "hover" && !npc.Hovering() {
			return
		}

		if npc.attackTimer < 0 {
			npc.NewAttack()
		}
	} else {
		npc.modeTimer = 0
	}
}

func (npc *Npc) GetAttack() *AttackData {
	d := npcData[npc.id]
	mode := d.Modes[npc.mode]
	attackLen := len(mode.Attacks)

	if attackLen > 0 {
		return &mode.Attacks[int(npc.attack)%attackLen]
	}
	return nil
}

func Max(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func (npc *Npc) CanAttack() bool {
	dist := float32(Distance(npc.target, npc))
	attack := npc.GetAttack()
	var attack_range float32 = 0

	if attack != nil {
		for _, proj := range attack.Projectiles {
			attack_range = Max(attack_range, GetProjectileData(proj.ID).Range)
		}
		if len(attack.Bombs) > 0 {
			attack_range = 32
		}
		if len(attack.Summons) > 0 {
			attack_range = Max(attack_range, GetNpcData(npc.id).Range)
		}

		return (npc.InCombat() &&
			dist < attack_range)
	} else {
		return false
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
	case "hover":
		npc.Hover(delta)
	case "turret":
		npc.Turret(delta)
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
		id:     id,

		target:    nil,
		looking:   nil,
		movement:  "wander",
		mode:      0,
		usedModes: make([]bool, len(GetNpcData(id).Modes)),
		modeTimer: 0,

		attack:      0,
		attackTimer: 0,
		nearby:      make(map[uint32]*Character),
		damage:      make(map[uint32]float32),
		Dead:        false,
	}
}
