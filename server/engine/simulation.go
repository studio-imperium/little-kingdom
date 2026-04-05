package engine

import "time"

var render_distance int = 16

func (simulation *Engine) StartSimulation(clientID uint32, instance *Engine, clientCharacter *Character) {
	delta := time.Millisecond * 50
	ticker := time.NewTicker(delta)
	for {
		if clientCharacter.Dead {
			return
		} else {
			clientCharacter.Tick(float32(delta) / float32(time.Second))
		}

		// cull characters
		simulation.ForEachCharacter(func(id uint32, character *Character) {
			if id == clientID {
				return
			}
			if character.Dead || int(Distance(character, clientCharacter)) > render_distance {
				simulation.RemoveCharacter(id, false)
			}
		})
		// cull npcs
		simulation.ForEachNpc(func(id uint32, npc *Npc) {
			if npc.Dead {
				delete(simulation.Npcs, id)
				return
			}
			if int(Distance(npc, clientCharacter)) > render_distance {
				npc.ExitView(clientID, clientCharacter)
				delete(simulation.Npcs, id)
			}
		})

		// add characters
		instance.ForEachCharacter(func(id uint32, character *Character) {
			exists := simulation.HasCharacter(id)
			if !exists && int(Distance(character, clientCharacter)) <= render_distance {
				simulation.AddCharacter(id, character)
			}
		})
		// add npcs
		instance.ForEachNpc(func(id uint32, npc *Npc) {
			exists := simulation.HasNpc(id)
			if !exists && int(Distance(npc, clientCharacter)) <= render_distance {
				npc.EnterView(clientID, clientCharacter)
				simulation.AddNpc(id, npc)
			}
		})

		hits := make([]struct {
			projectile *Projectile
			target     Entity
		}, 0)

		simulation.mu.Lock()
		for npcID, npc := range simulation.Npcs {
			_ = npcID
			_ = npc
		}
		for id, projectile := range simulation.Projectiles {
			if projectile.Dead {
				delete(simulation.Projectiles, id)
			}
			if projectile.evil {
				if _, hit := projectile.hitlist[clientID]; hit || projectile.Dead {
					continue
				}
				if hitboxesIntersect(projectile, clientCharacter, false) {
					projectile.hitlist[clientID] = clientCharacter
					hits = append(hits, struct {
						projectile *Projectile
						target     Entity
					}{projectile: projectile, target: clientCharacter})
				}
			} else {
				for npcID, npc := range simulation.Npcs {
					if _, hit := projectile.hitlist[npcID]; hit || projectile.Dead {
						continue
					}
					if hitboxesIntersect(projectile, npc, true) {
						projectile.hitlist[npcID] = npc
						hits = append(hits, struct {
							projectile *Projectile
							target     Entity
						}{projectile: projectile, target: npc})
					}
				}
			}
		}
		simulation.mu.Unlock()

		for _, hit := range hits {
			Hit(hit.projectile, hit.target)
		}

		<-ticker.C
	}
}
