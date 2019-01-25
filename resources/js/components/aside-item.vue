<template lang="pug">
  div
    li(:class="{ active: active }" @click="folded")
      slot
      slot(name="item")
</template>

<script>
export default {
  name: 'AsideItemComponent',
  data () {
    return {
      active: false
    }
  },
  methods: {
    folded () {
      this.active = !this.active
    }
  }
}
</script>

<style lang="scss" scoped>
  li > ul {
    display: none;
    animation-duration: .1s;
    animation-fill-mode: both;
    animation-name: slideInUp;
  }

  li.active {
    position: relative;
  }

  li.active > ul {
    display: block;
    animation-duration: .1s;
    animation-fill-mode: both;
    animation-name: slideInDown;
  }

  @keyframes slideInDown {
    from {
      transform: translate3d(-20%, -20%, 0);
      visibility: visible;
    }

    to {
      transform: translate3d(0, 0, 0);
    }
  }

  @keyframes slideInUp {
    from {
      transform: translate3d(0, 20%, 0);
      visibility: visible;
    }

    to {
      transform: translate3d(0, 0, 0);
    }
  }
</style>
