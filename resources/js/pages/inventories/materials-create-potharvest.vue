<template lang="pug">
  .materials-create
    form(@submit.prevent="validateBeforeSubmit")
      .form-group
        .row
          .col-xs-6
            label(for="name") Name
            input.form-control#name(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('name') }" v-model="inventory.name" name="name")
            span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
          .col-xs-6
            label.control-label Produced by
            input.form-control#produced_by(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('produced_by') }" v-model="inventory.produced_by" name="produced_by")
            span.help-block.text-danger(v-show="errors.has('produced_by')") {{ errors.first('produced_by') }}
      .form-group
        .row
          .col-xs-6
            label(for="price_per_unit") Price per Unit
            .input-group.m-b
              span.input-group-addon &euro;
              input.form-control#price_per_unit(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('price_per_unit') }" v-model="inventory.price_per_unit" name="price_per_unit")
              span.help-block.text-danger(v-show="errors.has('price_per_unit')") {{ errors.first('price_per_unit') }}
          .col-xs-6
            label(for="is_expense") Add this Expense?
            .radio
              label.i-checks.i-checks-sm
                input#is_expense(type="radio" name="is_expense" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('is_expense') }" v-model="inventory.is_expense" value="true")
                i
                | Yes
              label.i-checks.i-checks-sm
                input#is_expense(type="radio" name="is_expense" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('is_expense') }" v-model="inventory.is_expense" value="false")
                i
                | No
            span.help-block.text-danger(v-show="errors.has('is_expense')") {{ errors.first('is_expense') }}
      .form-group
        .row
          .col-xs-6
            label.control-label(for="quantity") Quantity
            .input-group.m-b
              input.form-control#quantity(type="text" v-validate="'required|numeric|min:0'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="inventory.quantity" name="quantity")
              span.input-group-addon Pieces
            span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
          .col-xs-6
            label.control-label(for="notes") Additional Notes
            textarea.form-control#notes(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('notes') }" v-model="inventory.notes" name="notes" rows="3")
            span.help-block.text-danger(v-show="errors.has('notes')") {{ errors.first('notes') }}
      .form-group
        button.btn.btn-addon.btn-success.pull-right(type="submit")
          i.fa.fa-plus
          | Save
        button.btn.btn-default(style="cursor: pointer;" @click="closeModal()") Cancel
</template>

<script>
import { StubInventory } from '@/stores/stubs'
import { mapGetters, mapActions } from 'vuex'
import moment from 'moment';
export default {
  name: 'InventoriesMaterialsCreatePotHarvest',
  data () {
    return {
      inventory: Object.assign({}, StubInventory)
    }
  },
  methods: {
    ...mapActions([
      'createMaterial',
      'openPicker',
    ]),
    create () {
      this.inventory.expiration_date = moment().format('YYYY-MM-DD')
      this.inventory.type = "post_harvest_supply"
      this.inventory.quantity_unit = "PIECES"
      this.createMaterial(this.inventory)
        .then(this.$emit('closeModal'))
        .catch(({ data }) => this.message = data)
    },
    closeModal () {
      this.$emit('closeModal')
    },
    openPicker () {
      this.$refs.openCal.showCalendar()
    },
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    }
  }
}
</script>

