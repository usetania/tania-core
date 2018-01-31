<template lang="pug">
  .material.col
    .wrapper-md
      modal(v-if="showModal" @close="showModal = false")
        inventories-material-create-sc
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
import Modal from '@/components/modal.vue'
export default {
  name: 'InventoriesMaterial',
  data () {
    return {
      showModal: false
    }
  },
  components: {
    Modal,
    InventoriesMaterialCreate: () => import('./materials-create.vue'),
    InventoriesMaterialCreateAc: () => import('./materials-create-ac.vue'),
    InventoriesMaterialCreateGm: () => import('./materials-create-gm.vue'),
    InventoriesMaterialCreateLc: () => import('./materials-create-lc.vue'),
    InventoriesMaterialCreateO: () => import('./materials-create-o.vue'),
    InventoriesMaterialCreatePh: () => import('./materials-create-ph.vue'),
    InventoriesMaterialCreateS: () => import('./materials-create-s.vue'),
    InventoriesMaterialCreateSc: () => import('./materials-create-sc.vue')
  },

}
</script>
