<template lang="pug">
  .materials-create
    form(@submit.prevent="validateBeforeSubmit")
      .form-group
        label(for="name") Name
        input.form-control#name(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('name') }" v-model="inventory.name" name="name")
        span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
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
            label.control-labe Quantity
            .row
              .col-xs-6
                select.form-control#quantity_unit(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('quantity_unit') }" v-model="inventory.quantity_unit" name="quantity_unit")
                  option(v-for="unit in options.quantityUnits" v-bind:value="unit.key") {{ unit.label }}
                span.help-block.text-danger(v-show="errors.has('quantity_unit')") {{ errors.first('quantity_unit') }}
              .col-xs-6
                input.form-control#quantity(type="text" v-validate="'required|decimal|min:0'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="inventory.quantity" name="quantity")
                span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
          .col-xs-6
            label.control-label Produced by
            input.form-control#produced_by(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('produced_by') }" v-model="inventory.produced_by" name="produced_by")
            span.help-block.text-danger(v-show="errors.has('produced_by')") {{ errors.first('produced_by') }}
      .form-group
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
import { GrowingMediumQuantityUnits } from '@/stores/helpers/inventories/inventory'
import { mapGetters, mapActions } from 'vuex'
import moment from 'moment';
export default {
  name: 'InventoriesMaterialsCreateGrowingMedium',
  data () {
    return {
      inventory: Object.assign({}, StubInventory),
      options: {
        quantityUnits: Array.from(GrowingMediumQuantityUnits),
      }
    }
  },
  methods: {
    ...mapActions([
      'createMaterial',
    ]),
    create () {
      this.inventory.expiration_date = moment().format('YYYY-MM-DD')
      this.inventory.type = "growing_medium"
      this.createMaterial(this.inventory)
        .then(this.$emit('closeModal'))
        .catch(({ data }) => this.message = data)
    },
    closeModal () {
      this.$emit('closeModal')
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

