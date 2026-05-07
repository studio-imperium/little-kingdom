package engine

import (
	"bytes"
	"encoding/binary"
)

type Character struct {
	id    uint32
	x     float32
	y     float32
	angle uint16
	hand  uint8
	head  uint8
	body  uint8

	maxHealth float32
	health    float32
	regen     float32
	speed     float32
	damage    float32
	reload    float32

	inventory      map[uint8]uint8
	send           *chan []byte
	AttackCounter  uint8
	AttackCooldown float32
	Simulation     *Engine
	Dead           bool
}

func (c *Character) GetX() float32      { return c.x }
func (c *Character) GetY() float32      { return c.y }
func (c *Character) GetId() uint32      { return c.id }
func (c *Character) GetHitbox() float32 { return 1 }
func (c *Character) Damage(amount uint16) {
	c.health -= float32(amount)

	if c.health < 1 {
		*c.send <- []byte{11}
	}

	c.SetHealth(c.health)
}

func (character *Character) SetHealth(val float32) {
	character.health = val

	data := new(bytes.Buffer)
	data.WriteByte(byte(10))
	binary.Write(data, binary.LittleEndian, uint16(val))
	binary.Write(data, binary.LittleEndian, uint16(character.maxHealth))
	packet := data.Bytes()

	*character.send <- packet
}

func (engine *Engine) GetHand(id uint32) uint8 {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	char := engine.Characters[id]
	return char.hand
}

func (engine *Engine) GetSlot(id uint32, idx uint8) uint8 {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	char := engine.Characters[id]
	return char.inventory[idx]
}

func (engine *Engine) ChangeInventory(id uint32, to uint8, from uint8) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	char := engine.Characters[id]

	remove_gear := func(slot_type string, gear_slot *uint8) {
		if item, ok := char.inventory[to]; ok {
			slot := GetItemData(item).Slot
			if slot == slot_type && char.head != 0 {
				char.inventory[to] = *gear_slot
				*gear_slot = item
			}
		} else {
			char.inventory[to] = *gear_slot

			if slot_type == "head" {
				*gear_slot = 0
			} else {
				*gear_slot = 1
			}
		}
	}
	equip_gear := func(slot_type string, gear_slot *uint8) {
		if item, ok := char.inventory[from]; ok {
			slot := GetItemData(item).Slot
			if slot == slot_type && char.head != 0 {
				char.inventory[from] = *gear_slot
				*gear_slot = item
			} else if slot == slot_type {
				delete(char.inventory, from)
				*gear_slot = item
			}
		}
	}

	if to == 24 && from == 25 || to == 25 && from == 24 {
		return
	} else if to == 24 {
		equip_gear("head", &char.head)
	} else if from == 24 {
		remove_gear("head", &char.head)
	} else if to == 25 {
		equip_gear("body", &char.body)
	} else if from == 25 {
		remove_gear("body", &char.body)
	} else if item1, ok := char.inventory[from]; ok {
		if item2, ok := char.inventory[to]; ok {
			char.inventory[from] = item2
		} else {
			delete(char.inventory, from)
		}
		char.inventory[to] = item1
	}
	char.Apply()
}

func (engine *Engine) SelectSlot(id uint32, new uint8) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	char := engine.Characters[id]
	char.hand = new
}

func (engine *Engine) PackCharacter(id uint32, packet_type uint8) []byte {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	char := engine.Characters[id]
	return char.PackFull(packet_type)
}

func (character *Character) PackFull(packet_type uint8) []byte {
	data := new(bytes.Buffer)

	data.WriteByte(packet_type)
	binary.Write(data, binary.LittleEndian, character.x)
	binary.Write(data, binary.LittleEndian, character.y)
	binary.Write(data, binary.LittleEndian, character.angle)
	binary.Write(data, binary.LittleEndian, uint16(character.health))
	binary.Write(data, binary.LittleEndian, uint16(character.maxHealth))
	data.WriteByte(character.hand)
	data.WriteByte(character.head)
	data.WriteByte(character.body)

	data.WriteByte(uint8(len(character.inventory)))
	for slot, id := range character.inventory {
		data.WriteByte(slot)
		data.WriteByte(id)
	}
	return data.Bytes()
}

func (character *Character) Pack() []byte {
	data := new(bytes.Buffer)

	binary.Write(data, binary.LittleEndian, character.x)
	binary.Write(data, binary.LittleEndian, character.y)
	binary.Write(data, binary.LittleEndian, character.angle)
	binary.Write(data, binary.LittleEndian, uint16(character.health))
	data.WriteByte(character.inventory[character.hand])
	data.WriteByte(character.head)
	data.WriteByte(character.body)

	return data.Bytes()
}

var weapon uint8 = 4

func DefaultCharacter(simulation *Engine, send *chan []byte, id uint32) *Character {
	weapon += 1
	weapon %= 4
	return &Character{
		id:     id,
		x:      0,
		y:      0,
		angle:  0,
		hand:   0,
		head:   0,
		body:   1,
		health: 1000,
		inventory: map[uint8]uint8{
			0: 4,
		},
		send:           send,
		AttackCounter:  0,
		AttackCooldown: 0,
		Simulation:     simulation,
		Dead:           false,
	}
}

func (character *Character) Tick(delta_sec float32) {
	character.AttackCooldown -= delta_sec

	if character.AttackCooldown < -1 {
		character.AttackCooldown = -1
	}

	tmp := uint16(character.health)
	character.health += character.regen * delta_sec

	if character.health > character.maxHealth {
		character.health = character.maxHealth
	}
	if tmp != uint16(character.health) {
		character.SetHealth(character.health)
	}
}

func (character *Character) Move(x float32, y float32, angle uint16) {
	character.x = x
	character.y = y
	character.angle = angle
}

func NonZero(a float32) float32 {
	if a != 0 {
		return a
	}
	return 1
}

func (character *Character) Apply() {
	var health float32 = 20
	var regen float32 = 1
	var speed float32 = 1
	var damage float32 = 1
	var reload float32 = 1

	helmet := GetItemData(character.head)
	body := GetItemData(character.body)

	health *= NonZero(helmet.Stats.Health)
	health *= NonZero(body.Stats.Health)

	regen *= NonZero(helmet.Stats.Regen)
	regen *= NonZero(body.Stats.Regen)

	speed *= NonZero(helmet.Stats.Speed)
	speed *= NonZero(body.Stats.Speed)

	damage *= NonZero(helmet.Stats.Damage)
	damage *= NonZero(body.Stats.Damage)

	reload *= NonZero(helmet.Stats.Reload)
	reload *= NonZero(body.Stats.Reload)

	character.maxHealth = health
	character.regen = regen
	character.speed = speed
	character.damage = damage
	character.reload = reload

	data := character.PackFull(0)
	*character.send <- data
}

func (character *Character) Attack(x float32, y float32, angle uint16) {
	character.x = x
	character.y = y
	character.angle = angle
}
