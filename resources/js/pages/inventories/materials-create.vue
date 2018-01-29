<template lang="pug">
  .materials-create
    .modal-header
      span.h4.font-bold Add Material to Inventory
    .modal-body
      p.text-muted Material is a consumable product using in your farm such as seeds, growing medium, fertilizer, pesticide, and so on.
      form(@submit.prevent="validateBeforeSubmit")
        .line.line-dashed.b-b.line-lg
        .form-group
          label.control-label Choose type of material
          select.form-control
            option
            option Seed
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit")
            i.fa.fa-long-arrow-right
            | Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>

<script>
export default {
  name: 'InventoriesMaterialsCreate',
  methods: {
    validateBeforeSubmit () {

    }
  }
}
</script>
