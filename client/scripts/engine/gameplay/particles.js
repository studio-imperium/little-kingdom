const particles = {}

class ParticleEmitter {
  constructor(object, texture, delta) {
    this.object = object
    this.container = new PIXI.ParticleContainer({
      dynamicProperties: {
        vertex: true,
        position: true,
        color: true,
      },
    })
    this.particles = []
    this.time_since_last = 0
    this.delta = delta
    this.texture = texture

    app.stage.addChild(this.container)
    head_layer.attach(this.container)
  }

  radiate(max_radius, angles, radius = 1) {
    if (radius <= max_radius) {
      let angle = 2 * Math.PI
      let delta = angle / angles

      for (let i = 0; i < angles; i++) {
        angle -= delta
        let angledX = this.object.x + Math.cos(angle) * radius
        let angledY = this.object.y + Math.sin(angle) * radius

        this.add_particle(angledX, angledY, 0.8)
      }
      setTimeout(() => {
        this.radiate(max_radius, angles, radius + 0.7)
      }, 10)
    }
  }

  tick(deltaMS, emit = true) {
    this.time_since_last += deltaMS

    if (emit && this.time_since_last > this.delta * 1000) {
      this.add_particle(this.object.x - 2.5 / 64, this.object.y - 2.5 / 64, 1)
      this.time_since_last = 0
    }

    for (let i = this.particles.length - 1; i >= 0; i--) {
      let particle = this.particles[i]
      particle.elapsed_time += deltaMS
      particle.object.scaleX -= deltaMS / 64000
      particle.object.scaleY -= deltaMS / 64000
      particle.object.alpha -= deltaMS / 500

      if (particle.elapsed_time > particle.lifetime * 1000) {
        this.container.removeParticle(particle.object)
        this.particles.splice(i, 1)
      }
    }
  }

  add_particle(x, y, size) {
    let particle = new Particle(this.texture, x, y, size)

    this.particles.push(particle)
    this.container.addParticle(particle.object)
  }

  kill() {
    app.stage.removeChild(this.container)
  }
}

class Particle {
  constructor(texture, x, y, size) {
    this.object = build_particle(texture, size)
    this.lifetime = size
    this.elapsed_time = 0
    this.object.x = x - 2.5 / 8
    this.object.y = y - 2.5 / 8
  }
}

function build_particle(texture_id, size) {
  if (!cache[texture_id]) {
    cache[texture_id] = create_texture(textures[texture_id])
  }

  let particle = new PIXI.Particle({
    texture: cache[texture_id],
    x: 0,
    y: 0,
    scaleX: size / 64,
    scaleY: size / 64,
  })
  return particle
}
