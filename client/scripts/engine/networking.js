const addr = "localhost"
let CONNECTED = false
let socket
let token

const [
  HANDSHAKE,
  CHARACTER_POSITION,
  CHARACTER_ATTACK,
  RECIEVE_ATTACK,
  WORLD_STATE,
  DAMAGED,
  TILES,
  SELECT_SLOT,
] = [0, 1, 2, 3, 4, 5, 6, 7]

function handshake() {
  token = (Math.random() * 0x100000000) >>> 0
  const buffer = new ArrayBuffer(5)
  const data = new DataView(buffer)

  data.setUint8(0, HANDSHAKE)
  data.setUint32(1, token, true)

  socket.send(data)
}

let initialized_character = false
function set_character(data) {
  const [x, y, angle, health, hand, head, body] = [
    data.getFloat32(1, true),
    data.getFloat32(5, true),
    data.getUint16(9, true),
    data.getUint16(11, true),
    data.getUint8(13),
    data.getUint8(14),
    data.getUint8(15),
  ]
  const inventory = {}
  const slots = data.getUint8(16)
  for (let i = 0; i < slots; i++) {
    let offset = i * 2
    let slot = data.getUint8(17 + offset)
    let item = data.getUint8(18 + offset)
    inventory[slot] = item
  }

  init_character(x, y, angle, health, hand, head, body, inventory)
  start_engine()
}

function set_world(data) {
  let offset = 1
  const character_count = data.getUint16(offset, true)
  offset += 2

  for (let i = 0; i < character_count; i++) {
    let id = data.getUint32(offset, true)
    let x = data.getFloat32(4 + offset, true)
    let y = data.getFloat32(8 + offset, true)
    let angle = data.getUint16(12 + offset, true)
    let health = data.getUint16(14 + offset, true)
    let hand = data.getUint8(16 + offset)
    let head = data.getUint8(17 + offset)
    let body = data.getUint8(18 + offset)

    if (id == token) {
    } else if (characters[id]) {
      characters[id].update(x, y, angle, health, hand, head, body)
    } else {
      characters[id] = new Character(x, y, angle, health, hand, head, body)
    }
    offset += 19
  }

  const npc_count = data.getUint16(offset, true)
  offset += 2

  for (let i = 0; i < npc_count; i++) {
    let id = data.getUint32(offset, true)
    let which = data.getUint8(4 + offset)
    let x = data.getFloat32(5 + offset, true)
    let y = data.getFloat32(9 + offset, true)
    let health = data.getUint16(13 + offset, true)

    offset += 15

    if (npcs[id]) {
      npcs[id].update(x, y, health)
    } else {
      npcs[id] = new Npc(which, x, y, health)
    }
  }
}

function set_attack(data) {
  let offset = 1
  const id = data.getUint32(offset, true)
  offset += 4
  const animation = data.getUint8(offset)
  offset += 1
  const projectile_count = data.getUint16(offset, true)
  offset += 2

  if (characters[id]) {
    characters[id].animator.animate(animation)
  }
  if (npcs[id]) {
    npcs[id].animator.animate(animation)
  }

  for (let i = 0; i < projectile_count; i++) {
    const projectile_id = data.getUint32(offset, true)
    offset += 4
    const which = data.getUint8(offset)
    offset += 1
    const x = data.getFloat32(offset, true)
    offset += 4
    const y = data.getFloat32(offset, true)
    offset += 4
    const angle = data.getUint16(offset, true)
    offset += 2

    const projectile = new Projectile(which, x, y, angle)
    console.log(x, y, angle)
    projectiles[projectile_id] = projectile
  }
}

function set_tiles(data) {
  let offset = 1
  const origin_x = data.getInt32(offset, true)
  const origin_y = data.getInt32(4 + offset, true)
  const tile_count = data.getUint16(8 + offset, true)
  offset += 10

  for (let i = 0; i < tile_count; i++) {
    let x = data.getInt32(offset, true)
    let y = data.getInt32(4 + offset, true)
    let tile_id = data.getUint8(8 + offset)

    tiles[y * size + x] = tile_id
    offset += 9
  }
}

function damaged(data) {
  let offset = 1
  const id = data.getUint32(offset, true)

  if (characters[id]) {
    characters[id].damage()
  }
  if (npcs[id]) {
    npcs[id].damage()
  }
}

function send_position(x, y, angle) {
  const buffer = new ArrayBuffer(11)
  const data = new DataView(buffer)

  data.setUint8(0, CHARACTER_POSITION)
  data.setFloat32(1, x, true)
  data.setFloat32(5, y, true)
  data.setUint16(9, angle, true)

  socket.send(data)
}

function send_attack(x, y, angle) {
  const buffer = new ArrayBuffer(11)
  const data = new DataView(buffer)

  data.setUint8(0, CHARACTER_ATTACK)
  data.setFloat32(1, x, true)
  data.setFloat32(5, y, true)
  data.setFloat32(9, target_x, true)
  data.setFloat32(13, target_y, true)
  data.setUint16(17, angle, true)

  socket.send(data)
}

function select_slot(idx) {
  const buffer = new ArrayBuffer(3)
  const data = new DataView(buffer)

  data.setUint8(0, SELECT_SLOT)
  data.setUint(1, idx, true)

  socket.send(data)
}

function connect() {
  socket = new WebSocket("ws://" + addr + ":8082/connect")
  socket.binaryType = "arraybuffer"

  function open() {
    CONNECTED = true
    handshake()
  }

  function closed(e) {
    console.log(e)
    CONNECTED = false
  }

  function handle_packet(e) {
    const data = new DataView(e.data)
    const packet_type = data.getUint8(0)

    switch (packet_type) {
      case HANDSHAKE:
        set_character(data)
        break
      case WORLD_STATE:
        set_world(data)
        break
      case RECIEVE_ATTACK:
        set_attack(data)
        break
      case DAMAGED:
        damaged(data)
        break
      case TILES:
        set_tiles(data)
        break
      default:
        console.log("Bad packet recieved: ", packet_type)
    }
  }

  socket.addEventListener("open", open)
  socket.addEventListener("close", closed)
  socket.addEventListener("error", closed)
  socket.addEventListener("message", handle_packet)
}
