<template lang="pug">
  .materials-create
    .modal-header
      span.h4.font-bold Add Material to Inventory
    .modal-body
      p.text-muted Material is a consumable product using in your farm such as seeds, growing medium, fertilizer, pesticide, and so on.
      form
        .line.line-dashed.b-b.line-lg
        .form-group
          label.control-label Choose type of material
          select.form-control#material_type(@change="typeChanged($event.target.value)")
            option(v-for="inventory in options.inventoryTypes" v-bind:value="inventory.key") {{ inventory.label }}
        InventoriesMaterialCreateAgrochemical(v-if="showAgrochemical" @closeModal="closeModal")
        InventoriesMaterialCreateGrowingMedium(v-if="showGrowingMedium" @closeModal="closeModal")
        InventoriesMaterialCreateLabelCrop(v-if="showLabelCrop" @closeModal="closeModal")
        InventoriesMaterialCreateOther(v-if="showOther" @closeModal="closeModal")
        InventoriesMaterialCreatePlant(v-if="showPlant" @closeModal="closeModal")
        InventoriesMaterialCreatePotHarvest(v-if="showPotHarvest" @closeModal="closeModal")
        InventoriesMaterialCreateSeed(v-if="showSeed" @closeModal="closeModal")
        InventoriesMaterialCreateSeedContainer(v-if="showSeedContainer" @closeModal="closeModal")
</template>

<script>
import { InventoryTypes } from '@/stores/helpers/inventories/inventory'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: 'InventoriesMaterialsCreate',
  components: {
    InventoriesMaterialCreateAgrochemical: () => import('./materials-create-agrochemical.vue'),
    InventoriesMaterialCreateGrowingMedium: () => import('./materials-create-growingmedium.vue'),
    InventoriesMaterialCreateLabelCrop: () => import('./materials-create-labelcrop.vue'),
    InventoriesMaterialCreateOther: () => import('./materials-create-other.vue'),
    InventoriesMaterialCreatePlant: () => import('./materials-create-plant.vue'),
    InventoriesMaterialCreatePotHarvest: () => import('./materials-create-potharvest.vue'),
    InventoriesMaterialCreateSeed: () => import('./materials-create-seed.vue'),
    InventoriesMaterialCreateSeedContainer: () => import('./materials-create-seedcontainer.vue')
  },
  data () {
    return {
      showAgrochemical: false,
      showGrowingMedium: false,
      showLabelCrop: false,
      showOther: false,
      showPlant: false,
      showPotHarvest: false,
      showSeed: true,
      showSeedContainer: false,
      options: {
        inventoryTypes: Array.from(InventoryTypes),
      }
    }
  },
  methods: {
    ...mapActions([
      'typeChanged',
    ]),
    closeModal () {
      this.$parent.$emit('close')
    },
    typeChanged (type) {
      this.showAgrochemical = false
      this.showGrowingMedium = false
      this.showLabelCrop = false
      this.showOther = false
      this.showPlant = false
      this.showPotHarvest = false
      this.showSeed = false
      this.showSeedContainer = false
      if (type == "seed") {
        this.showSeed = true
      } else if (type == "growing_medium") {
        this.showGrowingMedium = true
      } else if (type == "agrochemical") {
        this.showAgrochemical = true
      } else if (type == "label_and_crop_support") {
        this.showLabelCrop = true
      } else if (type == "seeding_container") {
        this.showSeedContainer = true
      } else if (type == "post_harvest_supply") {
        this.showPotHarvest = true
      } else if (type == "other") {
        this.showOther = true
      } else if (type == "plant") {
        this.showPlant = true
      }
    }
  }
}
</script>
