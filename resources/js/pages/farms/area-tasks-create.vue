<template lang="pug">
  .area-tasks-create
    .modal-header
      h4.text-lt Add Task
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .form-group
          small.text-muted.pull-right (max. 200 char)
          label(for="description") Task Description
          textarea.form-control#description(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('description') }" v-model="task.description" name="description" placeholder="What task do you want to work on?" rows="3")
          span.help-block.text-danger(v-show="errors.has('description')") {{ errors.first('description') }}
        .row
          .col-xs-6
            .form-group
              label Due Date
              .input-group
                input.form-control(type="text" datepicker-popup="" ng-model="dt" is-open="opened" datepicker-options="dateOptions" ng-required="true" close-text="Close")
                span.input-group-btn
                  button.btn.btn-primary(type="button" ng-click="open($event)")
                    i.glyphicon.glyphicon-calendar
          .col-xs-6
            .form-group
              label Is this task urgent?
              .radio
                label.i-checks.i-checks-sm
                  input(type="radio" name="urgenttask" value="yes" checked="checked")
                  i
                  |                           Yes
              .radio
                label.i-checks.i-checks-sm
                  input(type="radio" name="urgenttask" value="no")
                  i
                  |                           No
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>

<script>
import { AreaTypes, AreaLocations, AreaSizeUnits } from '@/stores/helpers/farms/area'
import { StubTask } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: "FarmAreaTasksCreate",
  data () {
    return {
      task: Object.assign({}, StubTask),
    }
  },
  methods: {
    ...mapActions([
      'createAreaTask',
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


