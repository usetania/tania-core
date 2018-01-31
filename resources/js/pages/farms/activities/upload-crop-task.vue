<template lang="pug">
  .upload-crop-task
    .modal-header
      span.h4.font-bold Upload Photo
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label Choose photo
          .row
            .col-xs-12.text-truncate
              label.btn.btn-default.btn-file Browse
                input(type="file" @change="processFile($event)" style="display: none;")
              span.text-muted {{ filename }}
        .form-group
          small.text-muted.pull-right (max. 200 char)
          label(for="description") Describe a bit about this photo
          textarea.form-control#description(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('description') }" v-model="task.description" name="description" rows="3")
          span.help-block.text-danger(v-show="errors.has('description')") {{ errors.first('description') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>


<script>
import { mapGetters, mapActions } from 'vuex'
import { StubTask } from '@/stores/stubs'
export default {
  name: "UploadCropTask",
  data () {
    return {
      task: Object.assign({}, StubTask),
      filename: '',
    }
  },
  methods: {
    ...mapActions([
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
    processFile (event) {
      this.task.photo = event.target.files[0]
      this.filename = event.target.files[0].name
    }
  }
}
</script>
