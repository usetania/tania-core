<template lang="pug">
  .upload-crop-task
    .modal-header
      h4.font-bold
        translate Watering
    .modal-body
      b-form(@submit.prevent="validateBeforeSubmit")
        .form-group
          label(for="type")
            translate Choose type of watering
          select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('type') }" v-model="task.type" name="type" @change="typeChanged($event.target.value)")
            option(value="ALL")
              translate All
            option(value="PARTIAL")
              translate Partial
          span.help-block.text-danger(v-show="errors.has('type')") {{ errors.first('type') }}
        .form-group
          label
            translate Which crop do you want to water?
          .checkbox(v-for="crop in crops")
            label.i-checks
              input(type="checkbox" name="selected crops" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('selected crops') }" v-model="task.crops" v-bind:value="crop.uid")
              i
              | {{ crop.inventory.name }}
              |
              |
              span.identifier-sm {{ crop.batch_id }}
          span.help-block.text-danger(v-show="errors.has('selected crops')") {{ errors.first('selected crops') }}
        .form-group
          button.btn.btn-success.float-right(type="submit")
            i.fa.fa-check
            translate SAVE
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')")
            i.fa.fa-close
            translate Cancel
</template>

<script>
import { mapActions } from 'vuex';
import moment from 'moment-timezone';
import { StubTask } from '../../../stores/stubs';

export default {
  name: 'WaterTaskModal',
  props: ['crops', 'area'],
  data() {
    return {
      task: Object.assign({}, StubTask),
    };
  },
  created() {
    this.task.type = 'PARTIAL';
    this.task.crops = [];
  },
  methods: {
    ...mapActions([
      'waterCrop',
    ]),
    create() {
      this.task.source_area_id = this.area.uid;
      this.task.watering_date = moment().format('YYYY-MM-DD HH:ss');
      for (let i = 0; i < this.task.crops.length; i += 1) {
        this.task.obj_uid = this.task.crops[i];
        this.waterCrop(this.task)
          .then(() => this.$parent.$emit('close'))
          .catch(() => this.$toasted.error('Error in water task submission'));
      }
    },
    typeChanged(type) {
      if (type === 'ALL') {
        for (let i = 0; i < this.crops.length; i += 1) {
          this.task.crops.push(this.crops[i].uid);
        }
      } else {
        this.task.crops = [];
      }
    },
    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.create();
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
i.fa.fa-check,
i.fa.fa-close {
  width: 30px;
}
</style>
