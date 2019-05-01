<template lang="pug">
b-button(
  variant="primary",
  size="md"
  :class="customClass"
)
  i.fa.fa-plus
  | {{ title }}
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
  text-align: left;
  margin-right: 10px;
  width: 20px;
}
</style>
