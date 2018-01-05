<template lang="pug">
  .alert.alert-danger(v-if="message")
    p {{ message }}
</template>

<script>
export default {
  name: 'ErrorComponent',
  props: {
    message: { type: String, default: null}
  }
}
</script>

