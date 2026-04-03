let character = null
let speed = 1
let inventory = {}
let velocity = {
    "x" : 0,
    "y" : 0,
}

let mouse_x = 0
let mouse_y = 0
let attacking = false
let attack_cooldown = 0
let attack_counter = 0

function init_character(x, y, angle, health, hand, head, body, _inventory) {
    character = new Character(x, y, angle, health, hand, head, body)
    character.object.zIndex = 2
    inventory = _inventory
    init_combat()
    
    const send_rate = 50
    let last_send = 0
    app.ticker.add(({deltaMS, deltaTime}) => {
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
        attack_cooldown -= deltaMS/1000
        
        character.object.x += speed * velocity.x * deltaTime / 32
        character.object.y += speed * velocity.y * deltaTime / 32

        const mouse = app.renderer.events.pointer.global
        const centerX = app.screen.width / 2
        const centerY = app.screen.height / 2
        const dx = mouse.x - centerX
        const dy = mouse.y - centerY
        character.object.angle = (Math.atan2(dy, dx) * (180 / Math.PI)) + 90

        app.stage.pivot.set(character.object.x, character.object.y)
        app.stage.position.set(centerX, centerY)

        const now = Date.now()
        if (now - last_send >= send_rate) {
            send_position(character.object.x, character.object.y, character.object.angle)
            last_send = now
        }
    })
}

function attack() {
    if (item_data[character.hand]) {
        const data = item_data[character.hand]

        send_attack(character.object.x, character.object.y, character.object.angle)

        if (attack_cooldown < -0.1) {
            attack_counter = 0
        }

        const attacks = data.attacks
        const attack = attacks[attack_counter]
        attack_cooldown = attack.reload

        for (let proj of attack.projectiles) {
            const id = (Math.random() * 0xFFFFFFFF) >>> 0
            const projectile = new Projectile(proj.id, character.object.x, character.object.y, character.object.angle + proj.angle)
            projectiles[id] = projectile
        }
        character.animator.animate(attack.animation)

        attack_counter += 1
        attack_counter %= attacks.length
    }
}

function init_combat() {
    app.canvas.addEventListener('pointerdown', (event) => {
        attacking = true
    })
    app.canvas.addEventListener('pointerup', (event) => {
        attacking = false
    })
    document.addEventListener("mousemove", (e) => {
        mouse_x = e.clientX
        mouse_y = e.clientY
    })
}

document.addEventListener("keydown", (e) => {
    if (e.repeat) return
    if (e.key == "d") {
        velocity.x += 1
    }
    if (e.key == "a") {
        velocity.x -= 1
    }
    if (e.key == "w") {
        velocity.y -= 1
    }
    if (e.key == "s") {
        velocity.y += 1
    }
})
document.addEventListener("keyup", (e) => {
    if (e.key == "d") {
        velocity.x -= 1
    }
    if (e.key == "a") {
        velocity.x += 1
    }
    if (e.key == "w") {
        velocity.y += 1
    }
    if (e.key == "s") {
        velocity.y -= 1
    }
})
