<template lang="pug">
b-button(
  variant="primary",
  size="md"
  type="submit"
  :class="customClass"
)
  i.fa.fa-long-arrow-alt-right
  | {{ title }}
</template>

<script>
export default {
  name: 'BtnContinue',
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
