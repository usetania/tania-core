<template lang="pug">
  .harvest-crop-task
    .modal-header
      span.h4.font-bold Move Crops
      span.pull-right.text-muted(style="cursor: pointer;" @click="$parent.$emit('close')")
        i.fa.fa-close
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="source_area_id") 
            | Choose area where you want to harvest 
            span.identifier-sm {{ crop.batch_id }}
          select.form-control#source_area_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('source_area_id') }" v-model="task.source_area_id" name="source_area_id")
            option(value="") Please select area
            option(v-for="area in areas" :value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('source_area_id')") {{ errors.first('source_area_id') }}
        .form-group
          label(for="type")
            | Choose type of harvesting
          select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="task.type" name="type")
            option(value="all") All
            option(value="partial") Partial
          span.help-block.text-danger(v-show="errors.has('type')") {{ errors.first('type') }}
        .row
          .col-xs-6
            .form-group
              label(for="quantity") Quantity
              input.form-control#quantity(type="text" v-validate="'required|alpha_num_space|min:1'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="task.quantity" name="quantity")
              span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
          .col-xs-6
            .form-group
              label(for="units") Units
              select.form-control#unit(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('unit') }" v-model="task.unit" name="unit")
                option(value="g") Grams
                option(value="kg") Kilograms
              span.help-block.text-danger(v-show="errors.has('unit')") {{ errors.first('unit') }}
        .form-group
          .text-center.m-t
            button.btn.btn-primary(type="submit")
              i.fa.fa-check
              |  OK
</template>


<script>
import { mapGetters, mapActions } from 'vuex'
import { StubTask } from '@/stores/stubs'
export default {
  name: "HarvestCropTask",
  computed : {
    ...mapGetters({
      areas: 'getAllAreas',
    })
  },
  data () {
    return {
      task: Object.assign({}, StubTask),
    }
  },
  props: ['crop'],
  mounted () {
    this.fetchAreas()
  },
  created () {
  },
  methods: {
    ...mapActions([
      'fetchAreas',
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
    },
  }
}
</script>