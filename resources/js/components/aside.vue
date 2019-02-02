<template lang="pug">
  aside#aside.app-aside.hidden-xs.bg-dark
    .aside-wrap
      .navi-wrap
        nav.navi.clearfix
          ul.nav
            aside-item-component
              router-link(:to="{ name: 'Home' }")
                i
                  font-awesome-icon.icon.text-primary(icon="home")
                span.hidden-folded.font-bold.text-primary Dashboard
            aside-item-component
              a.auto#production(@click.prevent="")
                span.hidden-folded.float-right.text-muted
                  i.fw.text
                    font-awesome-icon(icon="angle-right")
                  i.fw.text-active
                    font-awesome-icon(icon="angle-down")
                i
                  font-awesome-icon(icon="leaf")
                span.hidden-folded Production
              slot(name="item")
                ul.nav.nav-sub.dk
                  li: router-link#areas(:to="{name: 'FarmAreas'}")
                    span Areas
                  li: router-link#crops(:to="{name: 'FarmCrops'}")
                    span Crops
                  li: router-link#reservoirs(:to="{name: 'FarmReservoirs'}")
                    span Reservoirs
            aside-item-component
              router-link#tasks(:to="{name: 'Task'}")
                i
                  font-awesome-icon(icon="clipboard")
                span.hidden-folded Tasks
            aside-item-component
              a#inventories.auto(@click.prevent="")
                span.hidden-folded.float-right.text-muted
                  i.fw.text
                    font-awesome-icon(icon="angle-right")
                  i.fw.text-active
                    font-awesome-icon(icon="angle-down")
                i
                  font-awesome-icon(icon="archive")
                span.hidden-folded Inventories
              slot(name="item")
                ul.nav.nav-sub.dk
                  li: router-link#materials(:to="{name: 'InventoriesMaterials'}")
                    span.hidden-folded Materials
            li.hidden-folded.padder.m-t.m-b-sm.text-muted.text-xs
              span Settings
            aside-item-component
              router-link#account(:to="{ name: 'Account' }")
                i
                  font-awesome-icon(icon="user")
                span.hidden-folded Account
            li.line.dk
            .aside-item
              li
                a.no-shadow.navbar-btn(:class="folded == true ? 'active': ''" @click="setFolded")
                  i.text
                    font-awesome-icon(icon="chevron-circle-left")
                  i.text-active
                    font-awesome-icon(icon="chevron-circle-right")
                  span.hidden-folded.text-muted Collapse Menu

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
