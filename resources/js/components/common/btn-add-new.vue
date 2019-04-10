<template lang="pug">
b-button(
  variant="primary",
  size="md"
)
  i.fa.fa-plus
  | {{ this.title }}
</template>

<script>
export default {
  name: 'BtnAddNew',
  props: {
    title: {
      type: String,
      required: true,
    },
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
