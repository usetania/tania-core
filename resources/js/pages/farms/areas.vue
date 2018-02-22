<template lang="pug">
  .areas.col
    modal(v-if="showModal" @close="showModal = false")
      farmAreaCreate(:data="data")
    .wrapper-md
      a.btn.m-b-xs.btn-primary.btn-addon.pull-right(style="cursor: pointer;" @click="openModal()")
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
                a.pull-right(style="cursor: pointer;" @click="openModal(area)")
                  i.fa.fa-edit
              small.text-muted {{ getType(area.type).label }}
            .panel-body
              .row
                .col-xs-4
                  small.text-muted.block Size ( {{ getSizeUnit(area.size.unit.symbol).label }} )
                  span.text-md {{ area.size.value }}
                .col-xs-4
                  small.text-muted.block Grow Batches
                  span.text-md {{ area.total_crop_batch }}
                .col-xs-4
                  small.text-muted.block Plant Quantity
                  span.text-md {{ area.plant_quantity }}

</template>

<script>
import { FindAreaType, FindAreaSizeUnit } from '@/stores/helpers/farms/area'
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal'
export default {
  name: "FarmAreas",
  computed: {
    ...mapGetters({
      areas: 'getAllAreas'
    })
  },
  components: {
    FarmAreaCreate: () => import('./areas-form.vue'),
    Modal
  },
  mounted () {
    this.fetchAreas()
  },
  data () {
    return {
      showModal: false,
      data: {},
    }
  },
  methods: {
    ...mapActions([
      'fetchAreas',
    ]),
    getType(key) {
      return FindAreaType(key)
    },
    getSizeUnit(key) {
      return FindAreaSizeUnit(key)
    },
    openModal(data) {
      this.showModal = true
      if (data) {
        this.data = data
      } else {
        this.data = {}
      }
    }
  }
}
</script>
