<template lang="pug">
  .crops-create.col
    .wrapper-md
      .row
        .col-sm-12.col-xs-12
          h1.m-n.font-thin.h3.text-black Crops
    .wrapper.md
      .row
        .col-sm-6.col-sm-offset-3
          .panel.panel-default(style="min-height: 100px;")
            .panel-heading
              span.h4.font-bold Add a New Batch
            .panel-body
              small.text-muted Crop Batch is lorem ipsum
              form(@submit.prevent="validateBeforeSubmit")
                .line.line-dashed.b-b.line-lg
                .row
                  .form-group
                    label.col-xs-12.control-label Select activity type of this crop batch
                    .row
                      .col-sm-6(v-for="type in options.types")
                        .radio.m-l
                          label.i-checks
                            input(type="radio" name="activityType" v-bind:value="type.key")
                            i
                            | {{ type.label }}
                .row
                  .col-xs-12
                    .form-group
                      label Area
                      select.form-control(name="arealist")
                        option - select area to grow -
                        option Florania (Nursery / Seeding)
                        option Leafy Veggie (Growing Area)
                .row
                  .col-xs-6
                    .form-group
                      label Plant Type
                      select.form-control(name="planttype")
                        option - select plant type -
                        option Vegetable
                        option Fruit
                        option Herb
                        option Flower
                        option Tree &amp; Shrub
                        option Other
                  .col-xs-6
                    .form-group
                      label.control-label Crop Variety
                      select.form-control(name="cropvariety")
                        option - select crop variety -
                        option(value="BSG") Broccoli Super Green
                        option(value="KNT") Kale Nero Toscana
                        option(value="LAR") Lettuce Abby Red
                        option(value="BLH") Bayam Lu Hsien
                        option(value="TAC") Tomato Ailsa Craig
                        option(value="ZP") Zucchini Philip
                .row
                  .col-xs-6
                    .form-group
                      label.control-label Container Quantity
                      input.form-control(type="text" name="containerquantity")
                  .col-xs-6
                    .form-group
                      label.control-label Container Type
                      select.form-control(name="plantunit")
                        option - select unit -
                        option(value="trays") Trays
                        option(value="pots") Pots
                    input.form-control(type="text" placeholder="How many cells your tray has?" name="cells")
                .row
                  .form-group
                    .col-md-12.m-t
                      button.btn.btn-default CANCEL
                      button.btn.btn-success.pull-right SAVE
</template>

<script>
import { AreaTypes } from '@/stores/helpers/farms/area'
import { mapActions } from 'vuex'
export default {
  name: "FarmCropCreate",
  data () {
    return {
      options: {
        types: Array.from(AreaTypes)
      }
    }
  },
  methods: {
  }
}
</script>
