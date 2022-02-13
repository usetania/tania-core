<template lang="pug">
  .mapbox
    .mapbox-button
      label#label-location
        translate Location
      button.btn.btn-default.pull-right(@click="findMe" type="button")
        i.fa.fa-crosshairs
    v-map.map(:zoom="13" ref="map" :center="location")
      v-tile-layer(url="http://{s}.tile.osm.org/{z}/{x}/{y}.png" v-on:click="onMapClick")
      v-marker(:lat-lng="location")
</template>

<script>
import L from 'leaflet';
import Vue2Leaflet from 'vue2-leaflet';

const iconRetina = require('../../images/marker-icon-2x.png');
const iconMarker = require('../../images/marker-icon.png');
const iconShadow = require('../../images/marker-shadow.png');

// Build icon assets.
/* eslint no-underscore-dangle: ["error", { "allow": ["_getIconUrl"] }] */
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.imagePath = '';
L.Icon.Default.mergeOptions({
  iconRetinaUrl: iconRetina,
  iconUrl: iconMarker,
  shadowUrl: iconShadow,
});

export default {
  name: 'MapboxComponent',
  components: {
    'v-map': Vue2Leaflet.Map,
    'v-tile-layer': Vue2Leaflet.TileLayer,
    'v-marker': Vue2Leaflet.Marker,
  },
  props: {
    latitude: {
      default: -8.4960936,
    },
    longitude: {
      default: 115.2485298,
    },
  },
  data() {
    return {
      location: [-8.4960936, 115.2485298],
    };
  },
  // watcher props if the props value is not equal location state
  // whe need to change the location data from the props
  watch: {
    latitude(value) {
      if (value && value !== this.location[0] && this.isFloat(value)) {
        this.location = [parseFloat(value), parseFloat(this.location[1])];
      }
    },
    longitude(value) {
      if (value && value !== this.location[1] && this.isFloat(value)) {
        this.location = [parseFloat(this.location[0]), parseFloat(value)];
      }
    },
  },
  created() {
    this.location = Array.from([
      this.latitude !== '' ? parseFloat(this.latitude) : -8.4960936,
      this.longitude !== '' ? parseFloat(this.longitude) : 115.2485298,
    ]);
  },
  mounted() {
    this.$refs.map.mapObject.on('click', this.onMapClick);
  },

  methods: {
    onMapClick(e) {
      this.location = [e.latlng.lat, e.latlng.lng];
      this.publish();
    },

    // publish the change event, so the parent component can trigger and catch the data
    publish() {
      this.$emit('change', {
        latitude: this.location[0],
        longitude: this.location[1],
      });
    },

    findMe() {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition((position) => {
          this.location = [
            position.coords.latitude,
            position.coords.longitude,
          ];
          this.publish();
        }, error => error);
      }
    },

    isFloat(value) {
      const regexp = /^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$/;
      return value.search(regexp) === 0;
    },

  },
};
</script>


<style lang="scss" scoped>
  .mapbox-button {
    margin-bottom: 15px;
  }
  .map {
    margin-top: 10px;
    height: 400px;
    width: 100%;
  }
  #label-location {
    margin-right: 20px;
  }
</style>
