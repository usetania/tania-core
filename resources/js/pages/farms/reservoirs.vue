<template lang="pug">
  .reservoirs.col
    modal(v-if="showModal" @close="showModal = false")
      farmReservoirCreate
    .wrapper-md
      a.btn.m-b-xs.btn-primary.btn-addon.pull-right(style="cursor: pointer;" id="show-modal" @click="showModal = true")
        i.fa.fa-plus
        | Add Reservoir
      h1.m-n.font-thin.h3.text-black Water Reservoir
    .wrapper
      .panel.no-border
        .panel-heading.wrapper
          span.h4.text-lt All Reservoirs on This Farm
        table.table.m-b
          thead
            tr
              th Name
              th Created On
              th Source Type
              th Capacity
              th Used In
          tbody
            tr(v-for="reservoir in reservoirs")
              td: router-link(:to="{ name: 'FarmReservoir', params: { id: reservoir.uid } }") {{ reservoir.name }}
              td {{ reservoir.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
              td {{ getType(reservoir.water_source.type).label }}
              td {{ reservoir.water_source.capacity }}
              td
                span(v-for="(area, index) in reservoir.installed_to_areas")
                  | {{ area.name }}
                  span(v-if="index+1 < reservoir.installed_to_areas.length") , 
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
    FarmReservoirCreate: () => import('./reservoirs-create.vue'),
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
    }
  }
}
</script>


