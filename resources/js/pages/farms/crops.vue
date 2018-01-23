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
          btn.btn-sm.btn-primary.btn-addon.pull-right(style="cursor: pointer;" id="show-modal" @click="showModal = true")
            i.fa.fa-plus
            | Add a New Batch
        table.table.m-b
          thead
            tr
              th(style="width: 13%") Crop Variety
              th(style="width: 12%") Batch ID
              th(style="width: 10%") Creation Date
              th(style="width: 10%") Activity Type
              th(style="width: 10%") Initial Area
              th(style="width: 10%") Current Area
          tbody
            tr(v-for="crop in crops")
              td: router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") {{ crop.inventory.variety }}
              td: span.identifier {{ crop.batch_id }}
              td {{ crop.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
              td {{ crop.type.code }}
              td: span.areatag {{ crop.initial_area.name }}
              td
                span.areatag(v-for="area in crop.current_area") {{ area.name }}
              td
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
