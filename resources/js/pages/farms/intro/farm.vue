<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../../images/logobig.png")
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
                label(for="name") Farm Name
                input.form-control#name(type="text" v-validate="'required|alpha_num|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('farm.name') }" placeholder="" v-model="farm.name" name="farm.name")
                span.help-block.text-danger(v-show="errors.has('farm.name')") {{ errors.first('farm.name') }}
              .form-group
                label(for="description") Farm Description
                textarea.form-control#description(placeholder="" rows="3" v-model="farm.description")
              .form-group
                label(for="type") Farm Type
                select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('farm.farm_type') }" v-model="farm.farm_type" name="farm.farm_type")
                  option(value="") Please select type
                  option(v-for="type in types" :value="type.code") {{ type.name }}
                span.help-block.text-danger(v-show="errors.has('farm.farm_type')") {{ errors.first('farm.farm_type') }}
              .row
                .col-xs-12
                  .form-group
                    label Location
                    button.btn.btn-default.pull-right(@click="findMe" type="button")
                      i.fa.fa-crosshairs
              .form-group
                .row
                  .col-xs-6
                    input.form-control#latitude(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('farm.latitude') }" placeholder="Latitude" v-model="farm.latitude" name="farm.latitude")
                    span.help-block.text-danger(v-show="errors.has('farm.latitude')") {{ errors.first('farm.latitude') }}
                  .col-xs-6
                    input.form-control#longitude(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('farm.latitude') }" placeholder="Longitude" v-model="farm.longitude" name="farm.longitude")
                    span.help-block.text-danger(v-show="errors.has('farm.latitude')") {{ errors.first('farm.latitude') }}
              .form-group
                .row
                  .col-xs-6
                    label(for="country") Country
                    select.form-control#country(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('farm.country_code') }" v-model="farm.country_code" name="farm.country_code" @change="countryChanged($event.target.value)")
                      option(value="") Please select country
                      option(v-for="country in countries" :value="country.ID") {{ country.Name }}
                    span.help-block.text-danger(v-show="errors.has('farm.country_code')") {{ errors.first('farm.country_code') }}
                  .col-xs-6
                    label City
                    select.form-control#city(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('farm.city_code') }" v-model="farm.city_code" name="farm.city_code")
                      option(value="") Please select city
                      option(v-for="city in cities" :value="city.ID") {{ city.Name }}
                    span.help-block.text-danger(v-show="errors.has('farm.city_code')") {{ errors.first('farm.city_code') }}
              .form-group.text-center.m-t
                button.btn.btn-addon.btn-primary(type="submit")
                  i.fa.fa-long-arrow-right
                  | Continue

</template>

<script>
import { mapGetters } from 'vuex'
import mixins from '../mixins'
export default {
  name: 'FarmIntro',
  mixins: [mixins],

  computed: {
    ...mapGetters({
      currentFarm: 'getCurrentFarm',
    })
  },

  mounted () {
    // if user comming from reservoir intro aka back button, we need to check if on the state
    // have current farm we need to passed into farm data
    if (this.currentFarm && this.currentFarm.id !== '') {
      this.farm = Object.assign({}, this.currentFarm)
    }
  },

  methods: {
    redirectTo (data) {
      this.$router.push({ name: 'IntroReservoirCreate' })
    },
  }
}
</script>

