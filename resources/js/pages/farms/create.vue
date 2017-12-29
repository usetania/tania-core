<template lang="pug">
  form(@submit.prevent="create")
    .col(v-if="intro === false")
      .wrapper-md
        .m-n.font-thin.h3.text-black Home
      .wrapper-md
        .row
          .col-sm-8.col-sm-offset-2.col-md-6.col-md-offset-3
            include _form.pug
    .intro(v-else)
      h3.text-lt.text-center.wrapper.m-t Hello! Can you tell me a little <br/> about your farm?
      .m-b-lg
        .wrapper
          include _form.pug

</template>

<script>
import stub from '@/stores/stubs/farm'
import { clone } from 'lodash'
import { mapActions, mapGetters } from 'vuex'

export default {
  name: 'FarmCreate',
  data () {
    return {
      farm: clone(stub)
    }
  },
  computed: {
    ...mapGetters({
      countries: 'getCountries',
      cities: 'getCities',
      types: 'getAllFarmTypes'
    })
  },
  props: {
    intro: { type: Boolean, default: false}
  },
  methods: {
    ...mapActions([
      'createFarm',
      'fetchCitiesByCountryCode'
    ]),

    create () {
      this.createFarm(this.farm)
        .then(response => {

        })
        .catch(error => console.log(error))
    },

    findMe () {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(position => {
          this.farm.latitude = position.coords.latitude
          this.farm.longitude = position.coords.longitude
        }, error => console.log(error))
      }
    },

    countryChanged (countryCode) {
      this.fetchCitiesByCountryCode(countryCode)
    }
  }
}
</script>
