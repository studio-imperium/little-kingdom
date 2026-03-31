package packets

import (
	"fmt"
	"server/engine"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	id         uint32
	conn       *websocket.Conn
	send       chan []byte
	character  *engine.Character
	simulation *engine.Engine
}

func (client *Client) destroy() {
	fmt.Println("Client destroyed: ", client.id)
	client.conn.Close()
	delete(clients, client.id)
	engine.Game.RemoveCharacter(client.id)
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

var render_distance int = 16

func (client *Client) startSimulation() {
	for {
		if client.character.Dead {
			return
		}
		gameCharacters := make(map[uint32]struct{})
		gameNpcs := make(map[uint32]struct{})

		engine.Game.ForEachCharacter(func(id uint32, character *engine.Character) {
			gameCharacters[id] = struct{}{}
			if id != client.id && character.Dead {
				client.simulation.RemoveCharacter(id)
			} else if id != client.id && int(engine.Distance(character, client.character)) > render_distance {
				client.simulation.RemoveCharacter(id)
			}
		})
		engine.Game.ForEachNpc(func(id uint32, npc *engine.Npc) {
			gameNpcs[id] = struct{}{}
			if npc.Dead {
				client.simulation.RemoveNpc(id)
			} else if int(engine.Distance(npc, client.character)) > render_distance {
				npc.ExitView(client.id, client.character)
				client.simulation.RemoveNpc(id)
			}
		})

		engine.Game.ForEachCharacter(func(id uint32, character *engine.Character) {
			exists := client.simulation.HasCharacter(id)

			if id != client.id && !exists && int(engine.Distance(character, client.character)) <= render_distance {
				client.simulation.AddCharacter(id, character)
			}
		})
		engine.Game.ForEachNpc(func(id uint32, npc *engine.Npc) {
			exists := client.simulation.HasNpc(id)

			if !exists && int(engine.Distance(npc, client.character)) <= render_distance {
				npc.EnterView(client.id, client.character)
				client.simulation.AddNpc(id, npc)
			}
		})

		var removeCharacters []uint32
		client.simulation.ForEachCharacter(func(id uint32, character *engine.Character) {
			if _, exists := gameCharacters[id]; !exists {
				removeCharacters = append(removeCharacters, id)
			}
		})
		for _, id := range removeCharacters {
			client.simulation.RemoveCharacter(id)
		}

		var removeNpcs []uint32
		client.simulation.ForEachNpc(func(id uint32, npc *engine.Npc) {
			if _, exists := gameNpcs[id]; !exists {
				removeNpcs = append(removeNpcs, id)
			}
		})
		for _, id := range removeNpcs {
			client.simulation.RemoveNpc(id)
		}

		time.Sleep(time.Second)
	}
}
