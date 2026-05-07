package packets

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

type Message struct {
	sender   uint32
	contents string
}

func (msg *Message) Pack() []byte {
	return []byte{}
}

func (client *Client) SendMessage(id uint32, sender string, msg string) {
	data := new(bytes.Buffer)

	data.WriteByte(uint8(CHAT_MESSAGE))
	binary.Write(data, binary.LittleEndian, id)

	n := len(sender)
	data.WriteByte(uint8(n))
	for i := 0; i < n; i++ {
		data.WriteByte(sender[i])
	}

	m := len(msg)
	data.WriteByte(uint8(m))
	for i := 0; i < m; i++ {
		data.WriteByte(msg[i])
	}

	client.send <- data.Bytes()
}

func (client *Client) ProcessMessage(msg string) {
	words := strings.Split(msg, " ")
	if len(msg) == 0 {
		return
	}

	if msg[0] == '/' {
		found := true

		// non admin
		switch words[0] {
		case "/tp":
			id, _ := strconv.Atoi(words[1])

			//using id for some reason. We should use username really
			client.SendMessage(1, "System", "Teleported to "+Clients[uint32(id)].username)
		default:
			found = false
		}

		// admin
		if client.admin && !found {
			found = true
			switch words[0] {
			case "/spawn":
				which, _ := strconv.Atoi(words[1])
				client.instance.SpawnNpc(uint8(which), client.character.GetX(), client.character.GetY())
				client.SendMessage(1, "System", "Spawned "+strconv.Itoa(which))
			default:
				found = false
			}
		}

		if !found {
			client.SendMessage(0, "System", "Invalid command")
		}
	} else {
		for _, recipient := range Clients {
			recipient.SendMessage(client.id, client.username, msg)
		}
	}
}
