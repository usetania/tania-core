<template lang="pug">
  .crops.col
    modal(v-if="showModal" @close="showModal = false")
      farmCropCreate
    .wrapper-md
      h1.m-t.font-thin.h3.text-black Crops
      .row
        .col-sm-3.m-t
          .hbox.bg-white-only.wrapper(style="min-height: 100px;")
            small.text-muted Harvested Produce This Month 
            a.pull-right(href=""): i.fa.fa-question-circle
            .h3.m-b.m-t 12.25 kilograms
        .col-sm-3.m-t
          .hbox.bg-white-only.wrapper(style="min-height: 100px;")
            small.text-muted Planted Varieties
            a.pull-right(href=""): i.fa.fa-question-circle
            .h3.m-b.m-t {{ crops.length }}
    .wrapper
      .m-b
        ul.nav.nav-tabs.h4
          li(role="presentation" class="active"): a(href="#") Batch
          li: a(href="#") Archives
      .panel.no-border
        .panel-heading.wrapper.m-b
          span.h4.text-lt All Growing Batches on This Farm
          a.btn.btn-sm.btn-primary.btn-addon.pull-right(style="cursor: pointer;" id="show-modal" @click="showModal = true")
            i.fa.fa-plus
            | Add a New Batch
        table.table.m-b
          thead
            tr
              th Crop Variety
              th Batch ID
              th Days Since Seeding
              th Initial Quantity
              th Status
          tbody
            tr(v-for="crop in crops")
              td: router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") {{ crop.inventory.name }}
              td: span.identifier {{ crop.batch_id }}
              td {{ crop.initial_area.created_date | moment("from", new Date()) }}
              td {{ crop.initial_area.initial_quantity }}
              td {{ crop.area_status.seeding }} Seeding, {{ crop.area_status.growing }} Growing, {{ crop.area_status.dumped }} Dumped

</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal'
export default {
  name: "FarmCrops",
  components: {
    FarmCropCreate: () => import('./crops-create.vue'),
    Modal
  },
  computed: {
    ...mapGetters({
      crops: 'getAllCrops'
    })
  },
  mounted () {
    this.fetchCrops()
  },
  data () {
    return {
      showModal: false
    }
  },
  methods: {
    ...mapActions([
      'fetchCrops'
    ]),
  }
}
</script>
