<template lang="pug">
.table-responsive
  table.table
    thead
      tr
        th
          translate Crop Variety
        th
          translate Batch ID
        th(v-if="domain == 'AREA'")
          translate Seeding Date
        th
          translate Days Since Seeding
        th(v-if="domain == 'AREA'")
          translate Quantity
        th(v-if="domain == 'CROPS'")
          translate Initial Quantity
        th(v-if="domain == 'CROPS'")
          translate Status
        th(v-if="domain == 'AREA'")
          translate Last Watering
        th(v-if="domain == 'CROPS'")
    tbody
      tr(v-if="crops.length == 0")
        td(colspan="6")
          translate No crops available.
      tr(v-for="crop in crops")
        td: router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") {{ crop.inventory.name }}
        td: span.identifier {{ crop.batch_id }}
        td(v-if="domain == 'AREA'") {{ crop.seeding_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
        td(v-if="domain == 'CROPS'")  {{ crop.initial_area.created_date | moment("from", new Date()) }}
        td(v-if="domain == 'AREA'")  {{ crop.days_since_seeding }}
        td(v-if="domain == 'AREA'") {{ crop.current_quantity }} {{ getCropContainer(crop.container.type, crop.container.quantity) }}
        td(v-if="domain == 'CROPS'") {{ crop.initial_area.initial_quantity }}
        td(v-if="domain == 'CROPS'")
          | {{ crop.area_status.seeding }}
          | &nbsp;
          translate Seeding
          | ,
          | &nbsp;
          | {{ crop.area_status.growing }}
          | &nbsp;
          translate Growing
          | ,
          | &nbsp;
          | {{ crop.area_status.dumped }}
          | &nbsp;
          translate Dumped
        td(v-if="domain == 'AREA' && crop.last_watered") {{ crop.last_watered | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
        td(v-if="domain == 'AREA' && !crop.last_watered") -
        td(v-if="domain == 'CROPS' && batch && crop.initial_area.initial_quantity == crop.initial_area.current_quantity")
          a.pull-right(style="cursor: pointer;" @click="editCropModal(crop)")
            i.fa.fa-edit
</template>

<script>
import { FindContainer } from '../../stores/helpers/farms/crop';

export default {
  name: 'FarmCropsListing',
  props: ['batch', 'crops', 'domain'],
  methods: {
    editCropModal(crop) {
      this.$emit('editCrop', crop);
    },
    getCropContainer(key, count) {
      return FindContainer(key).label + ((count !== 1) ? 's' : '');
    },
  },
};
</script>
