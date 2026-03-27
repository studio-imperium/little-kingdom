let character = null

function init_character(data) {
    const {hand, gear, angle, x, y} = data
    character = new Character(hand, gear, angle, x, y)
}

document.addEventListener("keydown", (e) => {
    if (e.key == "d") {
        character.object.x += 5
    }
    if (e.key == "a") {
        character.object.x -= 5
    }
    if (e.key == "w") {
        character.object.y -= 5
    }
    if (e.key == "s") {
        character.object.y += 5
    }
    if (e.key == "q") {
        app.stage.angle += 5
        character.object.angle -= 5
    }
    if (e.key == "e") {
        app.stage.angle -= 5
        character.object.angle += 5
    }
})