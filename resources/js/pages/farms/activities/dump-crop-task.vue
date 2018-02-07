<template lang="pug">
  .dump-crop-task
    .modal-header
      span.h4.font-bold Dump Crops
      span.pull-right.text-muted(style="cursor: pointer;" @click="$parent.$emit('close')")
        i.fa.fa-close
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="quantity")
            | How many plants you want to dump?
          input.form-control#quantity(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="task.quantity" name="quantity")
          span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
        .form-group
          label(for="type")
            | Choose area
          select.form-control#source_area_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="task.source_area_id" name="source_area_id")
            option(value="") Please select area
            option(v-for="area in areas" :value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('source_area_id')") {{ errors.first('source_area_id') }}
        .form-group
          label(for="notes") Notes
          textarea.form-control#notes(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('notes') }" placeholder="Leave optional notes of the harvest" v-model="task.notes" name="notes" rows="2")
          span.help-block.text-danger(v-show="errors.has('notes')") {{ errors.first('notes') }}
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
  name: "DumpCropTask",
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
      'dumpCrop',
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
  }
}
</script>