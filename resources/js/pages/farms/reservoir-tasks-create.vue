<template lang="pug">
  .reservoir-tasks-create
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
              label(for="due_date") Due Date
              .input-group
                datepicker#due_date(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('due_date') }" v-model="task.due_date" name="due_date" input-class="form-control" ref="openCal")
                span.input-group-btn
                  button.btn.btn-primary(type="button" v-on:click="openPicker")
                    i.glyphicon.glyphicon-calendar
                span.help-block.text-danger(v-show="errors.has('due_date')") {{ errors.first('due_date') }}
          .col-xs-6
            .form-group
              label Assign to
              input.form-control
        .row
          .col-xs-6
            .form-group
              label(for="priority") Is this task urgent?
              .radio
                label.i-checks.i-checks-sm
                  input#priority(type="radio" name="priority" value="yes" checked="checked" v-model="task.priority" v-validate="'required'")
                  i
                  | Yes
              .radio
                label.i-checks.i-checks-sm
                  input#priority(type="radio" name="priority" value="no" v-model="task.priority" v-validate="'required'")
                  i
                  | No
        .form-group
          button.btn.btn-addon.btn-success.pull-right(type="submit") Save
          button.btn.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')") Cancel
</template>

<script>
import { StubTask } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
import Datepicker from 'vuejs-datepicker';

export default {
  name: "FarmReservoirTasksCreate",
  data () {
    return {
      task: Object.assign({}, StubTask),
    }
  },
  components: {
      Datepicker
  },
  methods: {
    ...mapActions([
      'createReservoirTask',
      'openPicker',
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    openPicker () {
      this.$refs.openCal.showCalendar()
    },
    create () {
    },
  }
}
</script>


