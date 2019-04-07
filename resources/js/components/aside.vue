<template lang="pug">
aside#aside.bg-dark
  .brand-image.text-center
    img(src="../../images/logo.png" alt="Tania Logo")

  b-nav(vertical="vertical")
    b-nav-item(:to="{ name: 'Home' }" :class="active ? 'active': ''" @click="clickHandler")
      i.fa.fa-home
      translate Dashboard

    b-nav-item(:to="{name: 'FarmReservoirs'}" :class="active ? 'active': ''" @click="clickHandler")
      i.fa.fa-tint
      translate Reservoirs

    b-nav-item(:to="{name: 'FarmAreas'}" :class="active ? 'active': ''" @click="clickHandler")
      i.fa.fa-grip-horizontal
      translate Areas

    b-nav-item(
      :to="{ name: 'InventoriesMaterials' }"
      :class="active ? 'active': ''"
      @click="clickHandler"
    )
      i.fa.fa-archive
      translate Inventories

    b-nav-item(:to="{name: 'FarmCrops'}" :class="active ? 'active': ''" @click="clickHandler")
      i.fa.fa-leaf
      translate Crops

    b-nav-item(:to="{ name: 'Task' }" :class="active ? 'active': ''" @click="clickHandler")
      i.fa.fa-clipboard
      translate Tasks

    b-nav-item(:to="{ name: 'Account' }" :class="active ? 'active': ''" @click="clickHandler")
      i.fa.fa-user
      translate Account
</template>

<script>
export default {
  name: 'AppAsideComponent',
  props: {
    folded: { type: Boolean, default: false },
  },
  data() {
    return {
      active: false,
    };
  },
  methods: {
    setFolded() {
      this.$emit('header-folded');
    },
    clickHandler() {
      this.active = !this.active;
    },
  },
};
</script>

<style lang="scss" scoped>
aside#aside {
  height: 100%;

  .nav-link {
    color: rgba(255, 255, 255, 0.5);
  }

  i {
    width: 30px;
  }

  .brand-image {
    padding: 15px 0;

    img {
      width: 100px;
    }
  }
}
</style>
