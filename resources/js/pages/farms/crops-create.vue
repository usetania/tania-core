<template lang="pug">
  .crops-create(v-if="loading === false")
    .modal-header
      span.h4.font-bold Add a New Batch
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
                  input#crop_type(type="radio" name="crop_type" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop_type') }" v-model="crop.crop_type" v-bind:value="type.key")
                  i
                  | {{ type.label }}
        .form-group
          label Area
          select.form-control#initial_area(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('initial_area') }" v-model="crop.initial_area" name="initial_area")
            option(value="") - select area to grow -
            option(v-for="area in areas" v-bind:value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('initial_area')") {{ errors.first('initial_area') }}
        .row
          .col-xs-6
            .form-group
              label Plant Type
              select.form-control#plant_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('plant_type') }" v-model="crop.plant_type" name="plant_type" v-on:change="onChange")
                option(value="") - select plant type -
                option(v-for="type in inventories" v-bind:value="type.plant_type") {{ type.plant_type }}
              span.help-block.text-danger(v-show="errors.has('plant_type')") {{ errors.first('plant_type') }}
          .col-xs-6
            .form-group
              label.control-label Crop Variety
              select.form-control#variety(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('variety') }" v-model="crop.variety" name="variety")
                option(value="") - select crop variety -
                option(v-for="variety in cropVarieties" v-bind:value="variety") {{ variety }}
              span.help-block.text-danger(v-show="errors.has('variety')") {{ errors.first('variety') }}
        .row

          .col-xs-6
            .form-group
              label.control-label Container Quantity
              input.form-control#container_quantity(type="text" v-validate="'required|numeric|min:0'" :class="{'input': true, 'text-danger': errors.has('container_quantity') }" v-model="crop.container_quantity" name="container_quantity")
              span.help-block.text-danger(v-show="errors.has('container_quantity')") {{ errors.first('container_quantity') }}
          .col-xs-6
            .form-group
              label.control-label Container Type
              select.form-control#container_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('container_type') }" v-model="crop.container_type" name="container_type" @change="typeChanged($event.target.value)")
                option(value="") - select unit -
                option(v-for="container in options.cropContainers" v-bind:value="container.key") {{ container.label }}s
              span.help-block.text-danger(v-show="errors.has('container_type')") {{ errors.first('container_type') }}
        .row(v-if="crop.container_type == 'tray'")
          .col-xs-6.pull-right
            .form-group
              input.form-control#container_cell(type="text" placeholder="How many cells your tray has?" v-validate="'required|numeric|min:0'" :class="{'input': true, 'text-danger': errors.has('container_cell') }" v-model="crop.container_cell" name="container_cell")
              span.help-block.text-danger(v-show="errors.has('container_cell')") {{ errors.first('container_cell') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit")
            i.fa.fa-long-arrow-right
            | Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>

<script>
import { AreaTypes } from '@/stores/helpers/farms/area'
import { CropContainers } from '@/stores/helpers/farms/crop'
import { StubCrop } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: "FarmCropCreate",
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
            this.crop.name = ""
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
        cropContainers: Array.from(CropContainers)
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
      'createCrop',
      'typeChanged',
      'fetchFarmInventories'
    ]),
    onChange: function () {
      this.cropVarieties = this.cropVarieties
    },
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
      this.createCrop(this.crop)
        .then(this.$parent.$emit('close'))
        .catch(({ data }) => this.message = data)
    },
    typeChanged (type) {
      if (type === 'tray') {
        this.crop.container_cell = ''
      } else {
        this.crop.container_cell = 0
      }
    }
  }
}
</script>
