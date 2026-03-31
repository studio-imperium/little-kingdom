package engine

import (
	"bytes"
	"encoding/binary"
)

type Character struct {
	x         float32
	y         float32
	angle     uint16
	health    uint16
	hand      uint8
	head      uint8
	body      uint8
	inventory map[uint8]uint8
	Dead      bool
}

func (c Character) GetX() float32 { return c.x }
func (c Character) GetY() float32 { return c.y }

func (character *Character) PackFull(packet_type uint8) []byte {
	data := new(bytes.Buffer)

	data.WriteByte(packet_type)
	binary.Write(data, binary.LittleEndian, character.x)
	binary.Write(data, binary.LittleEndian, character.y)
	binary.Write(data, binary.LittleEndian, character.angle)
	binary.Write(data, binary.LittleEndian, character.health)
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
	binary.Write(data, binary.LittleEndian, character.health)
	data.WriteByte(character.hand)
	data.WriteByte(character.head)
	data.WriteByte(character.body)

	return data.Bytes()
}

func DefaultCharacter() *Character {
	return &Character{
		x:     0,
		y:     0,
		angle: 0,
		hand:  6,
		head:  0,
		body:  2,
		inventory: map[uint8]uint8{
			0: 6,
			1: 5,
		},
		Dead: false,
	}
}

func (character *Character) Move(x float32, y float32, angle uint16) {
	character.x = x
	character.y = y
	character.angle = angle
}
