<template lang="pug">
  .materials-create
    form(@submit.prevent="validateBeforeSubmit")
      .form-group
        .row
          .col-xs-6
            label#label-name Name
            input.form-control#name(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('name') }" v-model="inventory.name" name="name")
            span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
          .col-xs-6
            label#label-produced-by Produced by
            input.form-control#produced_by(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('produced by') }" v-model="inventory.produced_by" name="produced by")
            span.help-block.text-danger(v-show="errors.has('produced by')") {{ errors.first('produced by') }}
      .form-group
        .row
          .col-xs-6
            label#label-price-per-unit Price per Unit
            .input-group.m-b
              span.input-group-addon &euro;
              input.form-control#price_per_unit(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('price') }" v-model="inventory.price_per_unit" name="price")
            span.help-block.text-danger(v-show="errors.has('price')") {{ errors.first('price') }}
          .col-xs-6
            label#label-quantity Quantity
            .input-group.m-b
              input.form-control#quantity(type="text" v-validate="'required|decimal|min:0'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="inventory.quantity" name="quantity")
              span.input-group-addon Pieces
            span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
      .form-group
        label#label-notes Additional Notes
        textarea.form-control#notes(type="text" :class="{'input': true, 'text-danger': errors.has('notes') }" v-model="inventory.notes" name="notes" rows="3")
        span.help-block.text-danger(v-show="errors.has('notes')") {{ errors.first('notes') }}
      .form-group
        button.btn.btn-addon.btn-success.pull-right(type="submit")
          i.fa.fa-plus
          | Save
        button.btn.btn-default(type="button" style="cursor: pointer;" @click="closeModal()") Cancel
</template>

<script>
import { StubInventory } from '@/stores/stubs'
import { mapGetters, mapActions } from 'vuex'
import moment from 'moment';
export default {
  name: 'InventoriesMaterialsFormOther',
  data () {
    return {
      inventory: Object.assign({}, StubInventory)
    }
  },
  methods: {
    ...mapActions([
      'submitMaterial',
    ]),
    submit () {
      this.inventory.expiration_date = moment().format('YYYY-MM-DD')
      this.inventory.type = "other"
      this.inventory.quantity_unit = "PIECES"
      this.submitMaterial(this.inventory)
        .then(() => this.$emit('closeModal'))
        .catch(() => this.$toasted.error('Error in material submission'))
    },
    closeModal () {
      this.$emit('closeModal')
    },
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.submit()
        }
      })
    }
  },
  mounted () {
    if (typeof this.data.uid != "undefined") {
      this.inventory.uid = this.data.uid
      this.inventory.name = this.data.name
      this.inventory.produced_by = this.data.produced_by
      this.inventory.quantity = this.data.quantity.value
      this.inventory.price_per_unit = this.data.price_per_unit.amount
      this.inventory.notes = this.data.notes
    }
  },
  props: ['data'],
}
</script>

