<template lang="pug">
  .crops-create
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
                  input#crop_type(type="radio" name="crop.crop_type" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop.initial_area') }" v-model="crop.crop_type" v-bind:value="type.key")
                  i
                  | {{ type.label }}
        .form-group
          label Area
          select.form-control#initial_area(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop.initial_area') }" v-model="crop.initial_area" name="crop.initial_area" value="")
            option(value="") - select area to grow -
            option(v-for="area in areas" v-bind:value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('crop.initial_area')") {{ errors.first('crop.initial_area') }}
        .row
          .col-xs-6
            .form-group
              label Plant Type
              select.form-control#plant_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop.plant_type') }" v-model="crop.plant_type" name="crop.plant_type" value="" v-on:change="onChange")
                option(value="") - select plant type -
                option(v-for="type in inventories" v-bind:value="type.plant_type.code") {{ type.plant_type.code }}
              span.help-block.text-danger(v-show="errors.has('crop.plant_type')") {{ errors.first('crop.plant_type') }}
          .col-xs-6
            .form-group
              label.control-label Crop Variety
              select.form-control#variety(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop.variety') }" v-model="crop.variety" name="crop.variety")
                option(value="") - select crop variety -
                option(v-for="variety in cropVarieties" v-bind:value="variety") {{ variety }}
              span.help-block.text-danger(v-show="errors.has('crop.variety')") {{ errors.first('crop.variety') }}
        .row

          .col-xs-6
            .form-group
              label.control-label Container Quantity
              input.form-control#container_quantity(type="text" v-validate="'required|alpha_num|min:0'" :class="{'input': true, 'text-danger': errors.has('crop.container_quantity') }" v-model="crop.container_quantity" name="crop.container_quantity" value="")
              span.help-block.text-danger(v-show="errors.has('crop.container_quantity')") {{ errors.first('crop.container_quantity') }}
          .col-xs-6
            .form-group
              label.control-label Container Type
              select.form-control#container_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('crop.container_type') }" v-model="crop.container_type" name="crop.container_type" @change="typeChanged($event.target.value)" value="")
                option(value="") - select unit -
                option(v-for="container in options.cropContainers" v-bind:value="container.key") {{ container.label }}s
              span.help-block.text-danger(v-show="errors.has('crop.container_type')") {{ errors.first('crop.container_type') }}
        .row(v-if="crop.container_type == 'tray'")
          .col-xs-6.pull-right
            .form-group
              input.form-control#container_cell(type="text" placeholder="How many cells your tray has?" v-validate="'required|alpha_num|min:0'" :class="{'input': true, 'text-danger': errors.has('crop.container_cell') }" v-model="crop.container_cell" name="crop.container_cell")
              span.help-block.text-danger(v-show="errors.has('crop.container_cell')") {{ errors.first('crop.container_cell') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") Save
            i.fa.fa-long-arrow-right
          btn.btn-addon.btn-default(style="cursor: pointer;") Cancel
</template>

<script>
import { AreaTypes } from '@/stores/helpers/farms/area'
import { CropContainers } from '@/stores/helpers/farms/crop'
import { StubCrop } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal'
export default {
  name: "FarmCropCreate",
  computed: {
    ...mapGetters({
      areas: 'getAllAreas',
      inventories: 'getAllFarmInventories'
    }),
    cropVarieties: {
      get() {
        let cropVarieties = []
        for (var inventory in this.inventories) {
          if (this.inventories[inventory].plant_type.code === this.crop.plant_type) {
            cropVarieties = this.inventories[inventory].varieties
            this.crop.variety = ""
            break
          }
        }
        return cropVarieties
      },
      set(value) {
      }
    }
  },
  components : {
    Modal
  },
  data () {
    return {
      crop: Object.assign({}, StubCrop),
      options: {
        areaTypes: Array.from(AreaTypes),
        cropContainers: Array.from(CropContainers)
      }
    }
  },
  methods: {
    ...mapActions([
      'createCrop',
      'typeChanged'
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
        .then(data => this.$router.push({ name: 'FarmCrops' }))
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
