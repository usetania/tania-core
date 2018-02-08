<template lang="pug">
  .move-crop-task
    .modal-header
      span.h4.font-bold Move 
        span.identifier {{ crop.batch_id }}
      span.pull-right.text-muted(style="cursor: pointer;" @click="$parent.$emit('close')")
        i.fa.fa-close
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="dest_area_id") Select area to move this crop
          select.form-control#destination_area_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('destination_area_id') }" v-model="task.destination_area_id" name="destination_area_id")
            option(value="") Please select area
            option(v-for="area in areas" :value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('destination_area_id')") {{ errors.first('destination_area_id') }}
        .form-group
          label(for="quantity")
            | How many plants do you want to move?
          vue-slider(v-model="task.quantity" v-bind:min="1" v-bind:max="max_value")
          span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
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
  name: "MoveCropTask",
  components: {
    vueSlider
  },
  computed : {
    ...mapGetters({
      areas: 'getAllAreas',
    })
  },
  created () {
    this.task.quantity = 1
    this.current_areas = this.areas
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
  methods: {
    ...mapActions([
      'fetchAreas',
      'moveCrop',
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
      if (this.crop.moved_area.length > 0) {
        this.task.source_area_id = this.crop.moved_area[this.crop.moved_area.length - 1].area.uid;
      } else {
        this.task.source_area_id = this.crop.initial_area.area.uid;
      }
      this.task.obj_uid = this.crop.uid
      this.moveCrop(this.task)
        .then(this.$parent.$emit('close'))
        .catch(({ data }) => this.message = data)
    },
    areaChange (area_id) {
      // TODO: change the max value here
    },
  }
}
</script>