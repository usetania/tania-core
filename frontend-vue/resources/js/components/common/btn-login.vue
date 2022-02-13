<template lang="pug">
b-button(
  type="submit"
  variant="primary"
  size="lg"
)
  i.fas.fa-unlock
  translate Login
</template>

<script>
export default {
  name: 'BtnLogin',
};
</script>

<style lang="scss" scoped>
i {
  text-align: center;
  margin-right: 10px;
}
</style>
