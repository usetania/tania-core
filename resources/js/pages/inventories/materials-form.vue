<template lang="pug">
  .materials-form
    .modal-header
      span.h4.font-bold(v-if="data.uid")
        translate Update Material
      span.h4.font-bold(v-else)
        translate Add Material to Inventory
    .modal-body
      p.text-muted
        translate Material is a consumable product using in your farm such as seeds, growing medium, fertilizer, pesticide, and so on.
      form
        .line.line-dashed.b-b.line-lg
        .form-group
          label#label-material-type
            translate Choose type of material
          select.form-control#material_type(@change="typeChanged($event.target.value)" v-model="materialType")
            option(v-for="inventory in options.inventoryTypes" v-bind:value="inventory.key") {{ inventory.label }}
        InventoriesMaterialFormAgrochemical(v-if="showAgrochemical" @closeModal="closeModal" :data="data")
        InventoriesMaterialFormGrowingMedium(v-if="showGrowingMedium" @closeModal="closeModal" :data="data")
        InventoriesMaterialFormLabelCrop(v-if="showLabelCrop" @closeModal="closeModal" :data="data")
        InventoriesMaterialFormOther(v-if="showOther" @closeModal="closeModal" :data="data")
        InventoriesMaterialFormPlant(v-if="showPlant" @closeModal="closeModal" :data="data")
        InventoriesMaterialFormPotHarvest(v-if="showPotHarvest" @closeModal="closeModal" :data="data")
        InventoriesMaterialFormSeed(v-if="showSeed" @closeModal="closeModal" :data="data")
        InventoriesMaterialFormSeedContainer(v-if="showSeedContainer" @closeModal="closeModal" :data="data")
</template>

<script>
import { InventoryTypes } from '../../stores/helpers/inventories/inventory'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: 'InventoriesMaterialsForm',
  components: {
    InventoriesMaterialFormAgrochemical: () => import('./materials-form-agrochemical.vue'),
    InventoriesMaterialFormGrowingMedium: () => import('./materials-form-growingmedium.vue'),
    InventoriesMaterialFormLabelCrop: () => import('./materials-form-labelcrop.vue'),
    InventoriesMaterialFormOther: () => import('./materials-form-other.vue'),
    InventoriesMaterialFormPlant: () => import('./materials-form-plant.vue'),
    InventoriesMaterialFormPotHarvest: () => import('./materials-form-potharvest.vue'),
    InventoriesMaterialFormSeed: () => import('./materials-form-seed.vue'),
    InventoriesMaterialFormSeedContainer: () => import('./materials-form-seedcontainer.vue')
  },
  data () {
    return {
      showAgrochemical: false,
      showGrowingMedium: false,
      showLabelCrop: false,
      showOther: false,
      showPlant: false,
      showPotHarvest: false,
      showSeed: false,
      showSeedContainer: false,
      materialType: 'seed',
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
      this.materialType = type
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
  },
  mounted () {
    if (typeof this.data.uid != "undefined") {
      this.typeChanged(this.data.type.code.toLowerCase())
    } else {
      this.typeChanged("seed")
    }
  },
  props: ['data'],
}
</script>
