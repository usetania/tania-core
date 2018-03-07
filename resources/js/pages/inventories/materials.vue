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
        .panel-footer(v-if="pages > 1")
          .text-center
            ul.pagination.pagination-sm.m-t-none.m-b-none
              li(v-bind:class="{ disabled: currentPage == 1 }")
                a(href="")
                  i.fa.fa-chevron-left
              li(v-if="pages < 6" v-for="pageNumber in pages")
                a(v-bind:class="{ active: currentPage == pageNumber }") {{ pageNumber }}
              li(v-bind:class="{ disabled: currentPage == pages }")
                a(href="")
                  i.fa.fa-chevron-right
</template>

<script>
import Modal from '@/components/modal.vue'
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
    Modal,
    InventoriesMaterialForm: () => import('./materials-form.vue'),
  },
  data () {
    return {
      data: {},
      currentPage: 1,
      showModal: false,
    }
  },
  methods: {
    ...mapActions([
      'fetchMaterials'
    ]),
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
    }
  },
  mounted () {
    this.fetchMaterials()
    console.log(this.pages)
  },
}
</script>
