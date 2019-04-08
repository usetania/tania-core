<template lang="pug">
  .container-fluid.bottom-space
    .row
      .col
        h3.title-page
          translate Materials

    .row
      .col
        modal(v-if="showModal" @close="showModal = false")
          InventoriesMaterialForm(:data="data")

        a#materialsform.btn.btn-primary(@click="openModal()")
          i.fa.fa-plus
          translate Add Material

    .table-responsive.table-wrapper
      table.table
        thead
          tr
            th
              translate Category
            th
              translate Name
            th
              translate Price
            th
              translate Produced By
            th
              translate Quantity
            th
              translate Additional Notes
            th
        tbody
          tr(v-if="materials.length == 0")
            td(colspan="7")
              translate No Inventories Available
          tr(v-for="material in materials")
            td {{ getType(material.type.code) }}
            td {{ material.name }}
            td {{ material.price_per_unit.amount }} {{ material.price_per_unit.symbol}}
            td {{ material.produced_by }}
            td {{ material.quantity.value }} {{ getQuantityUnit(material.quantity.unit) }}
            td {{ material.notes }}
            td
              a(@click="openModal(material)")
                i.fa.fa-edit
      Pagination(:pages="pages" @reload="getMaterials")
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import Modal from '../../components/modal.vue';
import Pagination from '../../components/pagination.vue';
import { FindInventoryType, FindQuantityUnit } from '../../stores/helpers/inventories/inventory';

export default {
  name: 'InventoriesMaterial',
  components: {
    InventoriesMaterialForm: () => import('./materials-form.vue'),
    Modal,
    Pagination,
  },
  data() {
    return {
      data: {},
      showModal: false,
    };
  },
  computed: {
    ...mapGetters({
      materials: 'getAllMaterials',
      pages: 'getMaterialsNumberOfPages',
    }),
  },
  mounted() {
    this.getMaterials();
  },
  methods: {
    ...mapActions([
      'fetchMaterials',
    ]),
    getMaterials() {
      let pageId = 1;
      if (typeof this.$route.query.page !== 'undefined') {
        pageId = parseInt(this.$route.query.page, 10);
      }
      this.fetchMaterials({ pageId });
    },
    getType(key) {
      return FindInventoryType(key);
    },
    getQuantityUnit(key) {
      return FindQuantityUnit(key);
    },
    openModal(data) {
      this.showModal = true;
      if (data) {
        this.data = data;
      } else {
        this.data = {};
      }
    },
  },
};
</script>

<style lang="scss" scoped>
h3.title-page {
  margin: 20px 0 30px 0;
}

i {
  text-align: left;
  width: 30px;
}

.table-wrapper {
  margin-top: 20px;
}

.bottom-space {
  padding-bottom: 60px;
}
</style>
