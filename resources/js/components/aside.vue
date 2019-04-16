<template lang="pug">
aside#aside
  .brand-image.text-center
    img(src="../../images/logo.png" alt="Tania Logo")

  b-nav(vertical="vertical")
    AsideItem(
      :title="$gettext('Dashboard')"
      fontawesome="fa fa-home"
      :routeName="{ name: 'Home' }"
      :isActive="homeActive"
      v-on:click.native="homeClickHandler"
    )

    AsideItem(
      :title="$gettext('Reservoirs')"
      fontawesome="fa fa-tint"
      :routeName="{ name: 'FarmReservoirs' }"
      :isActive="reservoirsActive"
      v-on:click.native="reservoirsClickHandler"
    )

    AsideItem(
      :title="$gettext('Areas')"
      fontawesome="fa fa-grip-horizontal"
      :routeName="{ name: 'FarmAreas' }"
      :isActive="areasActive"
      v-on:click.native="areasClickHandler"
    )

    AsideItem(
      :title="$gettext('Materials')"
      fontawesome="fa fa-archive"
      :routeName="{ name: 'InventoriesMaterials' }"
      :isActive="materialsActive"
      v-on:click.native="materialsClickHandler"
    )

    AsideItem(
      :title="$gettext('Crops')"
      fontawesome="fa fa-leaf"
      :routeName="{ name: 'FarmCrops' }"
      :isActive="cropsActive"
      v-on:click.native="cropsClickHandler"
    )

    AsideItem(
      :title="$gettext('Tasks')"
      fontawesome="fa fa-clipboard"
      :routeName="{ name: 'Task' }"
      :isActive="tasksActive"
      v-on:click.native="tasksClickHandler"
    )

    AsideItem(
      :title="$gettext('Account')"
      fontawesome="fa fa-user"
      :routeName="{ name: 'Account' }"
      :isActive="accountActive"
      v-on:click.native="accountClickHandler"
    )
</template>

<script>
import AsideItem from './common/aside-item.vue';

export default {
  name: 'AppAsideComponent',
  components: {
    AsideItem,
  },
  data() {
    return {
      homeActive: true,
      reservoirsActive: false,
      areasActive: false,
      materialsActive: false,
      cropsActive: false,
      tasksActive: false,
      accountActive: false,
    };
  },
  methods: {
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
aside#aside {
  background-color: #513969;
  height: 100%;

  .brand-image {
    padding: 15px 0;
    margin-bottom: 30px;

    img {
      width: 100px;
    }
  }
}
</style>
