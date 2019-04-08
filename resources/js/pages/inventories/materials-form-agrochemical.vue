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
          label#label-chemical-type
            translate Chemical Type
          select.form-control#chemical_type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('chemical type') }" v-model="inventory.chemical_type" name="chemical type")
            option(v-for="chemical in options.chemicalTypes" v-bind:value="chemical.key") {{ chemical.label }}
          span.help-block.text-danger(v-show="errors.has('chemical type')") {{ errors.first('chemical type') }}
        .col-6
          label#label-price-per-unit
            translate Price per Unit
          b-input-group(:prepend="$gettext('€')")
            input.form-control#price_per_unit(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('price') }" v-model="inventory.price_per_unit" name="price")
          span.help-block.text-danger(v-show="errors.has('price')") {{ errors.first('price') }}
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
        .col-xs-12.col-sm-12.col-md-7
          label#label-expiration-date
            translate Expiration date
          .input-group
            datepicker#expiration_date(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('expiration date') }" v-model="inventory.expiration_date" name="expiration date" input-class="form-control" ref="openCal")
            span.input-group-btn
              button.btn.btn-primary(type="button" v-on:click="openPicker")
                i.fa.fa-calendar
          span.help-block.text-danger(v-show="errors.has('expiration date')") {{ errors.first('expiration date') }}
        .col-xs-12.col-sm-12.col-md-5
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
        button.btn.btn-addon.btn-success.pull-right(type="submit")
          i.fa.fa-plus
          translate Save
        button.btn.btn-default(type="button" style="cursor: pointer;" @click="closeModal()")
          translate Cancel
</template>

<script>
import { mapActions } from 'vuex';
import Datepicker from 'vuejs-datepicker';
import moment from 'moment';
import { StubInventory } from '../../stores/stubs';
import { AgrochemicalQuantityUnits, ChemicalTypes } from '../../stores/helpers/inventories/inventory';

export default {
  name: 'InventoriesMaterialsFormAgrochemical',
  components: {
    Datepicker,
  },
  props: ['data'],
  data() {
    return {
      inventory: Object.assign({}, StubInventory),
      options: {
        chemicalTypes: Array.from(ChemicalTypes),
        quantityUnits: Array.from(AgrochemicalQuantityUnits),
      },
    };
  },
  mounted() {
    if (typeof this.data.uid !== 'undefined') {
      this.inventory.uid = this.data.uid;
      this.inventory.name = this.data.name;
      this.inventory.chemical_type = this.data.type.type_detail.chemical_type.code;
      this.inventory.produced_by = this.data.produced_by;
      this.inventory.quantity = this.data.quantity.value;
      this.inventory.quantity_unit = this.data.quantity.unit;
      this.inventory.price_per_unit = this.data.price_per_unit.amount;
      this.inventory.expiration_date = moment(this.data.expiration_date).format('YYYY-MM-DD');
      this.inventory.notes = this.data.notes;
    } else {
      this.inventory.quantity_unit = this.options.quantityUnits[0].key;
      this.inventory.chemical_type = this.options.chemicalTypes[0].key;
    }
  },
  methods: {
    ...mapActions([
      'submitMaterial',
      'openPicker',
    ]),
    submit() {
      this.inventory.expiration_date = moment(this.inventory.expiration_date).format('YYYY-MM-DD');
      this.inventory.type = 'agrochemical';
      this.submitMaterial(this.inventory)
        .then(() => this.$emit('closeModal'))
        .catch(() => this.$toasted.error('Error in material submission'));
    },
    closeModal() {
      this.$emit('closeModal');
    },
    openPicker() {
      this.$refs.openCal.showCalendar();
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
