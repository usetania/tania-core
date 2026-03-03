/**
 * Quantity Slider Hook for Phoenix LiveView
 *
 * Usage in HEEx:
 *   <div id="move-slider" phx-hook="QuantitySlider"
 *        data-min="1" data-max="100" data-value="50"
 *        data-event="slider_changed">
 *   </div>
 */
const QuantitySlider = {
  mounted() {
    const min = parseInt(this.el.dataset.min) || 0
    const max = parseInt(this.el.dataset.max) || 100
    const value = parseInt(this.el.dataset.value) || min
    const eventName = this.el.dataset.event || "slider_changed"

    // Build slider UI using safe DOM methods
    const wrapper = document.createElement("div")
    wrapper.className = "flex items-center gap-4"

    const input = document.createElement("input")
    input.type = "range"
    input.min = min
    input.max = max
    input.value = value
    input.className = "flex-1 h-2 rounded-lg appearance-none cursor-pointer accent-[#7136A1]"
    this.updateGradient(input, value, min, max)

    const display = document.createElement("span")
    display.className = "text-lg font-semibold text-gray-700 min-w-[3rem] text-right"
    display.textContent = value

    wrapper.appendChild(input)
    wrapper.appendChild(display)
    this.el.appendChild(wrapper)

    input.addEventListener("input", (e) => {
      const val = parseInt(e.target.value)
      display.textContent = val
      this.updateGradient(input, val, min, max)
      this.pushEvent(eventName, { value: val })
    })

    // Handle server-driven updates
    this.handleEvent("update_slider", ({ min: newMin, max: newMax, value: newValue }) => {
      if (newMin !== undefined) input.min = newMin
      if (newMax !== undefined) input.max = newMax
      if (newValue !== undefined) {
        input.value = newValue
        display.textContent = newValue
        const m = parseInt(input.min)
        const M = parseInt(input.max)
        this.updateGradient(input, newValue, m, M)
      }
    })
  },

  updateGradient(input, value, min, max) {
    const percent = ((value - min) / (max - min)) * 100
    input.style.background = `linear-gradient(to right, #7136A1 0%, #7136A1 ${percent}%, #dee2e6 ${percent}%, #dee2e6 100%)`
  }
}

export default QuantitySlider
