package engine

func Hit(projectile *Projectile, target Entity) {
	target.Damage(projectile.damage)
	if !projectileData[projectile.id].Piercing {
		projectile.Dead = true
	}
}

func hitboxesIntersect(projectile *Projectile, npc *Npc) bool {
	projHitbox := float32(projectileData[projectile.id].Hitbox)
	npcHitbox := float32(npcData[npc.id].Hitbox)

	dx := projectile.x - npc.x
	dy := projectile.y - npc.y
	radius := projHitbox + npcHitbox
	return (dx*dx + dy*dy) <= radius*radius
}
