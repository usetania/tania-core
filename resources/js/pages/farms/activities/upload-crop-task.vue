<template lang="pug">
  .upload-crop-task
    .modal-header
      span.h4.font-bold Upload Photo
      span.pull-right.text-muted(style="cursor: pointer;" @click="$parent.$emit('close')")
        i.fa.fa-close
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label Choose photo
          UploadComponent(@fileSelelected="fileSelelected")
        .form-group
          small.text-muted.pull-right (max. 200 char)
          label(for="description") Describe a bit about this photo
          textarea.form-control#description(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('description') }" v-model="task.description" name="description" rows="3")
          span.help-block.text-danger(v-show="errors.has('description')") {{ errors.first('description') }}
        .form-group
          .text-center.m-t
            button.btn.btn-primary(type="submit")
              i.fa.fa-check
              |  OK
</template>


<script>
import { StubTask } from '@/stores/stubs'
import UploadComponent from '@/components/upload'
export default {
  name: "UploadCropTask",
  components: {
    UploadComponent
  },
  data () {
    return {
      task: Object.assign({}, StubTask),
      filename: '',
    }
  },
  methods: {
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
    },
    fileSelelected (file) {
      this.task.photo = file
    }
  }
}
</script>
