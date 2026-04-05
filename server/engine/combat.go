package engine

func Hit(projectile *Projectile, target Entity) {
	target.Damage(projectile.damage)
	if !projectileData[projectile.id].Piercing {
		projectile.Dead = true
	}
}

func hitboxesIntersect(projectile *Projectile, entity Entity, npc bool) bool {
	projHitbox := float32(projectileData[projectile.id].Hitbox)
	entityHitbox := entity.GetHitbox()

	dx := projectile.x - entity.GetX()
	dy := projectile.y - entity.GetY()
	radius := projHitbox + entityHitbox
	return (dx*dx + dy*dy) <= radius*radius
}
