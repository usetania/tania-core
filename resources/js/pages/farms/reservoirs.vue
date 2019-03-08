<template lang="pug">
  .reservoirs.col
    modal(v-if="showModal" @close="showModal = false")
      FarmReservoirForm(:data="data")
    .wrapper-md
      a#reservoirsform.btn.m-b-xs.btn-primary.btn-addon.pull-right(style="cursor: pointer;" @click="openModal()")
        i.fa.fa-plus
        | Add Reservoir
      h1.m-n.font-thin.h3.text-black Water Reservoir
    .wrapper
      .panel.no-border
        table.table.m-b
          thead
            tr
              th Name
              th Created On
              th Source Type
              th Capacity
              th Used In
              th
          tbody
            tr(v-for="reservoir in reservoirs")
              td: router-link(:to="{ name: 'FarmReservoir', params: { id: reservoir.uid } }")
                u {{ reservoir.name }}
              td {{ reservoir.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
              td {{ getType(reservoir.water_source.type).label }}
              td {{ reservoir.water_source.capacity }}
              td
                span(v-for="(area, index) in reservoir.installed_to_area")
                  | {{ area.name }}
                  span(v-if="index+1 < reservoir.installed_to_area.length") , 
              td
                a.pull-right(style="cursor: pointer;" @click="openModal(reservoir)")
                  i.fa.fa-edit
</template>

<script>
import { FindReservoirType } from '@/stores/helpers/farms/reservoir'
import { mapGetters, mapActions } from 'vuex'
import Modal from '@/components/modal'
export default {
  name: "FarmReservoirs",
  computed : {
    ...mapGetters({
      reservoirs: 'getAllReservoirs'
    })
  },
  components: {
    FarmReservoirForm: () => import('./reservoirs-form.vue'),
    Modal
  },
  data () {
    return {
      showModal: false
    }
  },
  mounted () {
    this.fetchReservoirs()
  },
  methods: {
    ...mapActions([
      'fetchReservoirs'
    ]),
    getType(key) {
      return FindReservoirType(key)
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


