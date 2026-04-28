const hotbar = document.getElementById("hotbar")

function populate_hotbar() {
  for (let i = 0; i < 6; i++) {
    const slot = document.createElement("div")
    slot.className = "slot"

    const sprite = document.createElement("button")
    sprite.className = "slot_sprite"
    sprite.style.backgroundPosition = "999px -10px"

    slot.onclick = () => {
      select_slot(i)
    }

    slot.appendChild(sprite)
    hotbar.appendChild(slot)
  }
}

function reset_hotbar() {
  for (let slot of hotbar.childNodes) {
    let sprite = slot.querySelector(".slot_sprite")
    sprite.style.backgroundPosition = "999px -10px"
    slot.className = "slot"
  }
}

function get_inventory_slot(slot) {
  if (slot < 6) {
    // hotbar nodes in inventory
  } else {
    // for inventory nodes
    return null
  }
}

function get_hotbar_slot(slot) {
  if (slot < 6) {
    return hotbar.querySelectorAll(".slot")[slot]
  }
  return null
}

function set_slot(slot_node, item_id) {
  let data = item_data[item_id]
  let sprite_node = slot_node.querySelector(".slot_sprite")

  sprite_node.style = `
    background-position: ${-data.sprite.x}px -${data.sprite.y}px
    `
}

function refresh_inventory(_inventory, hand) {
  inventory = _inventory

  reset_hotbar()
  hotbar.childNodes[hand].className = "slot selected"

  for (let slot of Object.keys(inventory)) {
    let node = get_inventory_slot(slot)

    if (slot < 6) {
      set_slot(get_hotbar_slot(slot), inventory[slot])
    }
  }
}

populate_hotbar()
