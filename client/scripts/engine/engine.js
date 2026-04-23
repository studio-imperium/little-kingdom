const characters = {}
const npcs = {}
const projectiles = {}

const render_dist = 32
const size = 1000
const tiles = new Uint8Array(size * size)
const added = new Uint8Array(size * size)

function outside_range(x, y) {
  let player_x = character.object.x
  let player_y = character.object.y
  let halfdist = render_dist / 2

  return (
    x < player_x - halfdist ||
    x >= player_x + halfdist ||
    y < player_y - halfdist ||
    y >= player_y + halfdist
  )
}
function distance(x1, y1, x2, y2) {
  return Math.sqrt(Math.pow(x1 - x2, 2) + Math.pow(y1 - y2, 2))
}

let last_x
let last_y
function diff(x, y) {
  let difference = Math.abs(last_x - x) + Math.abs(last_y - y)
  last_x = x
  last_y = y
  return difference
}

let timer = 0
let elapsed = 0
function start_engine() {
  app.ticker.add(({ deltaMS, deltaTime }) => {
    timer += deltaMS
    for (let id of Object.keys(characters)) {
      let { interpolator, animator, colorAnimator, object } = characters[id]
      interpolator.tick(deltaMS)
      animator.tick(deltaMS)
      if (colorAnimator) {
        colorAnimator.tick(deltaMS)
      }

      if (Date.now() > interpolator.last_frame + 400) {
        characters[id].kill(id)
      }
    }
    for (let id of Object.keys(npcs)) {
      let { interpolator, animator, colorAnimator, object } = npcs[id]
      interpolator.tick(deltaMS)
      animator.tick(deltaMS)
      if (colorAnimator) {
        colorAnimator.tick(deltaMS)
      }

      if (Date.now() > interpolator.last_frame + 400) {
        npcs[id].kill(id)
      }
    }
    for (let id of Object.keys(projectiles)) {
      const { object, which } = projectiles[id]
      const speed = projectile_data[which].speed
      const rad = (object.angle - 90) * (Math.PI / 180)
      const dx = Math.cos(rad)
      const dy = Math.sin(rad)
      object.x += (dx * speed * deltaTime) / 16
      object.y += (dy * speed * deltaTime) / 16

      object.alpha += speed / 32
      if (object.scale.x * 64 < 1) {
        object.scale.set(object.scale.x + 1 / 1500)
      } else {
        object.scale.set(1 / 64)
      }
    }

    if (timer >= 500) {
      timer = 0
      let player_x = Math.floor(character.object.x)
      let player_y = Math.floor(character.object.y)
      let halfdist = render_dist / 2

      if (diff(player_x, player_y) < 2) {
        return
      }

      for (let tile_id of Object.keys(tile_map)) {
        let tile = tile_map[tile_id]
        let x = tile.x
        let y = tile.y

        if (outside_range(x, y) || tiles[y * size + x] != tile.idx) {
          added[y * size + x] = false
          tile.mesh.destroy()
          delete tile_map[tile_id]
        }
      }

      for (let y = player_y - halfdist; y < player_y + halfdist; y++) {
        for (let x = player_x - halfdist; x < player_x + halfdist; x++) {
          let tile_id = tiles[y * size + x]
          let is_added = added[y * size + x]
          if (!is_added && tile_id >= 0) {
            add_tile(x, y, tile_id)
            added[y * size + x] = true
          }
        }
      }
      render_tiles()
    }

    elapsed += deltaMS / 1000
    const offsetX = Math.sin(elapsed * 0.8) * 0.5
    const offsetY = Math.cos(elapsed * 0.6) * 0.5

    for (const { tile, mesh, uvs } of Object.values(tile_map)) {
      if (tile == "lava" || tile == "water") {
        const uv_buffer = mesh.geometry.getBuffer("aUV")

        for (let i = 0; i < uv_buffer.data.length; i += 2) {
          uv_buffer.data[i] = uvs.modified[i] + offsetX
          uv_buffer.data[i + 1] = uvs.modified[i + 1] + offsetY
        }
        uv_buffer.update()
      }
    }
  })
}

class Cell {
  constructor(origin_x, origin_y) {
    this.visible = false
    this.tiles = []
    this.origin = [origin_x, origin_y]
    cells.push(this)
  }

  add_tile(x, y, tile_idx) {
    this.tiles.push([x, y, tile_idx])
  }

