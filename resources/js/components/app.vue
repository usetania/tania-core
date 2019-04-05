<template lang="pug">
  .app.app-header-fixed(:class="folded == true ? 'app-aside-folded': ''")
    AppHeaderComponent(v-show="authenticated")
    AppAsideComponent(:folded="folded" v-on:header-folded="setFolded" v-show="authenticated")

    #content(role="main" :class="authenticated ? 'app-content': ''")
      .app-content-body.app-content-full
        .hbox.hbox-auto-xs.hbox-auto-sm
          router-view
    AppFooterComponent(v-show="authenticated")
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
    };
  },
  computed: {
    ...mapGetters({
      authenticated: 'IsUserAllowSeeNavigator',
    }),
  },

  methods: {
    setFolded() {
      this.folded = !this.folded;
    },
  },
};
</script>
