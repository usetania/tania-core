/**
 * Leaflet Map Hook for Phoenix LiveView
 *
 * Usage in HEEx:
 *   <div id="farm-map" phx-hook="LeafletMap"
 *        data-lat="0" data-lng="0" data-zoom="13"
 *        phx-update="ignore"
 *        style="height: 400px; width: 100%;">
 *   </div>
 */
const LeafletMap = {
  mounted() {
    // Dynamically load Leaflet CSS
    if (!document.querySelector('link[href*="leaflet"]')) {
      const link = document.createElement("link")
      link.rel = "stylesheet"
      link.href = "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
      link.integrity = "sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
      link.crossOrigin = ""
      document.head.appendChild(link)
    }

    // Dynamically load Leaflet JS
    this.loadLeaflet().then(() => {
      this.initMap()
    })
  },

  loadLeaflet() {
    return new Promise((resolve) => {
      if (window.L) {
        resolve()
        return
      }

      const script = document.createElement("script")
      script.src = "https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"
      script.integrity = "sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo="
      script.crossOrigin = ""
      script.onload = () => resolve()
      document.head.appendChild(script)
    })
  },

  initMap() {
    const lat = parseFloat(this.el.dataset.lat) || 0
    const lng = parseFloat(this.el.dataset.lng) || 0
    const zoom = parseInt(this.el.dataset.zoom) || 13

    this.map = L.map(this.el).setView([lat, lng], zoom)

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      maxZoom: 19
    }).addTo(this.map)

    // Add marker if coordinates are set
    if (lat !== 0 || lng !== 0) {
      this.marker = L.marker([lat, lng]).addTo(this.map)
    }

    // Click handler to select location
    this.map.on("click", (e) => {
      const { lat, lng } = e.latlng

      if (this.marker) {
        this.marker.setLatLng([lat, lng])
      } else {
        this.marker = L.marker([lat, lng]).addTo(this.map)
      }

      this.pushEvent("map_clicked", { lat: lat, lng: lng })
    })

    // Handle updates from server
    this.handleEvent("update_map", ({ lat, lng, zoom }) => {
      if (lat && lng) {
        this.map.setView([lat, lng], zoom || this.map.getZoom())
        if (this.marker) {
          this.marker.setLatLng([lat, lng])
        } else {
          this.marker = L.marker([lat, lng]).addTo(this.map)
        }
      }
    })

    // Fix map display when container becomes visible
    setTimeout(() => this.map.invalidateSize(), 100)
  },

  destroyed() {
    if (this.map) {
      this.map.remove()
    }
  }
}

export default LeafletMap
