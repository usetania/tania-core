<template lang="pug">
  .materials-create
    b-form(@submit.prevent="validateBeforeSubmit")
      .form-group
        label#label-name
          translate Name
        input.form-control#name(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('name') }" v-model="inventory.name" name="name")
        span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
      .form-row
        .col-6
          label#label-price-per-unit
            translate Price per Unit
          b-input-group(:prepend="$gettext('€')")
            input.form-control#price_per_unit(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('price') }" v-model="inventory.price_per_unit" name="price")
          span.help-block.text-danger(v-show="errors.has('price')") {{ errors.first('price') }}
        .col-6
          label#label-quantity
            translate Quantity
          b-input-group(:append="$gettext('Pieces')")
            b-form-input#quantity(type="text" v-validate="'required|decimal|min:0'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="inventory.quantity" name="quantity")
          span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
      .form-group
        label#label-notes
          translate Additional Notes
        textarea.form-control#notes(type="text" :class="{'input': true, 'text-danger': errors.has('notes') }" v-model="inventory.notes" name="notes" rows="3")
        span.help-block.text-danger(v-show="errors.has('notes')") {{ errors.first('notes') }}
      .form-group
        button.btn.btn-addon.btn-success.float-right(type="submit")
          i.fa.fa-plus
          translate Save
        button.btn.btn-default(type="button" style="cursor: pointer;" @click="closeModal()")
          translate Cancel
</template>

<script>
import moment from 'moment';
import { mapActions } from 'vuex';
import { StubInventory } from '../../stores/stubs';

export default {
  name: 'InventoriesMaterialsFormLabelCrop',
  props: ['data'],
  data() {
    return {
      inventory: Object.assign({}, StubInventory),
    };
  },
  mounted() {
    if (typeof this.data.uid !== 'undefined') {
      this.inventory.uid = this.data.uid;
      this.inventory.name = this.data.name;
      this.inventory.produced_by = this.data.produced_by;
      this.inventory.quantity = this.data.quantity.value;
      this.inventory.price_per_unit = this.data.price_per_unit.amount;
      this.inventory.notes = this.data.notes;
    }
  },
  methods: {
    ...mapActions([
      'submitMaterial',
    ]),
    submit() {
      this.inventory.expiration_date = moment().format('YYYY-MM-DD');
      this.inventory.type = 'label_and_crop_support';
      this.inventory.quantity_unit = 'PIECES';
      this.submitMaterial(this.inventory)
        .then(() => this.$emit('closeModal'))
        .catch(() => this.$toasted.error('Error in material submission'));
    },
    closeModal() {
      this.$emit('closeModal');
    },
    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.submit();
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
i.fa.fa-plus {
  width: 30px;
}
</style>
