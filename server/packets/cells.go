package packets

import (
	"bytes"
	"encoding/binary"
	"server/engine"
	"time"
)

func (client *Client) sendTiles(cell *engine.Cell) {
	tiles := cell.Tiles
	data := new(bytes.Buffer)

	data.WriteByte(byte(TILES))
	binary.Write(data, binary.LittleEndian, int32(cell.Origin.X))
	binary.Write(data, binary.LittleEndian, int32(cell.Origin.Y))
	binary.Write(data, binary.LittleEndian, uint16(len(tiles)))

	for _, tile := range tiles {
		binary.Write(data, binary.LittleEndian, int32(tile.X))
		binary.Write(data, binary.LittleEndian, int32(tile.Y))
		data.WriteByte(byte(uint8(tile.Val)))
	}

	client.send <- data.Bytes()
}

func (client *Client) sendCell() {
	client.instance.Map.Mu.Lock()

	nearest := client.instance.Map.GetNearestCell(client.character)

	for _, cell := range append(nearest.GetAdjacentCells(), nearest) {
		for _, adj_cell := range cell.GetAdjacentCells() {
			cell.Characters[client.id] = client.character
			if _, okay := client.discoveredCells[adj_cell.Idx]; !okay {
				cell.Load()
				client.discoveredCells[adj_cell.Idx] = adj_cell
				client.sendTiles(adj_cell)
			}
		}
	}

	client.instance.Map.Mu.Unlock()
}

func (client *Client) UpdateCells() {
	delta := time.Millisecond * 500
	ticker := time.NewTicker(delta)

	for {
		<-ticker.C
		if client.character.Dead {
			return
		}

		client.sendCell()
	}
}
