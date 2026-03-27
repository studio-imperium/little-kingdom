const SIZE = 8
const app = new PIXI.Application()
app.stage.scale = SIZE

async function init() {
    await app.init({
        background: "#1f1f1f",
        width: window.innerWidth,
        height: window.innerHeight,
        useContextAlpha: false,
        antialias: true,
        autoDensity: true,
        resolution: 1,
    })
    document.body.appendChild(app.canvas)

    app.stage.scale = 8
    app.canvas.style.imageRendering = "pixelated";
    app.canvas.style.imageRendering = "crisp-edges";
    
    await load_textures()
    await init_tiles()
    initiate()

    app.ticker.add(() => {
        if (!character || !character.object) {
            return
        }
        
        const centerX = app.screen.width / 2
        const centerY = app.screen.height / 2
        app.stage.pivot.set(character.object.x, character.object.y)
        app.stage.position.set(centerX, centerY)
    })
}
