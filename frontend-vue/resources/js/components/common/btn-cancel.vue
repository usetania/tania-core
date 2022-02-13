<template lang="pug">
b-button(
  variant="light"
  size="md"
  :class="this.customClass"
)
  i.fa.fa-times
  translate Cancel
</template>

<script>
export default {
  name: 'BtnCancel',
  props: {
    customClass: {
      type: String,
      required: false,
      default: '',
    },
  },
};
</script>

<style lang="scss" scoped>
i {
  text-align: center;
  margin-right: 10px;
}
</style>
