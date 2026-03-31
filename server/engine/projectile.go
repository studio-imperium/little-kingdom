package engine

import (
	"bytes"
	"encoding/binary"
)

type Projectile struct {
	id      uint8
	evil    bool
	x       float32
	y       float32
	angle   uint16
	hitlist map[uint32]Object
	Dead    bool
}

func (p Projectile) GetX() float32 { return p.x }
func (p Projectile) GetY() float32 { return p.y }

func (projectile *Projectile) Pack(packet_type uint8) []byte {
	data := new(bytes.Buffer)

	data.WriteByte(packet_type)
	binary.Write(data, binary.LittleEndian, projectile.x)
	binary.Write(data, binary.LittleEndian, projectile.y)
	binary.Write(data, binary.LittleEndian, projectile.angle)

	return data.Bytes()
}

func DefaultProjectile(id uint8, x float32, y float32, angle uint16, evil bool) *Projectile {
	return &Projectile{
		id:      id,
		evil:    evil,
		x:       x,
		y:       y,
		angle:   angle,
		hitlist: make(map[uint32]Object),
		Dead:    false,
	}
}
