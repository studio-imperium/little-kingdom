package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"server/atlas"
	"time"
)

func (client *Client) sendTiles(cell *atlas.Cell) {
	tiles := cell.Tiles
	data := new(bytes.Buffer)

	data.WriteByte(byte(TILES))
	binary.Write(data, binary.LittleEndian, int32(cell.Origin.X))
	binary.Write(data, binary.LittleEndian, int32(cell.Origin.Y))
	binary.Write(data, binary.LittleEndian, uint16(len(tiles)))

	for _, tile := range tiles {
		binary.Write(data, binary.LittleEndian, int32(tile.X))
		binary.Write(data, binary.LittleEndian, int32(tile.Y))
		data.WriteByte(byte(uint8(tile.Value)))
	}

	client.send <- data.Bytes()
	fmt.Println("Sending chunk")
}

func (client *Client) UpdateCells() {
	delta := time.Millisecond * 50
	ticker := time.NewTicker(delta)

	for {
		if client.character.Dead {
			return
		}

		pt := atlas.Point{
			X: float64(client.character.GetX()),
			Y: float64(client.character.GetY()),
		}
		nearest := client.instance.World.GetNearestCell(pt)

		for _, cell := range append(nearest.GetAdjacentCells(), nearest) {
			for _, adj_cell := range cell.GetAdjacentCells() {
				if _, okay := client.discoveredCells[adj_cell.Origin]; !okay {
					client.discoveredCells[adj_cell.Origin] = adj_cell
					client.sendTiles(adj_cell)
				}
			}
		}
		<-ticker.C
	}
}
