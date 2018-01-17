<template lang="pug">
  .map
    v-map(:zoom="13" ref="map" :center="location")
      v-tile-layer(url="http://{s}.tile.osm.org/{z}/{x}/{y}.png" v-on:click="onMapClick")
      v-marker(:lat-lng="location")
</template>

<script>
import L from 'leaflet'
import Vue2Leaflet from 'vue2-leaflet'
// Build icon assets.
delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.imagePath = ''
L.Icon.Default.mergeOptions({
  iconRetinaUrl: require('../../images/marker-icon-2x.png'),
  iconUrl: require('../../images/marker-icon.png'),
  shadowUrl: require('../../images/marker-shadow.png')
})

export default {
  name: "MapboxComponent",
  data () {
    return {
      location: [-8.4960936, 115.2485298]
    }
  },
  components: {
    'v-map': Vue2Leaflet.Map,
    'v-tile-layer': Vue2Leaflet.TileLayer,
    'v-marker': Vue2Leaflet.Marker
  },
  // props: {
  //   lat: { default:'-8.4960936'},
  //   lng: { default: '115.2485298'}
  // },
  mounted () {
    this.$refs.map.mapObject.on('click', (e) => {
      this.location = [e.latlng.lat, e.latlng.lng]
      // this.centerLatLong = this.$refs.map.mapObject.getCenter()
      // console.log('map was moved')
    })
  },
  methods: {
    onMapClick (data) {
      console.log(data)
    }
  }
}
</script>


<style lang="scss" scoped>
  .map {
    height: 500px;
    width: 100%;
  }
</style>
