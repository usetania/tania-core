<template lang="pug">
  .container-fluid.bottom-space
    .row
      .col
        h3.title-page
          translate Crops

    .row
      .col
        modal(v-if="showModal" @close="showModal = false")
          farmCropForm(:data="data")

    b-card(no-body="no-body")
      b-tabs(card="card")
        b-tab(
          :title="$gettext('Batch')"
          v-bind:class="{ active: isActive('BATCH') }"
          @click="statusSelected('BATCH')"
        )
          BtnAddNew(:title="$gettext('Add New Batch')" v-on:click.native="showModal = true")

          FarmCropsListing(:crops="crops" :domain="'CROPS'" :batch="isActive('BATCH')" @editCrop="editCrop")
          Pagination(:pages="pages" :current="currentPage" @reload="getCrops")
        b-tab(
          :title="$gettext('Archives')"
          v-bind:class="{ active: isActive('ARCHIVES') }"
          @click="statusSelected('ARCHIVES')"
        )
          FarmCropsListing(:crops="crops" :domain="'CROPS'" :batch="isActive('ARCHIVES')" @editCrop="editCrop")
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import Modal from '../../components/modal.vue';
import Pagination from '../../components/pagination.vue';
import FarmCropForm from './crops-form.vue';
import FarmCropsListing from './crops-listing.vue';
import BtnAddNew from '../../components/common/btn-add-new.vue';

export default {
  name: 'FarmCrops',
  components: {
    FarmCropForm,
    FarmCropsListing,
    Modal,
    Pagination,
    BtnAddNew,
  },
  data() {
    return {
      currentPage: 1,
      data: {},
      showModal: false,
      status: 'BATCH',
    };
  },
  computed: {
    ...mapGetters({
      cropInformation: 'getInformation',
      crops: 'getAllCrops',
      pages: 'getCropsNumberOfPages',
    }),
  },
  mounted() {
    this.getCrops();
    this.getInformation();
  },
  methods: {
    ...mapActions([
      'fetchCrops',
      'getInformation',
    ]),
    editCrop(crop) {
      this.showModal = true;
      if (crop) {
        this.data = crop;
      } else {
        this.data = {};
      }
    },
    getCrops() {
      if (this.status === 'BATCH') {
        this.fetchCrops({ pageId: this.getCurrentPage(), status: 'ACTIVE' });
      } else {
        this.fetchCrops({ pageId: this.getCurrentPage(), status: 'ARCHIVED' });
      }
    },
    getCurrentPage() {
      let pageId = 1;
      if (typeof this.$route.query.page !== 'undefined') {
        pageId = parseInt(this.$route.query.page, 10);
      }
      this.currentPage = pageId;
      return pageId;
    },
    statusSelected(status) {
      this.status = status;
      this.$router.replace(this.$route.path);
      this.getCrops();
    },
    isActive(status) {
      return this.status === status;
    },
  },
};
</script>

<style lang="scss" scoped>
.title-page {
  margin: 20px 0 30px 0;
}

i.fa.fa-plus {
  text-align: left;
  width: 30px;
}

.bottom-space {
  padding-bottom: 60px;
}
</style>
