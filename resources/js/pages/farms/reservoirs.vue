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
              th(style="witdh: 25%") Name
              th(style="width: 15%") Date of Creation
              th(style="width: 20%") Source Type
              th Capacity
              th Connected To
              th Tasks
          tbody
            tr(v-for="reservoir in reservoirs")
              td {{ reservoir.name }}
              td {{ reservoir.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY HH:mm:ss') }}
              td {{ getType(reservoir.water_source.type).label }}
              td {{ reservoir.water_source.capacity }}
              td
                span(v-for="area in reservoir.installed_to_areas") {{ area.name }}
              td
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


