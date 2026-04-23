package engine

import (
	"bytes"
	"encoding/binary"
	"math/rand/v2"
	"server/atlas"
	"sync"
	"time"
)

var Worlds []*Engine = []*Engine{
	CreateEngine(CreateIsland(800)),
}

type Engine struct {
	Characters  map[uint32]*Character
	Npcs        map[uint32]*Npc
	Projectiles map[uint32]*Projectile
	World       *atlas.World
	mu          sync.Mutex
}

func (engine *Engine) Pack(packet_type uint8) []byte {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	data := new(bytes.Buffer)

	data.WriteByte(packet_type)
	binary.Write(data, binary.LittleEndian, uint16(len(engine.Characters)))
	for id, character := range engine.Characters {
		binary.Write(data, binary.LittleEndian, id)
		data.Write(character.Pack())
	}
	binary.Write(data, binary.LittleEndian, uint16(len(engine.Npcs)))
	for id, npc := range engine.Npcs {
		binary.Write(data, binary.LittleEndian, id)
		data.Write(npc.Pack())
	}
	return data.Bytes()
}

func CreateEngine(world *atlas.World) *Engine {
	return &Engine{
		Characters:  make(map[uint32]*Character),
		Npcs:        make(map[uint32]*Npc),
		Projectiles: make(map[uint32]*Projectile),
		World:       world,
	}
}

func (engine *Engine) AddCharacter(id uint32, character *Character) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	engine.Characters[id] = character
}

func (engine *Engine) RemoveCharacter(id uint32, lock bool) {
	if lock {
		engine.mu.Lock()
		defer engine.mu.Unlock()
	}

	character, exists := engine.Characters[id]
	if !exists {
		return
	}
	character.Dead = true

	for _, npc := range engine.Npcs {
		npc.ExitView(id, character)
	}

	delete(engine.Characters, id)
}

func (engine *Engine) CreateProjectile(which uint8, x float32, y float32, angle uint16, evil bool, damage uint16) uint32 {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	id := rand.Uint32()
	projectile := DefaultProjectile(which, x, y, angle, evil, damage)
	engine.Projectiles[id] = projectile

	return id
}

func (engine *Engine) AddProjectile(id uint32, projectile *Projectile) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	engine.Projectiles[id] = projectile
}

func (engine *Engine) RemoveNpc(id uint32) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	delete(engine.Npcs, id)
}

func (engine *Engine) HasCharacter(id uint32) bool {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	_, exists := engine.Characters[id]
	return exists
}

func (engine *Engine) HasNpc(id uint32) bool {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	_, exists := engine.Npcs[id]
	return exists
}

func (engine *Engine) ForEachCharacter(f func(id uint32, character *Character)) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	for id, character := range engine.Characters {
		f(id, character)
	}
}

func (engine *Engine) ForEachNpc(f func(id uint32, npc *Npc)) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	for id, npc := range engine.Npcs {
		f(id, npc)
	}
}

func (engine *Engine) MoveCharacter(character *Character, x float32, y float32, angle uint16) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	character.Move(x, y, angle)
}

func (engine *Engine) AddNpc(id uint32, npc *Npc) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	npc.entityID = id
	engine.Npcs[id] = npc
}

func (engine *Engine) SpawnNpc(which uint8, x float32, y float32) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	id := rand.Uint32()
	npc := DefaultNpc(which, x, y)
	npc.entityID = id
	engine.Npcs[id] = npc
}

func (engine *Engine) Run() {
	for i := 0; i < 100; i++ {
		engine.SpawnNpc(1, 400, 400)
	}
	for i := 0; i < 1000; i++ {
		engine.SpawnNpc(0, 400, 400)
	}

	delta := time.Millisecond * 50
	ticker := time.NewTicker(delta)
	for {
		engine.mu.Lock()
		for id, npc := range engine.Npcs {
			if len(npc.nearby) > 0 {
				npc.UpdateTarget()
				npc.Move(delta)
				npc.Tick()

				if npc.Dead {
					delete(engine.Npcs, id)
				}
			}
		}
		for id, projectile := range engine.Projectiles {
			projectile.Tick(delta)
			if projectile.Dead {
				delete(engine.Projectiles, id)
			}
		}
		engine.mu.Unlock()
		<-ticker.C
	}
}
