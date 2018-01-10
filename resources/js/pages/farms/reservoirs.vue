<template lang="pug">
  .reservoirs.col
    .wrapper-md
      router-link.btn.m-b-xs.btn-primary.btn-addon.pull-right(:to="{name: 'FarmReservoirsCreate'}")
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
export default {
  name: "FarmReservoirs",
  computed : {
    ...mapGetters({
      reservoirs: 'getAllReservoirs'
    })
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


