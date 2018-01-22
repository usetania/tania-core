<template lang="pug">
  .crops.col
    FarmCropsCreate(v-if="showModal")
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
          btn.btn-sm.btn-primary.btn-addon.pull-right(v-on:click="openModal")
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
export default {
  name: "FarmCrops",
  components: {
    FarmCropsCreate: () => import('./crops-create.vue'),
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
    openModal: function () {
      this.showModal = true
    },
  }
}
</script>
