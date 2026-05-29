package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"server/engine"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	id              uint32
	admin           bool
	username        string
	conn            *websocket.Conn
	send            chan []byte
	character       *engine.Character
	instance        *engine.Engine
	simulation      *engine.Engine
	discoveredCells map[uint16]*engine.Cell
	closeOnce       sync.Once
}

// cleanup tears the client down exactly once: it removes the client from the
// global registry, marks its character dead (which stops StartSimulation and
// UpdateCells on their next tick) and closes the socket (which unblocks the
// read and write pumps). Without this, a write error used to leave the client
// in Clients with nothing draining its send channel — the channel filled, and
// PropogateWorldState dropped every world state forever, so the player stayed
// connected with terrain loaded but never saw another entity again.
func (client *Client) cleanup() {
	client.closeOnce.Do(func() {
		fmt.Println("Client destroyed: ", client.id)

		clientsMu.Lock()
		delete(Clients, client.id)
		clientsMu.Unlock()

		client.instance.RemoveCharacter(client.id, true)
		client.conn.Close()
	})
}

func (client *Client) destroy() {
	client.cleanup()
}

var guests int = 0

func CreateClient(conn *websocket.Conn) {
	guests += 1
	client := Client{
		id:              0,
		admin:           true,
		username:        "Guest_" + strconv.Itoa(guests),
		conn:            conn,
		send:            make(chan []byte, 64),
		character:       nil,
		simulation:      engine.CreateEngine(),
		discoveredCells: make(map[uint16]*engine.Cell),
	}
	go client.recievePackets()
	go client.writePackets()
}

func (client *Client) writePackets() {
	conn := client.conn
	for {
		message, ok := <-client.send
		if !ok {
			return
		}

		// A bare 1-byte {11} is the internal "this player died" sentinel and
		// triggers teardown. A longer message starting with 11 is a real
		// CHARACTER_DEAD wire packet (e.g. an npc death, [11, id]) and must be
		// forwarded to the client normally, not treated as a disconnect.
		if len(message) == 1 && message[0] == 11 {
			for _, recipient := range snapshotClients() {
				recipient.SendMessage(0, "System", client.username+" has died!")
			}
			client.cleanup()
			return
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
			client.cleanup()
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
		clientsMu.RLock()
		target, ok := Clients[id]
		clientsMu.RUnlock()
		if ok {
			// Non-blocking: never stall a sender on a backed-up peer. A dropped
			// ally-attack relay just costs a missed animation; the world state
			// resyncs it. Blocking here used to wedge whole goroutines.
			select {
			case target.send <- payload:
			default:
			}
		}
	})
}

func (client *Client) characterAttack(x float32, y float32, targetX float32, targetY float32, angle uint16) {
	idx := client.instance.GetHand(client.id)
	item_id := client.instance.GetSlot(client.id, idx)
	item := engine.GetItemData(item_id)

	cooldown := client.character.AttackCooldown
	data := new(bytes.Buffer)
	data.WriteByte(uint8(ALLY_ATTACK))

	if item.OnUse != "" {
		client.instance.UseItem(client.id, item.OnUse)
		return
	}

	if cooldown < -0.1 {
		client.character.AttackCounter = 0
	}
	if cooldown <= 0.1 {
		counter := client.character.AttackCounter

		if item.Attacks != nil {
			attack := item.Attacks[int(counter)%len(item.Attacks)]
			animation := attack.Animation
			projectiles := attack.Projectiles
			bombs := attack.Bombs

			// can make this longer/shorter here
			reload := attack.Reload / client.character.Reload

			client.character.AttackCooldown = reload

			binary.Write(data, binary.LittleEndian, client.id)
			data.WriteByte(animation)
			binary.Write(data, binary.LittleEndian, uint16(reload*1000))

			binary.Write(data, binary.LittleEndian, uint16(len(projectiles)))

			for _, projectile := range projectiles {
				baseDamage := engine.GetProjectileData(projectile.ID).Damage

				damage := float32(baseDamage) * client.character.Power
				// we can freely scale damage up/down here based on whatever we want

				id := client.instance.CreateProjectile(projectile.ID, projectile.X+x, projectile.Y+y, (projectile.Angle+angle)%360, false, damage)
				proj := client.instance.Projectiles[id]
				packet := proj.Pack()

				binary.Write(data, binary.LittleEndian, id)
				data.Write(packet)

				client.simulation.AddProjectile(id, proj)
			}

			binary.Write(data, binary.LittleEndian, uint16(len(bombs)))

			for _, bomb := range bombs {
				baseDamage := engine.GetBombData(bomb.ID).Damage
				damage := baseDamage
				timer := engine.GetBombData(bomb.ID).Airtime
				// we can freely scale damage up/down here based on whatever we want
				// and scale timer

				id := client.instance.CreateBomb(bomb.ID, targetX, targetY, client.character, false, damage, timer)
				bomb := client.instance.Bombs[id]
				packet := bomb.Pack()

				binary.Write(data, binary.LittleEndian, id)
				data.Write(packet)

				client.simulation.AddBomb(id, bomb)
			}

			client.character.AttackCounter += 1
			client.character.AttackCounter %= uint8(len(item.Attacks))
			client.sendToNearby(data.Bytes(), false)
		}
	}
}

func (client *Client) sendCharacter() {
	data := client.instance.PackCharacter(client.id, byte(HANDSHAKE))
	client.send <- data
}
