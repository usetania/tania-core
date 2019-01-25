<template lang="pug">
  .app.app-header-fixed(:class="folded == true ? 'app-aside-folded': ''")
    AppHeaderComponent(v-show="authenticated")
    AppAsideComponent(:folded="folded" v-on:header-folded="setFolded" v-show="authenticated")

    #content(role="main" :class="authenticated ? 'app-content': ''")
      .app-content-body.app-content-full
        .hbox.hbox-auto-xs.hbox-auto-sm.bg-light
          router-view
    AppFooterComponent(v-show="authenticated")
</template>

<script>
import { event } from '@/services/bus'
import { mapGetters } from 'vuex'
import AppHeaderComponent from '@/components/header.vue'
import AppAsideComponent from '@/components/aside.vue'
import AppFooterComponent from '@/components/footer.vue'

export default {
  name: 'AppComponent',
  data () {
    return {
      appReady: false,
      folded: false
    }
  },
  computed: {
    ...mapGetters({
      authenticated: 'IsUserAllowSeeNavigator'
    })
  },
  components: {
    AppHeaderComponent,
    AppAsideComponent,
    AppFooterComponent
  },

  methods: {
    setFolded (payload) {
      this.folded = !this.folded
    }
  }
}
</script>

