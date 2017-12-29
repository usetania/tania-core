<template lang="pug">
  .app.app-header-fixed(:class="folded == true ? 'app-aside-folded': ''")
    AppHeaderComponent(:folded="folded" v-on:header-folded="setFolded" v-show="authenticated")
    AppAsideComponent(v-show="authenticated")

    #content(role="main" :class="authenticated ? 'app-content': ''")
      .app-content-body
        .hbox.hbox-auto-xs.hbox-auto-sm
          router-view

    AppFooterComponent(v-if="authenticated")
</template>

<script>
import { event } from '@/services'
import { mapActions, mapGetters } from 'vuex'
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
    AppHeaderComponent: () => import('./header'),
    AppAsideComponent: () => import('./aside'),
    AppFooterComponent: () => import('./footer')
  },

  mounted () {
    this.init()
  },

  created () {
    event.on({
      'tanibox:ready': () => {
        this.appReady = true
        // console.log('tanibox is ready')
      }
    })
  },

  methods: {
    ...mapActions([
      'fetchCountries',
      'fetchFarmTypes'
    ]),
    init () {
      try {
        // perform http initial request location services ?
        Promise.all([
          this.fetchCountries(),
          this.fetchFarmTypes()
        ]).then(response => {
          event.emit('tanibox:ready')
        })
      } catch (err) {

      }
    },

    setFolded (payload) {
      this.folded = !this.folded
    }
  }
}
</script>