  hide() {
    if (this.visible) {
      this.visible = false
      for (let tile of this.tiles) {
        let x = tile[0]
        let y = tile[1]
        let data = tile_map[`${x},${y}`]

        if (data) {
          data.mesh.destroy()
          delete tile_map[`${x},${y}`]
        }
      }
    }
  }
  show() {
    if (!this.visible) {
      this.visible = true
      let player_x = character.object.x
      let player_y = character.object.y

      for (let tile of this.tiles) {
        let x = tile[0]
        let y = tile[1]
        let tile_id = tile[2]

        if (distance(x, y, player_x, player_y) < 16) {
        }
        add_tile(x, y, tile_id)
      }
    }
  }
}

class Interpolator {
  constructor(object) {
    this.object = object
    this.frames = []
    this.last_frame = Date.now()
  }

  lerp_angle_degrees(a, b, t) {
    d1 * frame1[2] + d2 * frame2[2]
    const delta = ((b - a + 540) % 360) - 180
    return a + delta * t
  }

  tick(delta_ms) {
    const time = Date.now()

    while (this.frames.length > 2) {
      this.frames.shift()
    }
    if (this.frames.length > 1) {
      let frame1 = this.frames[0]
      let frame2 = this.frames[1]
      let diff = frame2[0] - frame1[0]

      const d1 = (frame2[0] - time) / diff
      const d2 = (time - frame1[0]) / diff

      // npcs
      if (frame1.length == 3 || frame2.length == 3) {
        const dx = frame2[1] - frame1[1]
        const dy = frame2[2] - frame1[2]

        if (frame1[3] === undefined) {
          frame1[3] = this.object.angle
        }

        if (dx !== 0 || dy !== 0) {
          frame2[3] = Math.atan2(dy, dx) * (180 / Math.PI) + 90
        } else if (frame2[3] === undefined) {
          frame2[3] = frame1[3]
        }
      }

      if (frame1[1] != frame2[1]) {
        this.object.x = d1 * frame1[1] + d2 * frame2[1]
      }
      if (frame1[2] != frame2[2]) {
        this.object.y = d1 * frame1[2] + d2 * frame2[2]
      }
      if (frame1[3] != frame2[3]) {
        if (frame1[3] - frame2[3] > 180) {
          this.object.angle = d1 * frame1[3] + d2 * (frame2[3] + 360)
        } else if (frame2[3] - frame1[3] > 180) {
          this.object.angle = d1 * (frame1[3] + 360) + d2 * frame2[3]
        } else {
          this.object.angle = d1 * (frame1[3] + 360) + d2 * (frame2[3] + 360)
        }
      }
    }
  }

  add_char_frame(x, y, angle) {
    this.last_frame = Date.now()
    this.frames.push([this.last_frame + 200, x, y, angle % 360])
  }
  add_npc_frame(x, y) {
    this.last_frame = Date.now()
    this.frames.push([this.last_frame + 200, x, y])
  }
}

class Animator {
  constructor(object) {
    this.object = object
    this.head = object.getChildByName("head")
    this.body = object.getChildByName("body")
    this.hand = object.getChildByName("hand")
    this.animation = null
    this.timestamp = 0
  }

  tick(delta_ms) {
    if (this.animation != null) {
      this.timestamp += delta_ms
      const frames = this.animation.frames || []

      if (frames.length === 0) {
        this.animation = null
        this.timestamp = 0
        return
      }

      const now = this.timestamp
      let frame1 = frames[0]
      let frame2 = null

      for (let i = 0; i < frames.length; i++) {
        const frame = frames[i]
        if (frame.time <= now) {
          frame1 = frame
        } else {
          frame2 = frame
          break
        }
      }

      const applyFrame = (frame) => {
        if (this.head && frame.head_angle !== undefined) {
          this.head.angle = frame.head_angle
        }
        if (this.body && frame.body_angle !== undefined) {
          this.body.angle = frame.body_angle
        }
        if (this.hand && frame.hand_angle !== undefined) {
          this.hand.angle = frame.hand_angle
        }
        if (this.hand && frame.hand_scale !== undefined) {
          if (this.hand.scale && typeof this.hand.scale.set === "function") {
            this.hand.scale.set(frame.hand_scale)
          } else {
            this.hand.scale = frame.hand_scale
          }
        }
        if (this.object && frame.object_scale !== undefined) {
          if (
            this.object.scale &&
            typeof this.object.scale.set === "function"
          ) {
            this.object.scale.set(frame.object_scale / 64)
          } else {
            this.object.scale = frame.object_scale / 64
          }
        }
      }

      if (!frame2) {
        applyFrame(frame1)
        this.animation = null
        this.timestamp = 0
        return
      }

      const diff = frame2.time - frame1.time
      const t = diff > 0 ? (now - frame1.time) / diff : 1
      const lerp = (a, b) => a + (b - a) * t

      applyFrame({
        object_scale:
          frame1.object_scale !== undefined && frame2.object_scale !== undefined
            ? lerp(frame1.object_scale, frame2.object_scale)
            : frame1.object_scale,
        head_angle:
          frame1.head_angle !== undefined && frame2.head_angle !== undefined
            ? lerp(frame1.head_angle, frame2.head_angle)
            : frame1.head_angle,
        body_angle:
          frame1.body_angle !== undefined && frame2.body_angle !== undefined
            ? lerp(frame1.body_angle, frame2.body_angle)
            : frame1.body_angle,
        hand_angle:
          frame1.hand_angle !== undefined && frame2.hand_angle !== undefined
            ? lerp(frame1.hand_angle, frame2.hand_angle)
            : frame1.hand_angle,
        hand_scale:
          frame1.hand_scale !== undefined && frame2.hand_scale !== undefined
            ? lerp(frame1.hand_scale, frame2.hand_scale)
            : frame1.hand_scale,
      })
    }
  }

