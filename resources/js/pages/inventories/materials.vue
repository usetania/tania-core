<template lang="pug">
  .material.col
    .wrapper-md
      modal(v-if="showModal" @close="showModal = false")
        InventoriesMaterialCreate
      a.btn.m-b-xs.btn-primary.btn-addon.pull-right(@click="showModal = true")
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
                a(href="#")
                  i.fa.fa-edit
</template>

<script>
import Modal from '@/components/modal.vue'
import { mapActions, mapGetters } from 'vuex'
import { FindInventoryType, FindQuantityUnit } from '@/stores/helpers/inventories/inventory'
export default {
  name: 'InventoriesMaterial',
  computed: {
    ...mapGetters({
      materials: 'getAllMaterials'
    })
  },
  components: {
    Modal,
    InventoriesMaterialCreate: () => import('./materials-create.vue'),
  },
  data () {
    return {
      showModal: false
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
  },
  mounted () {
    this.fetchMaterials()
  },
}
</script>
