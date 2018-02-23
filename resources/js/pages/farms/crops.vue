<template lang="pug">
  .crops.col
    modal(v-if="showModal" @close="showModal = false")
      farmCropForm(:data="data")
    .wrapper-md
      h1.m-t.font-thin.h3.text-black Crops
      .row
        .col-sm-3.m-t
          .hbox.bg-white-only.wrapper
            small.text-muted Harvested Produce This Month 
            a.pull-right(href=""): i.fa.fa-question-circle
            .h3.m-b.m-t {{ cropInformation.total_harvest_produced/1000 }} kilograms
        .col-sm-3.m-t
          .hbox.bg-white-only.wrapper
            small.text-muted Planted Varieties
            a.pull-right(href=""): i.fa.fa-question-circle
            .h3.m-b.m-t {{ cropInformation.total_plant_variety }}
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
        FarmCropsListing(:crops="crops" :domain="'CROPS'" @editCrop="editCrop")
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal'
export default {
  name: "FarmCrops",
  components: {
    FarmCropForm: () => import('./crops-form.vue'),
    FarmCropsListing: () => import('./crops-listing.vue'),
    Modal
  },
  computed: {
    ...mapGetters({
      crops: 'getAllCrops',
      cropInformation: 'getInformation'
    })
  },
  data () {
    return {
      showModal: false,
    }
  },
  methods: {
    ...mapActions([
      'fetchCrops',
      'getInformation',
    ]),
    editCrop (crop) {
      this.showModal = true
      if (crop) {
        this.data = crop
      } else {
        this.data = {}
      }
    }
  },
  mounted () {
    this.fetchCrops()
    this.getInformation()
  },
}
</script>
