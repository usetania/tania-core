<template lang="pug">
  aside#aside.app-aside.hidden-xs.bg-dark
    .aside-wrap
      .navi-wrap
        nav.navi.clearfix
          ul.nav
            aside-item-component
              router-link(:to="{ name: 'Home' }")
                i.fa.fa-home
                span
                  translate Dashboard
            aside-item-component
              a.auto#production(@click.prevent="")
                span.pull-right.text-muted
                  i.fa.fa-fw.fa-angle-right.text
                  i.fa.fa.fw.fa-angle-down.text-active
                i.fa.fa-leaf
                span Production
              slot(name="item")
                ul.nav.nav-sub.dk
                  li: router-link#areas(:to="{name: 'FarmAreas'}")
                    span
                      translate Areas
                  li: router-link#crops(:to="{name: 'FarmCrops'}")
                    span
                      translate Crops
                  li: router-link#reservoirs(:to="{name: 'FarmReservoirs'}")
                    span
                      translate Reservoirs
            aside-item-component
              router-link#tasks(:to="{name: 'Task'}")
                i.fa.fa-clipboard
                span
                  translate Tasks
            aside-item-component
              a#inventories.auto(@click.prevent="")
                span.pull-right.text-muted
                  i.fa.fa-fw.fa-angle-right.text
                  i.fa.fa.fw.fa-angle-down.text-active
                i.fa.fa-archive
                span
                  translate Inventories
              slot(name="item")
                ul.nav.nav-sub.dk
                  li: router-link#materials(:to="{name: 'InventoriesMaterials'}")
                    span
                      translate Materials
            //li.hidden-folded.padder.m-t.m-b-sm.text-muted.text-xs
              span Settings
            aside-item-component
              router-link#account(:to="{ name: 'Account' }")
                i.fa.fa-user
                span
                  translate Account
</template>

<script>
import AsideItemComponent from './aside-item.vue';

export default {
  name: 'AppAsideComponent',
  components: {
    AsideItemComponent,
  },
  props: {
    folded: { type: Boolean, default: false },
  },
  methods: {
    setFolded() {
      this.$emit('header-folded');
    },
  },
};
</script>
