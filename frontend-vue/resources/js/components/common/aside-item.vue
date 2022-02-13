<template lang="pug">
b-nav-item(:to="this.routeName" :class="isActive ? 'active' : ''")
  i(:class="this.fontawesome")
  | {{ this.title }}
</template>

<script>
export default {
  name: 'AsideItem',
  props: {
    routeName: {
      type: Object,
      required: true,
    },
    fontawesome: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      default: '',
    },
    isActive: {
      type: Boolean,
      default: false,
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

.nav-item {
  .nav-link {
    color: #fff;
  }
}

.nav-item.active {
  font-weight: 700;
}

@media (max-width: 767.98px) {
  .nav-item {
    .nav-link {
      color: #513969;
    }
  }
}
</style>
