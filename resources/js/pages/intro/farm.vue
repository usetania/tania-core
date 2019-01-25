<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../images/logobig.png")
    h3.text-lt.text-center.wrapper.m-t Hello! Can you tell me a little <br/> about your farm?
    .m-b-lg
      .wrapper
        form(@submit.prevent="validateBeforeSubmit")
          .panel.panel-default
            .panel-heading
              span.h4.font-bold Create Farm
            .panel-body
              error-message(:message="message.error_message")
              .form-group
                label#label-name(for="name") Farm Name
                input.form-control#name(type="text" v-validate="'required|alpha_num_space|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('name') }" placeholder="" v-model="farm.name" name="name")
                span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
              .form-group
                label#label-description(for="description") Farm Description
                textarea.form-control#description(placeholder="" rows="3" v-model="farm.description")
              .form-group
                label#label-type(for="type") Farm Type
                select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="farm.farm_type" name="type")
                  option(value="") Please select type
                  option(v-for="type in types" :value="type.code") {{ type.name }}
                span.help-block.text-danger(v-show="errors.has('type')") {{ errors.first('type') }}
              .form-group
                mapbox(:latitude="farm.latitude" :longitude="farm.longitude" @change="onLocationChange")
              .form-group
                .row
                  .col-xs-6
                    input.form-control#latitude(type="text" v-validate="'required|latitude'" :class="{'input': true, 'text-danger': errors.has('latitude') }" placeholder="Latitude" v-model="farm.latitude" name="latitude")
                    span.help-block.text-danger(v-show="errors.has('latitude')") {{ errors.first('latitude') }}
                  .col-xs-6
                    input.form-control#longitude(type="text" v-validate="'required|longitude'" :class="{'input': true, 'text-danger': errors.has('latitude') }" placeholder="Longitude" v-model="farm.longitude" name="longitude")
                    span.help-block.text-danger(v-show="errors.has('latitude')") {{ errors.first('latitude') }}
              .form-group
                .row
                  .col-xs-6
                    label#label-country(for="country") Country
                    select.form-control#country(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('country') }" v-model="farm.country" name="country")
                      option(value="") Please select country
                      option(v-for="country in countries" :value="country.name") {{ country.name }}
                    span.help-block.text-danger(v-show="errors.has('country')") {{ errors.first('country') }}
                  .col-xs-6
                    label#label-city(for="city") City
                    input.form-control#city(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('city') }" v-model="farm.city" name="city")
                    span.help-block.text-danger(v-show="errors.has('city')") {{ errors.first('city') }}
              .form-group.text-center.m-t
                button.btn.btn-addon.btn-primary(type="submit")
                  i.fas.fa-long-arrow-alt-right
                  | Continue

</template>

<script>
import Mapbox from '@/components/mapbox.vue'
import { StubFarm, StubMessage } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'

export default {
  data () {
    return {
      message: Object.assign({}, StubMessage),
      farm: Object.assign({}, StubFarm),
      mapbox: {}
    }
  },
  components : {
    Mapbox
  },

  computed: {
    ...mapGetters({
      countries: 'getCountries',
      cities: 'getCities',
      types: 'getAllFarmTypes',
      currentFarm: 'introGetFarm',
    })
  },

  mounted () {
    if (this.currentFarm) {
      this.farm = Object.assign({}, this.currentFarm)
    }
  },

  methods: {
    ...mapActions([
      'introSetFarm',
      'fetchCitiesByCountryCode'
    ]),

    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },

    create () {
      this.introSetFarm(this.farm)
      this.$router.push({ name: 'IntroReservoirCreate' })
    },

    onLocationChange (location) {
      this.farm.latitude = location.latitude
      this.farm.longitude = location.longitude
    },
  }
}

</script>


<style lang="scss" scoped>
 #map {
  width: 100%;
  height: 500px;
}
</style>
