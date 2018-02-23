<template lang="pug">
  table.table.m-b
    thead
      tr
        th Crop Variety
        th Batch ID
        th(v-if="domain == 'AREA'") Seeding Date
        th Days Since Seeding
        th(v-if="domain == 'AREA'") Quantity
        th(v-if="domain == 'CROPS'") Initial Quantity
        th(v-if="domain == 'CROPS'") Status
        th(v-if="domain == 'AREA'") Last Watering
        th(v-if="domain == 'CROPS'")
    tbody
      tr(v-if="crops.length == 0")
        td(colspan="6") No Crops Available
      tr(v-for="crop in crops")
        td: router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") {{ crop.inventory.name }}
        td: span.identifier {{ crop.batch_id }}
        td(v-if="domain == 'AREA'") {{ crop.seeding_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
        td(v-if="domain == 'CROPS'")  {{ crop.initial_area.created_date | moment("from", new Date()) }}
        td(v-if="domain == 'AREA'")  {{ crop.days_since_seeding }}
        td(v-if="domain == 'AREA'") {{ crop.current_quantity }} {{ getCropContainer(crop.container.type, crop.container.quantity) }}
        td(v-if="domain == 'CROPS'") {{ crop.initial_area.initial_quantity }}
        td(v-if="domain == 'CROPS'")  {{ crop.area_status.seeding }} Seeding, {{ crop.area_status.growing }} Growing, {{ crop.area_status.dumped }} Dumped
        td(v-if="domain == 'AREA' && crop.last_watered") {{ crop.last_watered | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
        td(v-if="domain == 'AREA' && !crop.last_watered") -
        td(v-if="domain == 'CROPS'") 
          a.pull-right(style="cursor: pointer;" @click="editCropModal(crop)")
            i.fa.fa-edit
</template>

<script>
import { FindContainer } from '@/stores/helpers/farms/crop'
export default {
  name: "FarmCropsListing",
  methods: {
    editCropModal (crop) {
      this.$emit('editCrop', crop)
    },
    getCropContainer (key, count) {
      return FindContainer(key).label + ((count != 1)? 's':'')
    },
  },
  props: ['domain', 'crops'],
}
</script>