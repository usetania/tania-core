<template lang="pug">
  .row
    .col-xs-12.text-truncate
      label.btn.btn-default.btn-file Browse
        input(type="file" @change="processFile($event)" style="display: none;")
      span.text-muted {{ filename }}
</template>

<script>
export default {
  name: "UploadComponent",
  data () {
    return {
      filename: 'No file chosen',
    }
  },
  methods: {
    processFile (event) {
      this.filename = event.target.files[0].name
      this.$emit('fileSelelected', event.target.files[0])
    }
  }
}
</script>
