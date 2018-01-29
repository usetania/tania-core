<template lang="pug">
  .material.col
    .wrapper-md
      modal(v-if="showModal" @close="showModal = false")
        inventories-material-create
      a.btn.m-b-xs.btn-primary.btn-addon.pull-right(@click="showModal = true")
        i.fa.fa-plus
        |Add Material
      h1.m-t.font-thin.h3.text-black Materials

    .wrapper-md
      .panel.no-border
        table.table.m-b
          thead
            tr
              th(style="width: 15%") Category
              th(style="width: 20%") Name
              th(style="width: 10%") Price
              th(style="width: 15%") Produced By
              th(style="width: 12%") Quantity
          tbody

</template>

<script>
import InventoriesMaterialCreate from '@/pages/inventories/materials-create.vue'
import Modal from '@/components/modal.vue'
export default {
  name: 'InventoriesMaterial',
  data () {
    return {
      showModal: false
    }
  },
  components: {
    Modal, InventoriesMaterialCreate
  },

}
</script>
