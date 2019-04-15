<template lang="pug">
b-button(
  variant="secondary",
  size="md"
  :class="customClass"
)
  i.fa.fa-tint
  translate Watering
</template>

<script>
export default {
  name: 'BtnWatering',
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
.btn-secondary {
  color: #fff;
}

i {
  text-align: left;
  margin-right: 10px;
  width: 20px;
}
</style>
