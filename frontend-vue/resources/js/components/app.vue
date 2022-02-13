<template lang="pug">
  #app
    .row.no-gutters(v-if="authenticated")
      .col-md-3.col-lg-2.d-none.d-md-block
        AppAsideComponent(:folded="folded" v-on:header-folded="setFolded" v-show="authenticated")

      .col-xs-12.col-sm-12.col-md-9.col-lg-10
        .main-content(:style="{ 'min-height': `${window.height}px` }")
          AppHeaderComponent(v-show="authenticated")

          .app-content-body
            router-view
          AppFooterComponent(v-show="authenticated")

    // When the user is not logged in
    .row.no-gutters(v-else)
      .col
        .main-content(:style="{ 'height': `${window.height}px` }")
          router-view
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  name: 'AppComponent',
  components: {
    AppHeaderComponent: () => import('./header.vue'),
    AppAsideComponent: () => import('./aside.vue'),
    AppFooterComponent: () => import('./footer.vue'),
  },
  data() {
    return {
      appReady: false,
      folded: false,
      window: {
        height: 0,
      },
    };
  },
  computed: {
    ...mapGetters({
      authenticated: 'IsUserAllowSeeNavigator',
    }),
  },
  created() {
    window.addEventListener('resize', this.handleResize);
    this.handleResize();
  },
  destroyed() {
    window.removeEventListener('resize', this.handleResize);
  },
  methods: {
    setFolded() {
      this.folded = !this.folded;
    },

    // Calculating browser height
    handleResize() {
      this.window.height = window.innerHeight;
    },
  },
};
</script>

<style lang="scss" scoped>
.main-content {
  background-color: #f6f8f8;
}
</style>