  animate(which) {
    this.animation = animation_data[which]
    this.timestamp = 0
  }
}

class ColorAnimator {
  constructor(object) {
    this.object = object
    this.active = false
    this.duration = 0
    this.elapsed = 0
    this.startR = 255
    this.startG = 255
    this.startB = 255
  }

  tick(delta_ms) {
    if (!this.active || !this.object) {
      return
    }

    this.elapsed += delta_ms
    const duration = this.duration
    const t = duration > 0 ? Math.min(this.elapsed / duration, 1) : 1

    const r = Math.round(this.startR + (255 - this.startR) * t)
    const g = Math.round(this.startG + (255 - this.startG) * t)
    const b = Math.round(this.startB + (255 - this.startB) * t)

    this.object.tint = (r << 16) | (g << 8) | b

    if (t >= 1) {
      this.object.tint = 0xffffff
      this.active = false
    }
  }

  animate(color, duration_ms) {
    const value = color >>> 0
    this.startR = (value >> 16) & 0xff
    this.startG = (value >> 8) & 0xff
    this.startB = value & 0xff
    this.duration = Math.max(0, duration_ms || 0)
    this.elapsed = 0
    this.active = true

    if (this.object) {
      this.object.tint = value
    }

    if (this.duration === 0) {
      if (this.object) {
        this.object.tint = 0xffffff
      }
      this.active = false
    }
  }
}

class Projectile {
  constructor(id, x, y, angle) {
    this.object = build_projectile(id)
    this.which = id
    this.object.x = x
    this.object.y = y
    this.object.angle = angle
    this.object.scale.set(0)

    add_object(this.object)
  }
}

class Character {
  constructor(x, y, angle, health, hand, head, body) {
    this.object = build_character(hand, head, body)
    this.kit = hand + head + body
    this.object.x = x
    this.object.y = y
    this.object.angle = angle
    this.object.health = health
    this.hand = hand
    this.head = head
    this.body = body
    this.interpolator = new Interpolator(this.object)
    this.animator = new Animator(this.object)
    this.colorAnimator = new ColorAnimator(this.object)

    add_object(this.object)
  }

  update(x, y, angle, health, hand, head, body) {
    if (hand + head + body != this.kit) {
      this.object.destroy()
      this.object = build_character(hand, head, body)
      this.kit = hand + head + body
      this.object.x = x
      this.object.y = y
      this.object.angle = angle
      this.interpolator.object = this.object
      this.colorAnimator.object = this.object
      add_object(this.object)
    }
    this.object.health = health
    this.hand = hand
    this.head = head
    this.body = body
    this.interpolator.add_char_frame(x, y, angle)
  }

  damage() {
    this.colorAnimator.animate(0xffb3b3, 300)
  }

  kill(id) {
    this.colorAnimator.animate(0xff0000, 300)
    this.animator.animate(0)
    this.interpolator.frames = []
    this.interpolator.last_frame = Date.now()

    setTimeout(() => {
      this.object.destroy()
      delete characters[id]
    }, 300)
  }
}

class Npc {
  constructor(id, x, y, health) {
    this.object = build_npc(id)
    this.object.angle = 0
    this.object.x = x
    this.object.y = y
    this.interpolator = new Interpolator(this.object)
    this.animator = new Animator(this.object)
    this.colorAnimator = new ColorAnimator(this.object)

    add_object(this.object)
  }

  update(x, y, health) {
    this.object.health = health
    this.interpolator.add_npc_frame(x, y)
  }

  damage() {
    this.colorAnimator.animate(0xffb3b3, 300)
  }

  kill(id) {
    this.colorAnimator.animate(0xffb3b3, 30000)
    this.animator.animate(0)
    this.interpolator.frames = []
    this.interpolator.last_frame = Date.now()

    setTimeout(() => {
      this.object.destroy()
      delete npcs[id]
    }, 300)
  }
}
