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
          input.form-control#quantity(type="text" v-validate="'required|numeric|min:0'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="task.quantity" name="quantity")
          span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
        .form-group
          button.btn.btn-primary.pull-right(type="submit")
            i.fa.fa-check
            |  OK
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')")
            i.fa.fa-close
            | Cancel
</template>


<script>
import { mapGetters, mapActions } from 'vuex'
import { StubTask } from '@/stores/stubs'
export default {
  name: "MoveCropTask",
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
  methods: {
    ...mapActions([
      'fetchAreas',
      'moveCrop',
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
  }
}
</script>