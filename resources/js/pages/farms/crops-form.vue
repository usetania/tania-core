<template lang="pug">
  .crops-form(v-if="loading === false")
    .modal-header
      span.h4.font-bold(v-if="crop.uid") Update Crop
      span.h4.font-bold(v-else) Add a New Batch
    .modal-body
      p.text-muted
        | Crop Batch is a quantity or consignment of crops done at one time.
      form(@submit.prevent="validateBeforeSubmit")
        .line.line-dashed.b-b.line-lg
        .form-group
          label.control-label Select activity type of this crop batch
          .row
            .col-sm-6(v-for="type in options.areaTypes")
              .radio
                label.i-checks
                  input#crop_type(type="radio" name="crop type" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop type') }" v-model="crop.crop_type" v-bind:value="type.key")
                  i
                  | {{ type.label }}
          span.help-block.text-danger(v-show="errors.has('crop type')") {{ errors.first('crop type') }}
        .form-group
          label Area
          select.form-control#initial_area(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('initial area') }" v-model="crop.initial_area" name="initial area")
            option(value="") - select area to grow -
            option(v-for="area in areas" v-bind:value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('initial area')") {{ errors.first('initial area') }}
        .row
          .col-xs-6
            .form-group
              label Plant Type
              select.form-control#plant_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('plant type') }" v-model="crop.plant_type" name="plant type" v-on:change="onChange")
                option(value="") - select plant type -
                option(v-for="type in inventories" v-bind:value="type.plant_type") {{ type.plant_type }}
              span.help-block.text-danger(v-show="errors.has('plant type')") {{ errors.first('plant type') }}
          .col-xs-6
            .form-group
              label.control-label Crop Variety
              select.form-control#crop-variety(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop variety') }" v-model="crop.name" name="crop variety")
                option(value="") - select crop variety -
                option(v-for="variety in cropVarieties" v-bind:value="variety") {{ variety }}
              span.help-block.text-danger(v-show="errors.has('crop variety')") {{ errors.first('crop variety') }}
        .row

          .col-xs-6
            .form-group
              label.control-label Container Quantity
              input.form-control#container_quantity(type="text" v-validate="'required|decimal|min:0'" :class="{'input': true, 'text-danger': errors.has('container quantity') }" v-model="crop.container_quantity" name="container quantity")
              span.help-block.text-danger(v-show="errors.has('container quantity')") {{ errors.first('container quantity') }}
          .col-xs-6
            .form-group
              label.control-label Container Type
              select.form-control#container_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('container type') }" v-model="crop.container_type" name="container type" @change="typeChanged($event.target.value)")
                option(value="") - select unit -
                option(v-for="container in options.containers" v-bind:value="container.key") {{ container.label }}s
              span.help-block.text-danger(v-show="errors.has('container type')") {{ errors.first('container type') }}
        .row(v-if="crop.container_type == 'TRAY'")
          .col-xs-6.pull-right
            .form-group
              input.form-control#container_cell(type="text" placeholder="How many cells your tray has?" v-validate="'required|decimal|min:0'" :class="{'input': true, 'text-danger': errors.has('container cell') }" v-model="crop.container_cell" name="container cell")
              span.help-block.text-danger(v-show="errors.has('container cell')") {{ errors.first('container cell') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") SAVE
          button.btn.btn-default(type="button" style="cursor: pointer;" @click="$parent.$emit('close')") CANCEL
</template>

<script>
import { AreaTypes } from '@/stores/helpers/farms/area'
import { Containers } from '@/stores/helpers/farms/crop'
import { StubCrop } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: "FarmCropForm",
  computed: {
    ...mapGetters({
      areas: 'getAllAreas',
    }),
    cropVarieties: {
      get() {
        let cropVarieties = []
        for (var inventory in this.inventories) {
          if (this.inventories[inventory].plant_type === this.crop.plant_type) {
            cropVarieties = this.inventories[inventory].names
            if (!cropVarieties.includes(this.crop.name)) {
              this.crop.name = ''
            }
            break
          }
        }
        return cropVarieties
      },
      set(value) {
      }
    }
  },
  data () {
    return {
      crop: Object.assign({}, StubCrop),
      inventories: [],
      loading: true,
      options: {
        areaTypes: Array.from(AreaTypes),
        containers: Array.from(Containers)
      }
    }
  },
  created () {
    this.fetchFarmInventories()
      .then(({ data }) =>  {
        this.loading = false
        this.inventories = data
      })
      .catch(error => console.log(error))
  },
  methods: {
    ...mapActions([
      'submitCrop',
      'typeChanged',
      'fetchFarmInventories'
    ]),
    onChange: function () {
      this.cropVarieties = this.cropVarieties
    },
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.submit()
        }
      })
    },
    submit () {
      this.submitCrop(this.crop)
        .then(() => this.$parent.$emit('close'))
        .catch(() => this.$toasted.error('Error in crop submission'))
    },
    typeChanged (type) {
      if (type === 'TRAY') {
        this.crop.container_cell = ''
      } else {
        this.crop.container_cell = 0
      }
    }
  },
  mounted () {
    if (typeof this.data.uid != "undefined") {
      this.crop.uid = this.data.uid
      this.crop.crop_type = this.data.type
      this.crop.initial_area = this.data.initial_area.area_id
      this.crop.plant_type = this.data.inventory.plant_type
      this.crop.name = this.data.inventory.name
      this.crop.container_quantity = this.data.container.quantity
      this.crop.container_type = this.data.container.type
      this.crop.container_cell = this.data.container.cell
    }
  },
  props: ['data'],
}
</script>
