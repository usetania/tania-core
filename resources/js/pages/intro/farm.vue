<template lang="pug">
  .container-fluid.container-intro
    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        .text-center
          img(
            src="../../../images/logobig.png"
            alt="Tania Logo"
            width="200"
          )

    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        h3.text-center
          translate Hello! Can you tell me a little about your farm?

    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        b-form(@submit.prevent="validateBeforeSubmit")
          b-card(
            :title="$gettext('Create Farm')"
          )
            error-message(:message="message.error_message")
            .form-group
              label#label-name(for="name")
                translate Farm Name
              input.form-control#name(
                type="text"
                v-validate="'required|alpha_num_space|min:5|max:100'"
                :class="{'input': true, 'text-danger': errors.has('name') }"
                placeholder=""
                v-model="farm.name"
                name="name"
              )
              span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
            .form-group
              label#label-description(for="description")
                translate Farm Description
              textarea.form-control#description(
                placeholder=""
                rows="3"
                v-model="farm.description"
              )
            .form-group
              label#label-type(for="type")
                translate Farm Type
              select.form-control#type(
                v-validate="'required'"
                :class="{'input': true, 'text-danger': errors.has('type') }"
                v-model="farm.farm_type"
                name="type"
              )
                option(value="")
                  translate Please select type
                option(v-for="type in types" :value="type.code") {{ type.name }}
              span.help-block.text-danger(v-show="errors.has('type')") {{ errors.first('type') }}
            .form-group
              mapbox(
                :latitude="farm.latitude"
                :longitude="farm.longitude"
                @change="onLocationChange"
              )
            .form-row.form-group
              .col
                input.form-control#latitude(
                  type="text" v-validate="'required|latitude'"
                  :class="{'input': true, 'text-danger': errors.has('latitude') }"
                  placeholder="Latitude"
                  v-model="farm.latitude"
                  name="latitude"
                )
                span.help-block.text-danger(v-show="errors.has('latitude')")
                  | {{ errors.first('latitude') }}
              .col
                input.form-control#longitude(
                  type="text"
                  v-validate="'required|longitude'"
                  :class="{'input': true, 'text-danger': errors.has('latitude') }"
                  placeholder="Longitude"
                  v-model="farm.longitude" name="longitude"
                )
                span.help-block.text-danger(v-show="errors.has('latitude')")
                  | {{ errors.first('latitude') }}
            .form-row.form-group
              .col
                label#label-country(for="country") Country
                select.form-control#country(
                  v-validate="'required'"
                  :class="{'input': true, 'text-danger': errors.has('country') }"
                  v-model="farm.country"
                  name="country"
                )
                  option(value="")
                    translate Please select country
                  option(v-for="country in countries" :value="country.name") {{ country.name }}
                span.help-block.text-danger(v-show="errors.has('country')")
                  | {{ errors.first('country') }}
              .col
                label#label-city(for="city")
                  translate City
                input.form-control#city(
                  v-validate="'required'"
                  :class="{'input': true, 'text-danger': errors.has('city') }"
                  v-model="farm.city"
                  name="city"
                )
                span.help-block.text-danger(v-show="errors.has('city')")
                  | {{ errors.first('city') }}
            .form-group
              BtnContinue(title="Continue" customClass="float-right")

</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import Mapbox from '../../components/mapbox.vue';
import { StubFarm, StubMessage } from '../../stores/stubs';
import BtnContinue from '../../components/common/btn-continue.vue';

export default {
  components: {
    Mapbox,
    BtnContinue,
  },

  data() {
    return {
      message: Object.assign({}, StubMessage),
      farm: Object.assign({}, StubFarm),
      mapbox: {},
    };
  },

  computed: {
    ...mapGetters({
      countries: 'getCountries',
      cities: 'getCities',
      types: 'getAllFarmTypes',
      currentFarm: 'introGetFarm',
    }),
  },

  mounted() {
    if (this.currentFarm) {
      this.farm = Object.assign({}, this.currentFarm);
    }
  },

  methods: {
    ...mapActions([
      'introSetFarm',
      'fetchCitiesByCountryCode',
    ]),

    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.create();
        }
      });
    },

    create() {
      this.introSetFarm(this.farm);
      this.$router.push({ name: 'IntroReservoirCreate' });
    },

    onLocationChange(location) {
      this.farm.latitude = location.latitude;
      this.farm.longitude = location.longitude;
    },
  },
};
</script>


<style lang="scss" scoped>
#map {
  width: 100%;
  height: 500px;
}

.container-intro {
  padding-top: 20px;
  padding-bottom: 20px;
}
</style>
