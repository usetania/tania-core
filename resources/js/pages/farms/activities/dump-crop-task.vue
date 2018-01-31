<template lang="pug">
  .dump-crop-task(v-if="loading === false")
    .modal-header
      span.h4.font-bold Move Crops
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="type") Select area to move this crop
          select.form-control#dest_area_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="task.dest_area_id" name="dest_area_id")
            option(value="") Please select area
            option(v-for="area in areas" :value="area.uid") {{ area.name }}
          span.help-block.text-danger(v-show="errors.has('dest_area_id')") {{ errors.first('dest_area_id') }}
        .form-group
          label(for="quantity")
            | How many of 
            span.identifier-sm {{ crop.batch_id }}
            |  do you want to move?
          input.form-control#quantity(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('quantity') }" v-model="task.quantity" name="quantity")
          span.help-block.text-danger(v-show="errors.has('quantity')") {{ errors.first('quantity') }}
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>


<script>
import { mapGetters, mapActions } from 'vuex'
import { StubTask, StubCrop } from '@/stores/stubs'
export default {
  name: "DumpCropTask",
  computed : {
    ...mapGetters({
      areas: 'getAllAreas',
    })
  },
  data () {
    return {
      loading: true,
      crop: Object.assign({}, StubCrop),
      task: Object.assign({}, StubTask),
    }
  },
  mounted () {
    this.fetchAreas()
  },
  created () {
    this.getCropByUid(this.$route.params.id)
      .then(({ data }) =>  {
        this.loading = false
        this.crop = data
      })
      .catch(error => console.log(error))
  },
  methods: {
    ...mapActions([
      'fetchAreas',
      'getCropByUid',
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
    },
  }
}
</script>
