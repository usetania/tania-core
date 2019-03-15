<template lang="pug">
  .crops.col
    modal(v-if="showModal" @close="showModal = false")
      farmCropForm(:data="data")
    .wrapper-md
      h1.m-t.font-thin.h3.text-black Crops
      //.row
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
          li(role="presentation"  v-bind:class="{ active: isActive('BATCH') }")
            a(style="cursor: pointer;" @click="statusSelected('BATCH')") Batch
          li(role="presentation"  v-bind:class="{ active: isActive('ARCHIVES') }")
            a(style="cursor: pointer;" @click="statusSelected('ARCHIVES')") Archives
      .panel.no-border
        .panel-heading.wrapper.m-b
          span.h4.text-lt All Growing Batches on This Farm
          a.btn.btn-sm.btn-primary.btn-addon.pull-right(style="cursor: pointer;" id="show-modal" @click="showModal = true")
            i.fa.fa-plus
            | Add a New Batch
        FarmCropsListing(:crops="crops" :domain="'CROPS'" :batch="isActive('BATCH')" @editCrop="editCrop")
        Pagination(:pages="pages" :current="currentPage" @reload="getCrops")
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal.vue'
import Pagination from '@/components/pagination.vue'
export default {
  name: "FarmCrops",
  components: {
    FarmCropForm: () => import('./crops-form.vue'),
    FarmCropsListing: () => import('./crops-listing.vue'),
    Modal,
    Pagination,
  },
  computed: {
    ...mapGetters({
      cropInformation: 'getInformation',
      crops: 'getAllCrops',
      pages: 'getCropsNumberOfPages',
    })
  },
  data () {
    return {
      currentPage: 1,
      data: {},
      showModal: false,
      status: "BATCH",
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
    },
    getCrops () {
      if (this.status == 'BATCH') {
        this.fetchCrops({ pageId : this.getCurrentPage(), status : 'ACTIVE' })
      } else {
        this.fetchCrops({ pageId : this.getCurrentPage(), status : 'ARCHIVED' })
      }
    },
    getCurrentPage () {
      let pageId = 1
      if (typeof this.$route.query.page != "undefined") {
        pageId = parseInt(this.$route.query.page)
      }
      this.currentPage = pageId
      return pageId
    },
    statusSelected (status) {
      this.status = status
      this.$router.replace(this.$route.path)
      this.getCrops()
    },
    isActive (status) {
      return this.status == status
    }
  },
  mounted () {
    this.getCrops()
    this.getInformation()
  },
}
</script>
