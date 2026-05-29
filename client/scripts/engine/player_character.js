let character = null
let speed = 4
let reload = 1
let inventory = {}
let velocity = {
  x: 0,
  y: 0,
}

let mouse_x = 0
let mouse_y = 0
let attacking = false
let attack_cooldown = 0
let attack_counter = 0

function init_character(x, y, angle, health, hand, head, body, _inventory) {
  character = new Character(x, y, angle, health, hand, head, body)
  character.interpolator = null
  character.object.zIndex = 2
  init_combat()
  update_preview()

  const send_rate = 50
  let last_send = 0
  app.ticker.add(({ deltaMS, deltaTime }) => {
    if (!character) {
      return
    }
    if (attacking && attack_cooldown <= 0) {
      const rect = app.canvas.getBoundingClientRect()
      const x = mouse_x - rect.left + character.object.x * 8
      const y = mouse_y - rect.top + character.object.y * 8

      attack()
    }

    character.animator.tick(deltaMS)
    character.colorAnimator.tick(deltaMS)
    attack_cooldown -= deltaMS / 1000

    if (!chat_focused) {
      const len = Math.sqrt(velocity.x * velocity.x + velocity.y * velocity.y)
      const nx = len > 0 ? velocity.x / len : 0
      const ny = len > 0 ? velocity.y / len : 0

      character.object.x += (speed * nx * deltaTime) / 32
      character.object.y += (speed * ny * deltaTime) / 32
    }

    const mouse = app.renderer.events.pointer.global
    const centerX = app.screen.width / 2
    const centerY = app.screen.height / 2
    const dx = mouse.x - centerX
    const dy = mouse.y - centerY
    character.object.angle = Math.atan2(dy, dx) * (180 / Math.PI) + 90

    app.stage.pivot.set(character.object.x, character.object.y)
    app.stage.position.set(centerX, centerY)

    const now = Date.now()
    if (now - last_send >= send_rate) {
      send_position(
        character.object.x,
        character.object.y,
        character.object.angle,
      )
      last_send = now
    }
  })
}

const healthbar = document.getElementById("health")
const health_label = document.getElementById("health_label")
function update_healthbar(health, max_health) {
  healthbar.style.width = (100 * health) / max_health + "%"
  health_label.innerHTML = health
}

function attack() {
  const data = item_data[character.hand]

  if (data && (data.type == "head" || data.type == "body")) {
    // Holding a helmet/armor and clicking equips it into the matching gear
    // slot, same as dragging it there (like clicking to drink a potion).
    if (attack_cooldown < 0) {
      attack_cooldown = 0.5
      change_inventory(data.type == "head" ? 24 : 25, selected_slot)
    }
  } else if (data && data.on_use && attack_cooldown < 0) {
    attack_cooldown = 0.5
    send_attack(character.object.x, character.object.y, character.object.angle)
  } else if (data && data.attacks) {
    character.object.angle = ((character.object.angle % 360) + 360) % 360
    send_attack(character.object.x, character.object.y, character.object.angle)

    if (attack_cooldown < -0.1) {
      attack_counter = 0
    }

    const attacks = data.attacks
    const attack = attacks[attack_counter % data.attacks.length]

    attack_cooldown = attack.reload / reload

    const attack_projectiles = attack.projectiles ? attack.projectiles : []
    const attack_bombs = attack.bombs ? attack.bombs : []

    for (let proj of attack_projectiles) {
      const id = (Math.random() * 0xffffffff) >>> 0
      const projectile = new Projectile(
        proj.id,
        character.object.x,
        character.object.y,
        character.object.angle + proj.angle,
        true,
      )
      projectiles[id] = projectile
    }
    for (let bom of attack_bombs) {
      const id = (Math.random() * 0xffffffff) >>> 0
      const [target_x, target_y] = get_mouse_target()
      const bomb = new Bomb(
        bom.id,
        character.object.x,
        character.object.y,
        target_x,
        target_y,
        true,
      )
      bombs[id] = bomb
    }

    character.animator.animate(attack.animation, attack.reload)

    attack_counter += 1
    attack_counter %= attacks.length
  }
}

function init_combat() {
  app.canvas.addEventListener("pointerdown", (event) => {
    attacking = true
  })
  app.canvas.addEventListener("pointerup", (event) => {
    attacking = false
  })
  document.addEventListener("mousemove", (e) => {
    mouse_x = e.clientX
    mouse_y = e.clientY
  })
}

// Track which movement keys are physically held and derive velocity from that
// set every event. Keying off e.code (physical key) instead of e.key avoids the
// classic stuck-key bug where Shift/CapsLock change e.key between the keydown
// ("d") and keyup ("D"), so the keyup never cancels the keydown.
const pressed_keys = new Set()
const movement_keys = new Set(["KeyW", "KeyA", "KeyS", "KeyD"])

function recompute_velocity() {
  velocity.x =
    (pressed_keys.has("KeyD") ? 1 : 0) - (pressed_keys.has("KeyA") ? 1 : 0)
  velocity.y =
    (pressed_keys.has("KeyS") ? 1 : 0) - (pressed_keys.has("KeyW") ? 1 : 0)
}

function reset_keys() {
  pressed_keys.clear()
  recompute_velocity()
}

document.addEventListener("keydown", (e) => {
  if (e.repeat) return
  if (movement_keys.has(e.code)) {
    pressed_keys.add(e.code)
    recompute_velocity()
  }
})
document.addEventListener("keyup", (e) => {
  if (pressed_keys.delete(e.code)) {
    recompute_velocity()
  }
})

// If the window loses focus (alt-tab, OS shortcut, clicking out) the browser
// stops delivering keyup events, which would otherwise leave a key "stuck".
window.addEventListener("blur", reset_keys)
document.addEventListener("visibilitychange", () => {
  if (document.hidden) reset_keys()
})

function get_mouse_target() {
  const mouse = app.renderer.events.pointer.global
  const point = app.stage.toLocal(mouse)
  return [point.x, point.y]
}
