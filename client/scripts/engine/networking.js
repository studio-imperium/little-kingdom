function initiate() {
    init_character({
        hand : "wind_scroll",
        gear : {
            "head" : "iron_helmet",
            "body" : "iron_platemail",
        },
        x: 0,
        y: 0,
        angle: 0,
        inventory : {}
    })
    new Character("shortsword", {
            "head" : null,
            "body" : null,
        }, 30, 32, 32)
    new Character("shortbow", {
            "head" : "iron_helmet",
            "body" : null,
        }, 30, -32, 32)
    new Npc("basilisk", 150, -32, -32)
    new Npc("sheep", 250, -64, -32)
}