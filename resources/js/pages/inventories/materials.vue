<template lang="pug">
  .material.col
    .wrapper-md
      modal(v-if="showModal" @close="showModal = false")
        InventoriesMaterialForm(:data="data")
      a.btn.m-b-xs.btn-primary.btn-addon.pull-right(@click="openModal()")
        i.fa.fa-plus
        |Add Material
      h1.m-t.font-thin.h3.text-black Materials

    .wrapper-md
      .panel.no-border
        table.table.m-b
          thead
            tr
              th Category
              th Name
              th Price
              th Produced By
              th Quantity
              th Additional Notes
              th 
          tbody
            tr(v-if="materials.length == 0")
              td(colspan="7") No Materials Available
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
import Modal from '@/components/modal.vue'
import Pagination from '@/components/pagination.vue'
import { mapActions, mapGetters } from 'vuex'
import { FindInventoryType, FindQuantityUnit } from '@/stores/helpers/inventories/inventory'
export default {
  name: 'InventoriesMaterial',
  computed: {
    ...mapGetters({
      materials: 'getAllMaterials',
      pages: 'getMaterialsNumberOfPages',
    })
  },
  components: {
    InventoriesMaterialForm: () => import('./materials-form.vue'),
    Modal,
    Pagination,
  },
  data () {
    return {
      data: {},
      showModal: false,
    }
  },
  methods: {
    ...mapActions([
      'fetchMaterials'
    ]),
    getMaterials() {
      let pageId = 1
      if (typeof this.$route.query.page != "undefined") {
        pageId = parseInt(this.$route.query.page)
      }
      this.fetchMaterials({ pageId : pageId })
    },
    getType(key) {
      return FindInventoryType(key)
    },
    getQuantityUnit(key) {
      return FindQuantityUnit(key)
    },
    openModal(data) {
      this.showModal = true
      if (data) {
        this.data = data
      } else {
        this.data = {}
      }
    },
  },
  mounted () {
    this.getMaterials()
  },
}
</script>
