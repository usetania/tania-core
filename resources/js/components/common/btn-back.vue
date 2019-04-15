<template lang="pug">
router-link(
  :to="route"
  :class="`btn btn-light btn-md ${customClass}`"
)
  i.fa.fa-long-arrow-alt-left
  translate Back
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
    route: {
      type: Object,
      required: true,
    }
  },
};
</script>

<style lang="scss" scoped>
i {
  text-align: center;
  margin-right: 10px;
}
</style>
