package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"server/engine"

	"github.com/gorilla/websocket"
)

type Client struct {
	id         uint32
	conn       *websocket.Conn
	send       chan []byte
	character  *engine.Character
	instance   *engine.Engine
	simulation *engine.Engine
}

func (client *Client) destroy() {
	fmt.Println("Client destroyed: ", client.id)
	client.conn.Close()
	client.instance.RemoveCharacter(client.id, true)
	delete(Clients, client.id)
}

func CreateClient(conn *websocket.Conn) {
	client := Client{
		id:         0,
		conn:       conn,
		send:       make(chan []byte),
		character:  nil,
		simulation: engine.CreateEngine(),
	}
	go client.recievePackets()
	go client.writePackets()
}

func (client *Client) writePackets() {
	conn := client.conn
	for {
		message, _ := <-client.send
		err := conn.WriteMessage(websocket.BinaryMessage, message)

		if err != nil {
			return
		}
	}
}

func (client *Client) recievePackets() {
	conn := client.conn
	defer client.destroy()

	for {
		messageType, r, err := conn.NextReader()

		if err != nil || messageType != websocket.BinaryMessage {
			return
		}
		HandlePacket(client, r)
	}
}

func (client *Client) sendToNearby(payload []byte, includeSelf bool) {
	client.simulation.ForEachCharacter(func(id uint32, _ *engine.Character) {
		if !includeSelf && id == client.id {
			return
		}
		if target, ok := Clients[id]; ok {
			target.send <- payload
		}
	})
}

func (client *Client) characterAttack(x float32, y float32, angle uint16) {
	item := client.character.GetHand()
	cooldown := client.character.AttackCooldown
	data := new(bytes.Buffer)
	data.WriteByte(uint8(ALLY_ATTACK))

	if cooldown < -0.1 {
		client.character.AttackCounter = 0
	}
	if cooldown <= 0.01 {
		counter := client.character.AttackCounter
		item := engine.GetItemData()[item]

		if item.Attacks != nil {
			attack := item.Attacks[counter]
			animation := attack.Animation
			projectiles := attack.Projectiles

			client.character.AttackCooldown = attack.Reload

			binary.Write(data, binary.LittleEndian, client.id)
			data.WriteByte(animation)

			binary.Write(data, binary.LittleEndian, uint16(len(projectiles)))

			for _, projectile := range projectiles {
				baseDamage := engine.GetProjectileData()[projectile.ID].Damage
				damage := baseDamage
				// we can freely scale damage up/down here based on whatever we want

				id := client.instance.CreateProjectile(projectile.ID, projectile.X+x, projectile.Y+y, (projectile.Angle+angle)%360, false, damage)
				proj := client.instance.Projectiles[id]
				packet := proj.Pack()

				binary.Write(data, binary.LittleEndian, id)
				data.Write(packet)

				client.simulation.AddProjectile(id, proj)
			}
		}
		client.sendToNearby(data.Bytes(), false)

		client.character.AttackCounter += 1
		client.character.AttackCounter %= uint8(len(item.Attacks))
	}
}
