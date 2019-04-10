<template lang="pug">
b-button(
  variant="secondary"
  size="md"
  :class="this.customClass"
)
  i.fa.fa-check
  translate Save
</template>

<script>
export default {
  name: 'BtnSave',
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

.btn-secondary {
  color: #fff;
}
</style>
