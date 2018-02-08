<template lang="pug">
  .upload-crop-task
    .modal-header
      span.h4.font-bold Upload Photo
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label Choose photo
          UploadComponent(@fileSelelected="fileSelelected")
        .form-group
          small.text-muted.pull-right (max. 200 char)
          label(for="description") Describe a bit about this photo
          textarea.form-control#description(type="text" :class="{'input': true, 'text-danger': errors.has('description') }" v-model="task.description" name="description" rows="3")
          span.help-block.text-danger(v-show="errors.has('description')") {{ errors.first('description') }}
        .form-group
          button.btn.btn-primary.btn-addon.pull-right(type="submit")
            i.fa.fa-check
            |  OK
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')")
            i.fa.fa-close
            |  Cancel
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
