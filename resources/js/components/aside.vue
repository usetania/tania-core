<template lang="pug">
  aside#aside.app-aside.hidden-xs.bg-dark
    .aside-wrap
      .navi-wrap
        nav.navi.clearfix
          ul.nav
            aside-item-component
              router-link(:to="{ name: 'Home' }")
                i.fa.fa-home.icon.text-primary
                span.font-bold.text-primary Dashboard
            aside-item-component
              a.auto(@click.prevent="")
                span.pull-right.text-muted
                  i.fa.fa-fw.fa-angle-right.text
                  i.fa.fa.fw.fa-angle-down.text-active
                i.fa.fa-leaf
                span Production
              slot(name="item")
                ul.nav.nav-sub.dk
                  li <router-link :to="{ name: 'FarmAreas'}"><span>Areas</span></router-link>
                  li <router-link :to="{ name: 'FarmCrops'}"><span>Crops</span></router-link>
                  li <router-link :to="{ name: 'FarmReservoirs'}"><span>Reservoir</span></router-link>
            aside-item-component
              router-link(:to="{name: 'Task'}")
                i.fa.fa-clipboard
                span Tasks
            aside-item-component
              a.auto(@click.prevent="")
                span.pull-right.text-muted
                  i.fa.fa-fw.fa-angle-right.text
                  i.fa.fa.fw.fa-angle-down.text-active
                i.fa.fa-archive
                span Inventories
              slot(name="item")
                ul.nav.nav-sub.dk
                  li <router-link :to="{ name: 'InventoriesMaterials'}"><span>Material</span></router-link>
                  li <a><span>Tools</span></a>
            li.hidden-folded.padder.m-t.m-b-sm.text-muted.text-xs
              span Settings
            aside-item-component
              a
                i.fa.fa-cog
                span General
            aside-item-component
              a
                i.fa.fa-user
                span Account
            aside-item-component
              a
                i.fa.fa-wrench
                span Preferences
            li.line.dk
            li
              a.btn.no-shadow.navbar-btn(:class="folded == true ? 'active': ''" @click="setFolded")
                i.fa.fa-chevron-circle-left.text
                i.fa.fa-chevron-circle-right.text-active
                span.text-muted Collapse Menu

</template>

<script>
import AsideItemComponent from './aside-item.vue'
export default {
  name: 'AppAsideComponent',
  components: {
    AsideItemComponent
  },
  props: {
    folded: { type: Boolean, default: false }
  },
  methods: {
    setFolded () {
      this.$emit('header-folded')
    }
  }
}
</script>
