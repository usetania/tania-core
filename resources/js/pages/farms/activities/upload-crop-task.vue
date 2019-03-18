<template lang="pug">
  .upload-crop-task
    .modal-header
      span.h4.font-bold Take Picture
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
          button.btn.btn-addon.btn-primary.pull-right(type="submit")
            i.fas.fa-check
            |  OK
          button.btn.btn-addon.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')")
            i.fas.fa-times
            |  Cancel
</template>


<script>
import { mapGetters, mapActions } from 'vuex'
import { StubTask } from '../../../stores/stubs'
import UploadComponent from '../../../components/upload.vue'
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
    ...mapActions([
      'photoCrop',
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
      this.photoCrop(this.task)
        .then(() => this.$parent.$emit('close'))
        .catch(() => this.$toasted.error('Error in crop image upload'))
    },
    fileSelelected (file) {
      this.task.photo = file
    }
  },
  props: ['crop'],
}
</script>
