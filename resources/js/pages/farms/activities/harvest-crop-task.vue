<template lang="pug">
  .harvest-crop-task
    .modal-header
      span.h4.font-bold Harvest 
        span.identifier {{ crop.batch_id }}
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .row
          .col-xs-6
            .form-group
              label(for="source_area_id") 
                | Choose area to be harvested
              select.form-control#source_area_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('source_area_id') }" v-model="task.source_area_id" name="source_area_id")
                option(value="") Please select area
                option(v-for="area in areas" :value="area.uid") {{ area.name }}
              span.help-block.text-danger(v-show="errors.has('source_area_id')") {{ errors.first('source_area_id') }}
          .col-xs-6
            .form-group
              label(for="harvest_type")
                | Choose type of harvesting
              select.form-control#harvest_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('harvest_type') }" v-model="task.harvest_type" name="harvest_type")
                option(value="ALL") All
                option(value="PARTIAL") Partial
              span.help-block.text-danger(v-show="errors.has('harvest_type')") {{ errors.first('harvest_type') }}
        .row
          .col-xs-6
            .form-group
              label(for="produced_quantity") Quantity
              input.form-control#produced_quantity(type="text" v-validate="'required|decimal'" :class="{'input': true, 'text-danger': errors.has('produced_quantity') }" v-model="task.produced_quantity" name="produced_quantity")
              span.help-block.text-danger(v-show="errors.has('produced_quantity')") {{ errors.first('produced_quantity') }}
          .col-xs-6
            .form-group
              label(for="units") Units
              select.form-control#produced_unit(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('produced_unit') }" v-model="task.produced_unit" name="produced_unit")
                option(value="Gr") Grams
                option(value="Kg") Kilograms
              span.help-block.text-danger(v-show="errors.has('produced_unit')") {{ errors.first('produced_unit') }}
        .form-group
          label(for="notes") Notes
          textarea.form-control#notes(type="text" :class="{'input': true, 'text-danger': errors.has('notes') }" placeholder="Leave optional notes of the harvest" v-model="task.notes" name="notes" rows="2")
          span.help-block.text-danger(v-show="errors.has('notes')") {{ errors.first('notes') }}
        .form-group
          button.btn.btn-addon.btn-primary.pull-right(type="submit")
            i.fas.fa-check
            |  OK
          button.btn.btn-addon.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')")
            i.fas.fa-times
            |  Cancel
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
      'harvestCrop',
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
      this.task.obj_uid = this.crop.uid
      this.harvestCrop(this.task)
        .then(this.$parent.$emit('close'))
        .catch(({ data }) => this.message = data)
    },
  }
}
</script>