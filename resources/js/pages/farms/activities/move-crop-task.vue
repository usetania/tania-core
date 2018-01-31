<template lang="pug">
  .move-crop-task.col
    .modal-header
      span.h4.font-bold Move Crops
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="type") Select area to move this crop
          select.form-control#area(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="task.area" name="area")
            option(value="") Please select area
            option(v-for="option in options" :value="option.key") {{ option.label }}
          span.help-block.text-danger(v-show="errors.has('type')") {{ errors.first('type') }}
          label(for="name") Reservoir Name
        .form-group
          input.form-control#name(type="text" v-validate="'required|alpha_num_space|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('name') }" v-model="reservoir.name" name="name")
          span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>


<script>
import { mapActions } from 'vuex'
export default {
  name: "MoveCropTask",
  computed : {
    ...mapGetters({
      areas: 'getAllAreas'
    })
  },
  mounted () {
    this.fetchAreas()
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
      'fetchAreas',
    },
  }
}
</script>
