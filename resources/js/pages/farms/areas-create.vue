<template lang="pug">
  .areas-create
    .modal-header
      span.h4.font-bold Add New Area
    .modal-body
      p.text-muted
        | Area is a space where you grow your plants. It could be a seeding tray, a garden bed, or a
        | pot or anything that describes the different physical locations in your facility.
      form(@submit.prevent="validateBeforeSubmit")
        .line.line-dashed.b-b.line-lg
        .row
          .col-xs-6
            .form-group
              label(for="name") Area Name
              input.form-control#name(type="text" v-validate="'required|alpha_num_space|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('name') }" v-model="area.name" name="name")
              span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
          .col-xs-6
            .from-group
              label Size
              .row
                .col-xs-6
                  input.form-control#size(type="text" v-validate="'required|numeric'" :class="{'input': true, 'text-danger': errors.has('size') }" v-model="area.size" name="size")
                  span.help-block.text-danger(v-show="errors.has('size')") {{ errors.first('size') }}
                .col-xs-6
                  select.form-control(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('size_unit') }" v-model="area.size_unit" name="size_unit")
                    option(v-for="size_unit in options.size_units" :value="size_unit.key") {{ size_unit.label }}
                  span.help-block.text-danger(v-show="errors.has('size_unit')") {{ errors.first('size_unit') }}
        .row
          .col-xs-6
            .form-group
              label(for="type") Type
              select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="area.type" name="type")
                option(v-for="type in options.types" :value="type.key") {{ type.label }}
              span.help-block.text-danger(v-show="errors.has('type')") {{ errors.first('type') }}
          .col-xs-6
            .form-group
              label(for="locations") Locations
              select.form-control#locations(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('location') }" v-model="area.location" name="location")
                option(v-for="location in options.locations" :value="location.key") {{ location.label }}
              span.help-block.text-danger(v-show="errors.has('location')") {{ errors.first('location') }}
        .row
          .col-xs-6
            .form-group
              label Select Reservoir
              select.form-control(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('reservoir') }" v-model="area.reservoir_id" name="reservoir")
                option Please select reservoir
                option(v-for="reservoir in reservoirs" :value="reservoir.uid") {{ reservoir.name }}
              span.help-block.text-danger(v-show="errors.has('reservoir')") {{ errors.first('reservoir') }}
          .col-xs-6
            .form-group
              label Select photo
                small.text-muted (if any)
              UploadComponent(@fileSelelected="fileSelelected")
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") SAVE
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") CANCEL
</template>

<script>
import { AreaTypes, AreaLocations, AreaSizeUnits } from '@/stores/helpers/farms/area'
import { StubArea, StubMessage } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
import UploadComponent from '@/components/upload'
export default {
  name: "FarmAreasCreate",
  components: {
    UploadComponent
  },
  data () {
    return {
      message: Object.assign({}, StubMessage),
      area: Object.assign({}, StubArea),
      filename : 'No file chosen',
      options: {
        types: Array.from(AreaTypes),
        locations: Array.from(AreaLocations),
        size_units: Array.from(AreaSizeUnits)
      }
    }
  },
  computed: {
    ...mapGetters({
      reservoirs: 'getAllReservoirs'
    })
  },
  mounted () {
    this.fetchReservoirs()
  },
  methods: {
    ...mapActions([
      'createArea',
      'fetchReservoirs',
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
      this.createArea(this.area)
        .then(this.$parent.$emit('close'))
        .catch(({ data }) => this.message = data)
    },
    fileSelelected (file) {
      this.area.photo = file
    }
  }
}
</script>


