<template lang="pug">
  .areas.col
    .wrapper-md
      router-link.btn.m-b-xs.btn-primary.btn-addon.pull-right(:to="{name: 'FarmAreasCreate'}")
        i.fa.fa-plus
        | Add Area
      h1.m-n.font-thin.h3.text-black Areas
    .wrapper-md
      .row
        .col-md-4.col-xs-12(v-for="area in areas")
          .panel.no-border
            .panel-heading.description
              .h3.text-lt.name
                router-link(:to="{ name: 'FarmArea', params: { id: area.uid } }") {{ area.name }}
              small.text-muted {{ getType(area.type).label }}
            .panel-body
              .row
                .col-xs-4
                  small.text-muted.block Size ( {{ getSizeUnit(area.size.symbol).label }} )
                  span.text-md {{ area.size.value }}
                .col-xs-4
                  small.text-muted.block Grow Batches
                  span.text-md 21
                .col-xs-4
                  small.text-muted.block Plant Quantity
                  span.text-md 1680

</template>

<script>
import { FindAreaType, FindAreaSizeUnit } from '@/stores/helpers/farms/area'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: "FarmAreas",
  computed: {
    ...mapGetters({
      areas: 'getAllAreas'
    })
  },
  mounted () {
    this.fetchAreas()
  },
  methods: {
    ...mapActions([
      'fetchAreas'
    ]),
    getType(key) {
      return FindAreaType(key)
    },
    getSizeUnit(key) {
      return FindAreaSizeUnit(key)
    }
  }
}
</script>


