<template lang="pug">
  .dump-crop-task
    .modal-header
      span.h4.font-bold Dump 
        span.identifier {{ crop.batch_id }}
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="type")
            | Choose area
          select.form-control#source_area_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="task.source_area_id" name="source_area_id" @change="areaChange($event.target.value)")
            option(value="") Please select area
            option(v-for="area in current_areas" :value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('source_area_id')") {{ errors.first('source_area_id') }}
        .form-group
          label(for="quantity")
            | How many plants you want to dump?
          vue-slider(v-model="task.quantity" v-bind:min="1" v-bind:max="max_value")
          span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
        .form-group
          label(for="notes") Notes
          textarea.form-control#notes(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('notes') }" placeholder="Leave optional notes of the harvest" v-model="task.notes" name="notes" rows="2")
          span.help-block.text-danger(v-show="errors.has('notes')") {{ errors.first('notes') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit")
            i.fa.fa-long-arrow-right
            | Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')")
            i.fa.fa-close
            | Cancel
</template>


<script>
import { mapGetters, mapActions } from 'vuex'
import { StubTask } from '@/stores/stubs'
import vueSlider from 'vue-slider-component';
export default {
  name: "DumpCropTask",
  components: {
    vueSlider
  },
  computed : {
    ...mapGetters({
      areas: 'getAllAreas',
    })
  },
  data () {
    return {
      max_value: 100,
      current_areas: [],
      task: Object.assign({}, StubTask),
    }
  },
  props: ['crop'],
  mounted () {
    this.fetchAreas()
  },
  created () {
    this.task.quantity = 1
    this.current_areas = this.areas
  },
  methods: {
    ...mapActions([
      'fetchAreas',
      'dumpCrop',
      'areaChange',
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
      this.dumpCrop(this.task)
        .then(this.$parent.$emit('close'))
        .catch(({ data }) => this.message = data)
    },
    areaChange (area_id) {
      // TODO: change the max value here
    },
  }
}
</script>