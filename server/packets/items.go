package packets

func (client *Client) selectSlot(idx uint8) {
	client.instance.SelectSlot(client.id, idx)
	client.sendCharacter()
}

func (client *Client) changeInventory(to uint8, from uint8) {
	client.instance.ChangeInventory(client.id, to, from)
	client.sendCharacter()
}
