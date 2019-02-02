<template lang="pug">
  nav#header.app-header.navbar.navbar-expand-lg.navbar-dark.bg-dark
    a.navbar-brand(href="/")
      img(src="../../images/logo.png")
      span.hidden-folded.m-l.xs Tania
    
    button.navbar-toggler(type="button", data-toggle="collapse", data-target="#navbarSupportedContent", aria-controls="navbarSupportedContent", aria-expanded="false", aria-label="Toggle navigation")
      span.navbar-toggler-icon

    #navbarSupportedContent.collapse.navbar-collapse
      ul.nav.navbar-nav.hidden-xs.mr-auto
        li.dropdown.farmswitch(:class="dropdown === true ? 'open': 'closed'")
          a.farm-current(href="#" @click.prevent="dropdownToggle")
            span {{ farm.name }}
            span.caret
          ul.dropdown-menu
            li.m-l.m-r.text-muted Switch Farm
            li(v-for="f in farms" :class="f.uid === farm.uid ? 'active': ''")
              a(href="#" @click.prevent="setFarm(f.uid)" :id="f.name")
                span
                  i.fa.fa-leaf(:class="{ 'text-success' : f.uid === farm.uid }")
                  | {{ f.name }}

      ul.navbar-nav
        li.nav-item
          a#signout.nav-link(href="#" @click.prevent="signout") Sign Out
</template>

<script>
import { ls } from '@/services'
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'AppHeaderComponent',
  data () {
    return {
      dropdown: false
    }
  },
  computed: {
    ...mapGetters({
      farm: 'getCurrentFarm',
      farms: 'getAllFarms'
    })
  },
  methods: {
    ...mapActions([
      'setCurrentFarm',
      'userSignOut'
    ]),
    setFarm (farmId) {
      this.setCurrentFarm(farmId)
        .then(data => this.dropdownToggle())
        .catch(error => console.log('Farm '+farmId+' is not found in the data'))
    },
    dropdownToggle() {
      this.dropdown = !this.dropdown
    },
    signout () {
      this.userSignOut()
        .then(data => {
          ls.remove('vuex')
          this.$router.push({ name: 'AuthLogin' })
        })
        .catch(err => console.log(error))
    }
  }
}
</script>

<style lang="scss" scoped>
  i.fa.fa-leaf {
    padding-right: 15px;
  }
</style>
