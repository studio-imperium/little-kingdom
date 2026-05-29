let current_tab = null

const inventoryNode = document.getElementById("inventory")

function close_tab() {
  if (current_tab != null) {
    current_tab.classList.add("hidden")
    current_tab = null
  }
}
function open_tab(tab) {
  current_tab = tab
  current_tab.classList.remove("hidden")
}

function toggle_inventory() {
  if (current_tab == null) {
    open_tab(inventoryNode)
    const rect = preview_canvas.getBoundingClientRect()
    preview.renderer.resize(rect.width, rect.height)
    update_preview()
  } else {
    close_tab()
  }
}

// Number keys 1-6 quick-select the matching hotbar slot. Keyed off e.code
// (physical Digit row) so Shift/other modifiers don't break it.
const hotbar_hotkeys = {
  Digit1: 0,
  Digit2: 1,
  Digit3: 2,
  Digit4: 3,
  Digit5: 4,
  Digit6: 5,
}

document.addEventListener("keydown", (e) => {
  if (chat_focused) {
    return
  }
  if (e.key == "Escape") {
    close_tab()
  }
  if (e.key == "e" || e.key == "b") {
    toggle_inventory()
  }
  if (!e.repeat && e.code in hotbar_hotkeys) {
    select_slot(hotbar_hotkeys[e.code])
  }
  // Drop the item in the currently selected hotbar slot.
  if (!e.repeat && e.key == "q") {
    drop_item(selected_slot)
  }
})
