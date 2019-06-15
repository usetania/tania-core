<template lang="pug">
  b-navbar(toggleable="lg" type="light" variant="white")
    b-navbar-brand(href="/")
      img.mobile-brand.d-md-none(src="../../images/logobig.png" alt="Tania Logo")

    b-navbar-toggle(target="nav-collapse")

    b-collapse(id="nav-collapse" is-nav="is-nav")
      b-navbar-nav.d-none.d-md-block
        b-nav-item(href="#") {{ farm.name }}

      b-navbar-nav.d-md-none
        AsideItem(
          title="Dashboard"
          fontawesome="fa fa-home"
          :routeName="{ name: 'Home' }"
          :isActive="homeActive"
          v-on:click.native="homeClickHandler"
        )

        AsideItem(
          title="Reservoirs"
          fontawesome="fa fa-tint"
          :routeName="{ name: 'FarmReservoirs' }"
          :isActive="reservoirsActive"
          v-on:click.native="reservoirsClickHandler"
        )

        AsideItem(
          title="Areas"
          fontawesome="fa fa-grip-horizontal"
          :routeName="{ name: 'FarmAreas' }"
          :isActive="areasActive"
          v-on:click.native="areasClickHandler"
        )

        AsideItem(
          title="Materials"
          fontawesome="fa fa-archive"
          :routeName="{ name: 'InventoriesMaterials' }"
          :isActive="materialsActive"
          v-on:click.native="materialsClickHandler"
        )

        AsideItem(
          title="Crops"
          fontawesome="fa fa-leaf"
          :routeName="{ name: 'FarmCrops' }"
          :isActive="cropsActive"
          v-on:click.native="cropsClickHandler"
        )

        AsideItem(
          title="Tasks"
          fontawesome="fa fa-clipboard"
          :routeName="{ name: 'Task' }"
          :isActive="tasksActive"
          v-on:click.native="tasksClickHandler"
        )

        AsideItem(
          title="Account"
          fontawesome="fa fa-user"
          :routeName="{ name: 'Account' }"
          :isActive="accountActive"
          v-on:click.native="accountClickHandler"
        )

      b-navbar-nav.ml-auto
        b-nav-item(href="#" @click.prevent="signout")
          i.fa.fa-power-off
          translate Sign Out
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import { ls } from '../services';
import AsideItem from './common/aside-item.vue';

export default {
  name: 'AppHeaderComponent',
  components: {
    AsideItem,
  },
  data() {
    return {
      dropdown: false,
      homeActive: true,
      reservoirsActive: false,
      areasActive: false,
      materialsActive: false,
      cropsActive: false,
      tasksActive: false,
      accountActive: false,
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
    homeClickHandler() {
      this.homeActive = true;
      this.reservoirsActive = false;
      this.areasActive = false;
      this.materialsActive = false;
      this.cropsActive = false;
      this.tasksActive = false;
      this.accountActive = false;
    },
    reservoirsClickHandler() {
      this.homeActive = false;
      this.reservoirsActive = true;
      this.areasActive = false;
      this.materialsActive = false;
      this.cropsActive = false;
      this.tasksActive = false;
      this.accountActive = false;
    },
    areasClickHandler() {
      this.homeActive = false;
      this.reservoirsActive = false;
      this.areasActive = true;
      this.materialsActive = false;
      this.cropsActive = false;
      this.tasksActive = false;
      this.accountActive = false;
    },
    materialsClickHandler() {
      this.homeActive = false;
      this.reservoirsActive = false;
      this.areasActive = false;
      this.materialsActive = true;
      this.cropsActive = false;
      this.tasksActive = false;
      this.accountActive = false;
    },
    cropsClickHandler() {
      this.homeActive = false;
      this.reservoirsActive = false;
      this.areasActive = false;
      this.materialsActive = false;
      this.cropsActive = true;
      this.tasksActive = false;
      this.accountActive = false;
    },
    tasksClickHandler() {
      this.homeActive = false;
      this.reservoirsActive = false;
      this.areasActive = false;
      this.materialsActive = false;
      this.cropsActive = false;
      this.tasksActive = true;
      this.accountActive = false;
    },
    accountClickHandler() {
      this.homeActive = false;
      this.reservoirsActive = false;
      this.areasActive = false;
      this.materialsActive = false;
      this.cropsActive = false;
      this.tasksActive = false;
      this.accountActive = true;
    },
  },
};
</script>

<style lang="scss" scoped>
.mobile-brand {
  width: 100px;
}

i.fa.fa-power-off {
  text-align: left;
  margin-right: 10px;
  width: 20px;
}

.bg-white {
  background-color: #fcfcfc !important;
}
</style>
