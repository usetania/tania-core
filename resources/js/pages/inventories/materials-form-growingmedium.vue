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
          label#label-quantity
            translate Quantity
          input.form-control#quantity(type="text" v-validate="'required|decimal|min:0'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="inventory.quantity" name="quantity")
          span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
        .col-6
          label#label-quantity-unit
            translate Unit
          select.form-control#quantity_unit(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('quantity_unit') }" v-model="inventory.quantity_unit" name="quantity_unit")
            option(v-for="unit in options.quantityUnits" v-bind:value="unit.key") {{ unit.label }}
          span.help-block.text-danger(v-show="errors.has('quantity_unit')") {{ errors.first('quantity_unit') }}
      .form-row
        .col-6
          label#label-price-per-unit
            translate Price per Unit
          b-input-group(:prepend="$gettext('€')")
            input.form-control#price_per_unit(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('price') }" v-model="inventory.price_per_unit" name="price")
          span.help-block.text-danger(v-show="errors.has('price')") {{ errors.first('price') }}
        .col-6
          label#label-produced-by
            translate Produced by
          input.form-control#produced_by(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('produced by') }" v-model="inventory.produced_by" name="produced by")
          span.help-block.text-danger(v-show="errors.has('produced by')") {{ errors.first('produced by') }}
      .form-group
        label#label-notes
          translate Additional Notes
        textarea.form-control#notes(type="text" :class="{'input': true, 'text-danger': errors.has('notes') }" v-model="inventory.notes" name="notes" rows="3")
        span.help-block.text-danger(v-show="errors.has('notes')") {{ errors.first('notes') }}
      .form-group
        BtnCancel(v-on:click.native="closeModal()")
        BtnSave(customClass="float-right")
</template>

<script>
import { mapActions } from 'vuex';
import moment from 'moment';
import { StubInventory } from '../../stores/stubs';
import { GrowingMediumQuantityUnits } from '../../stores/helpers/inventories/inventory';
import BtnCancel from '../../components/common/btn-cancel.vue';
import BtnSave from '../../components/common/btn-save.vue';

export default {
  name: 'InventoriesMaterialsFormGrowingMedium',
  components: {
    BtnCancel,
    BtnSave,
  },
  props: ['data'],
  data() {
    return {
      inventory: Object.assign({}, StubInventory),
      options: {
        quantityUnits: Array.from(GrowingMediumQuantityUnits),
      },
    };
  },
  mounted() {
    if (typeof this.data.uid !== 'undefined') {
      this.inventory.uid = this.data.uid;
      this.inventory.name = this.data.name;
      this.inventory.produced_by = this.data.produced_by;
      this.inventory.quantity = this.data.quantity.value;
      this.inventory.quantity_unit = this.data.quantity.unit;
      this.inventory.price_per_unit = this.data.price_per_unit.amount;
      this.inventory.notes = this.data.notes;
    } else {
      this.inventory.quantity_unit = this.options.quantityUnits[0].key;
    }
  },
  methods: {
    ...mapActions([
      'submitMaterial',
    ]),
    submit() {
      this.inventory.expiration_date = moment().format('YYYY-MM-DD');
      this.inventory.type = 'growing_medium';
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
