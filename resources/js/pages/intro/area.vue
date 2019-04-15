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
          translate Perfect! Let's create a new area.

    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        b-card(
          :title="$gettext('Add New Area')"
        )
          b-card-text
            translate
              | Area is a space where you grow your plants.
              |
              | It could be a seeding tray, a garden bed, or a pot or anything
              |
              | that describes the different physical locations in your facility.
          b-form(@submit.prevent="validateBeforeSubmit")
            .line.line-dashed.b-b.line-lg
            .form-row
              .col
                .form-group
                  label#label-name(for="name")
                    translate Area Name
                  input#name.form-control(
                    type="text"
                    v-validate="'required|alpha_num_space|min:5|max:100'"
                    :class="{'input': true, 'text-danger': errors.has('name') }"
                    v-model="area.name"
                    name="name"
                  )
                  span.help-block.text-danger(v-show="errors.has('name')")
                    | {{ errors.first('name') }}
              .col
                .from-group
                  label#label-size
                    translate Size
                  .form-row
                    .col
                      input#size.form-control(
                        type="text"
                        v-validate="'required|decimal'"
                        :class="{'input': true, 'text-danger': errors.has('size') }"
                        v-model="area.size"
                        name="size"
                      )
                      span.help-block.text-danger(v-show="errors.has('size')")
                        | {{ errors.first('size') }}
                    .col
                      select#size_unit.form-control(
                        v-validate="'required'"
                        :class="{'input': true, 'text-danger': errors.has('size_unit') }"
                        v-model="area.size_unit"
                        name="size_unit"
                      )
                        option(v-for="size_unit in options.size_units" :value="size_unit.key")
                          | {{ size_unit.label }}
                      span.help-block.text-danger(v-show="errors.has('size_unit')")
                        | {{ errors.first('size_unit') }}
            .form-row
              .col
                .form-group
                  label#label-type(for="type")
                    translate Type
                  select#type.form-control(
                    v-validate="'required'"
                    :class="{'input': true, 'text-danger': errors.has('type') }"
                    v-model="area.type"
                    name="type"
                  )
                    option(v-for="type in options.types" :value="type.key") {{ type.label }}
                  span.help-block.text-danger(v-show="errors.has('type')")
                    | {{ errors.first('type') }}
              .col
                .form-group
                  label#label-location(for="location")
                    translate Locations
                  select#location.form-control(
                    v-validate="'required'"
                    :class="{'input': true, 'text-danger': errors.has('location') }"
                    v-model="area.location"
                    name="location"
                  )
                    option(v-for="location in options.locations" :value="location.key")
                      | {{ location.label }}
                  span.help-block.text-danger(v-show="errors.has('location')")
                    | {{ errors.first('location') }}
            .form-row
              .col
                .form-group
                  label#label-reservoir(for="reservoir")
                    translate Select Reservoir
                  select#reservoir.form-control(
                    v-validate="'required'"
                    :class="{'input': true, 'text-danger': errors.has('reservoir') }"
                    v-model="area.reservoir_id"
                    name="reservoir"
                  )
                    option(value = "")
                      translate Please select reservoir
                    option(:value="reservoir.uid ? reservoir.uid : reservoir.name")
                      | {{ reservoir.name }}
                  span.help-block.text-danger(v-show="errors.has('reservoir')")
                    | {{ errors.first('reservoir') }}
              .col
                .form-group
                  label
                    translate Select photo
                    small.text-muted (if any)
                  UploadComponent(@fileSelelected="fileSelelected")
            .form-group
              BtnContinue(title="Finish Setup" customClass="float-right")
              BtnBack(:route="{name: 'IntroReservoirCreate'}")
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { AreaTypes, AreaLocations, AreaSizeUnits } from '../../stores/helpers/farms/area';
import { StubArea, StubMessage } from '../../stores/stubs';
import UploadComponent from '../../components/upload.vue';
import BtnContinue from '../../components/common/btn-continue.vue';
import BtnBack from '../../components/common/btn-back.vue';

export default {
  name: 'AreaIntro',
  components: {
    UploadComponent,
    BtnContinue,
    BtnBack,
  },
  data() {
    return {
      message: Object.assign({}, StubMessage),
      area: Object.assign({}, StubArea),
      options: {
        types: Array.from(AreaTypes),
        locations: Array.from(AreaLocations),
        size_units: Array.from(AreaSizeUnits),
      },
    };
  },
  computed: {
    ...mapGetters({
      reservoir: 'introGetReservoir',
      currentArea: 'introGetArea',
      currentFarm: 'introGetFarm',
    }),
  },
  mounted() {
    if (this.currentArea) {
      this.area = Object.assign({}, this.currentArea);
      this.area.size_unit = this.options.size_units[0].key;
      this.area.location = this.options.locations[0].key;
      this.area.type = this.options.types[0].key;
      this.area.reservoir_id = this.reservoir.uid ? this.reservoir.uid : this.reservoir.name;
    }

    if (this.reservoir.name === '') {
      this.$router.push({ name: 'IntroReservoirCreate' });
    }

    if (this.currentFarm.name === '') {
      this.$router.push({ name: 'IntroFarmCreate' });
    }
  },
  methods: {
    ...mapActions([
      'introSetArea',
      'introCreateFarm',
      'introCreateReservoir',
      'introCreateArea',
    ]),
    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.introSetArea(this.area);
          this.toCreateFarm();
        }
      });
    },
    toCreateFarm() {
      this.introCreateFarm()
        .then(() => {
          this.toCreateReservoir();
        }).catch(() => {
          this.$router.push({ name: 'IntroFarmCreate' });
        });
    },
    toCreateReservoir() {
      this.introCreateReservoir()
        .then(() => {
          this.toCreateArea();
        }).catch(() => {
          this.$router.push({ name: 'IntroReservoirCreate' });
        });
    },
    toCreateArea() {
      this.introCreateArea()
        .then(() => {
          this.$router.push({ name: 'Home' });
        }).catch(() => {
          this.$toasted.error('Error in area submission');
        });
    },
    fileSelelected(file) {
      this.area.photo = file;
    },
  },
};
</script>

<style lang="scss" scoped>
.container-intro {
  padding-top: 20px;
  padding-bottom: 20px;
}
</style>

