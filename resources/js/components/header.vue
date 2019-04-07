<template lang="pug">
  b-navbar(toggleable="lg" type="dark" variant="dark")
    b-navbar-brand(href="/")
      img.mobile-brand.d-md-none(src="../../images/logo.png" alt="Tania Logo")

    b-navbar-toggle(target="nav-collapse")

    b-collapse(id="nav-collapse" is-nav="is-nav")
      b-navbar-nav.d-none.d-md-block
        b-nav-item(href="#") {{ farm.name }}

      b-navbar-nav.d-md-none
        b-nav-item(:to="{ name: 'Home' }" :class="active ? 'active': ''" @click="clickHandler")
          i.fa.fa-home
          translate Dashboard

        b-nav-item(
          :to="{name: 'FarmReservoirs'}"
          :class="active ? 'active': ''"
          @click="clickHandler"
        )
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

      b-dropdown-divider.d-md-none

      b-navbar-nav.ml-auto
        b-nav-item(href="#" @click.prevent="signout")
          i.fa.fa-power-off
          translate Sign Out
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import { ls } from '../services';

export default {
  name: 'AppHeaderComponent',
  data() {
    return {
      dropdown: false,
      active: false,
    };
  },
  computed: {
    ...mapGetters({
      farm: 'getCurrentFarm',
      farms: 'getAllFarms',
    }),
  },
  methods: {
    ...mapActions([
      'setCurrentFarm',
      'userSignOut',
    ]),
    setFarm(farmId) {
      this.setCurrentFarm(farmId)
        .then(() => this.dropdownToggle())
        .catch(() => `Farm ${farmId} is not found in the data`);
    },
    dropdownToggle() {
      this.dropdown = !this.dropdown;
    },
    signout() {
      this.userSignOut()
        .then(() => {
          ls.remove('vuex');
          this.$router.push({ name: 'AuthLogin' });
        })
        .catch(error => error);
    },
    clickHandler() {
      this.active = !this.active;
    },
  },
};
</script>

<style lang="scss" scoped>
.mobile-brand {
  width: 100px;
}

i {
  width: 30px;
}
</style>
