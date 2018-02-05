<template lang="pug">

</template>
<template lang="pug">
  .materials-create
    form(@submit.prevent="validateBeforeSubmit")
      .form-group
        label(for="#name") Name
        input#name.form-control(type="text" name="name")
      .form-group
        .row
          .col-xs-6
            label(for="price") Price per Unit
            .input-group.m-b
              span.input-group-addon &euro;
              input#price.form-control(type="text" name="price")
          .col-xs-6
            label() Add this Expense?
            .radio
              label.i-checks.i-checks-sm
                input(type="radio" value="yes" checked)
                i
                | Yes
              label.i-checks.i-checks-sm
                input(type="radio" value="no")
                i
                | No
      .form-group
        .row
          .col-xs-6
            label.control-labe Quantity
            .input-group.m-b
              input.form-control(type="text")
              span.input-group-addon Pieces
          .col-xs-6
            label.control-label Produced by
            input.form-control(type="text")
      .form-group
        label.control-label Additional Notes
        textarea.form-control(rows="3")
      .form-group
        button.btn.btn-addon.btn-success.pull-right(type="submit")
          i.fa.fa-plus
          | Save
        button.btn.btn-default(style="cursor: pointer;" @click="closeModal()") Cancel
</template>

<script>
export default {
  name: 'InventoriesMaterialsCreateO',
  methods: {
    closeModal () {
      this.$emit('closeModal')
    },
    validateBeforeSubmit () {

    }
  }
}
</script>
