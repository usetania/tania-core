<template lang="pug">
  header#header.app-header.navbar(role="menu")
    .navbar-header.bg-dark
      button.pull-right.visible-xs.dk()
      button.pull-right.visible-xs()
      a.navbar-brand.text-lt(href="/")
        img(src="../../images/logo.png")
    .collapse.pos-rlt.navbar-collapse.box-shadow.bg-white-only
      ul.nav.navbar-nav.navbar-right
        li
          a#signout(href="#" @click.prevent="signout") Sign Out
      ul.nav.navbar-nav.hidden-xs
        li
          a {{ farm.name }}
        //li.dropdown.farmswitch(:class="dropdown === true ? 'open': 'closed'")
          a.farm-current(href="#" @click.prevent="dropdownToggle")
            span {{ farm.name }}
            span.caret
          ul.dropdown-menu
            li.m-l.m-r.text-muted Switch Farm
            li(v-for="f in farms" :class="f.uid === farm.uid ? 'active': ''")
              a(href="#" @click.prevent="setFarm(f.uid)" :id="f.name")
                span
                  i.fa.fa-leaf(:class="f.uid === farm.uid ? 'text-success': ''")
                  | {{ f.name }}
</template>

<script>
import { ls } from '../services'
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

<style type="sccs" scoped>
  i.fa.fa-leaf {
    padding-right: 15px;
  }
</style>
