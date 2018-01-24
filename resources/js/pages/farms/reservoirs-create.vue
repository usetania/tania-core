<template lang="pug">
  .reservoirs-create.col
    .modal-header
      span.h4.font-bold Add New Reservoir
    .modal-body
      small.text-muted Reservoir is a water source for your farm. It can be a direct line from well, or water tank that has volume/capacity.
      error-message(:message="message.error_message")
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="name") Reservoir Name
          input.form-control#name(type="text" v-validate="'required|alpha_num|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('reservoir.name') }" v-model="reservoir.name" name="reservoir.name")
          span.help-block.text-danger(v-show="errors.has('reservoir.name')") {{ errors.first('reservoir.name') }}
        .form-group
          label(for="type") Source
          select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('reservoir.type') }" v-model="reservoir.type" name="reservoir.type" @change="typeChanged($event.target.value)")
            option(value="") Please select source
            option(v-for="option in options" :value="option.key") {{ option.label }}
          span.help-block.text-danger(v-show="errors.has('reservoir.type')") {{ errors.first('reservoir.type') }}
        .form-group(v-if="reservoir.type == 'bucket'")
          input.form-control(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('reservoir.capacity') }" v-model="reservoir.capacity" placeholder="Capacity (litre)" name="reservoir.capacity")
          span.help-block.text-danger(v-show="errors.has('reservoir.capacity')") {{ errors.first('reservoir.capacity') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>


<script>
import { ReservoirTypes } from '@/stores/helpers/farms/reservoir'
import { StubReservoir, StubMessage } from '@/stores/stubs'
import { mapActions } from 'vuex'
export default {
  name: "FarmReservoirCreate",
  data () {
    return {
      message: Object.assign({}, StubMessage),
      reservoir: Object.assign({}, StubReservoir),
      options: Array.from(ReservoirTypes)
    }
  },
  methods: {
    ...mapActions([
      'createReservoir'
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
      this.createReservoir(this.reservoir)
        .then(({ data }) => this.$router.push({ name: 'FarmReservoirs'}))
        .catch(({ data }) => this.message = data)
    },
    typeChanged (type) {
      if (type === 'bucket') {
        this.reservoir.capacity = ''
      } else {
        this.reservoir.capacity = 0
      }
    },
  }
}
</script>
