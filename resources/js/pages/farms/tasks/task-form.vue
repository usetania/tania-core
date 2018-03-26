<template lang="pug">
  .tasks-form
    .modal-header
      h4.font-bold(v-if="task.uid") Update Task
      h4.font-bold(v-else-if="asset != 'General'")
        | {{ asset }}: Add New Task on 
        span.areatag {{ data.name }}
      h4.font-bold(v-else) Add New Task
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
        .row
          .col-xs-6
            .form-group
              label#label-due-date Due Date
              .input-group
                datepicker#due_date(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('due date') }" v-model="task.due_date" name="due date" input-class="form-control" ref="openCal")
                span.input-group-btn
                  button.btn.btn-primary(type="button" v-on:click="openPicker")
                    i.glyphicon.glyphicon-calendar
              span.help-block.text-danger(v-show="errors.has('due date')") {{ errors.first('due date') }}
          .col-xs-6
            .form-group
              label#label-priority Is this task urgent?
              .radio
                label.i-checks.i-checks-sm
                  input#priority(type="radio" name="priority" value="URGENT" checked="checked" v-model="task.priority" v-validate="'required'")
                  i
                  | Yes
              .radio
                label.i-checks.i-checks-sm
                  input#priority(type="radio" name="priority" value="NORMAL" v-model="task.priority" v-validate="'required'")
                  i
                  | No
              span.help-block.text-danger(v-show="errors.has('priority')") {{ errors.first('priority') }}
        .form-group
          label#label-category(for="category") 
            | Task Category
          select.form-control#category(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('category') }" v-model="task.category" name="category")
            option(value="") Please select category
            option(v-if="asset_name == 'AREA'" value="AREA") Area
            option(v-if="asset_name == 'RESERVOIR'" value="RESERVOIR") Reservoir
            option(v-if="asset_name == 'GENERAL'" value="GENERAL") General
            option(v-for="category in options.taskCategories" :value="category.key") {{ category.label }}
            option(v-if="asset_name == 'GENERAL'" value="INVENTORY") Inventory
          span.help-block.text-danger(v-show="errors.has('category')") {{ errors.first('category') }}
        .form-group
          label#label-title(for="title") Title
          input.form-control#title(type="text" v-validate="'required|max:100'" :class="{'input': true, 'text-danger': errors.has('title') }" v-model="task.title" name="title")
          span.help-block.text-danger(v-show="errors.has('title')") {{ errors.first('title') }}
        .form-group
          label#label-description(for="description") Description
          textarea.form-control#description(type="text" :class="{'input': true, 'text-danger': errors.has('description') }" v-model="task.description" name="description" rows="3")
          span.help-block.text-danger(v-show="errors.has('description')") {{ errors.first('description') }}
        .form-group
          button.btn.btn-addon.btn-primary.pull-right(type="submit")
            i.fas.fa-check
            |  OK
          button.btn.btn-addon.btn-default(style="cursor: pointer;" @click="$parent.$emit('close')")
            i.fas.fa-times
            |  Cancel
</template>

<script>
import { StubTask } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
import Datepicker from 'vuejs-datepicker';
import { TaskDomainCategories } from '@/stores/helpers/farms/task'

export default {
  name: "FarmTasksForm",
  components: {
      Datepicker
  },
  data () {
    return {
      task: Object.assign({}, StubTask),
      asset_name: '',
      options: {
        taskCategories: Array.from(TaskDomainCategories),
      }
    }
  },
  methods: {
    ...mapActions([
      'openPicker',
      'submitTask',
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.submit()
        }
      })
    },
    openPicker () {
      this.$refs.openCal.showCalendar()
    },
    submit () {
      if (typeof this.data.domain != "undefined") {
        this.task.domain = this.data.domain
      } else {
        this.task.asset_id = this.data.uid
        this.task.domain = this.asset.toUpperCase()
      }
      this.submitTask(this.task)
        .then(() => this.$parent.$emit('close'))
        .catch(() => this.$toasted.error('Error in task submission'))
    },
  },
  mounted () {
    if (typeof this.data.domain != "undefined") {
      this.task.uid = this.data.uid
      this.task.due_date = this.data.due_date
      this.task.priority = this.data.priority
      this.task.category = this.data.category
      this.task.title = this.data.title
      this.task.description = this.data.description
      this.asset_name = this.data.domain
    } else {
      this.asset_name = this.asset.toUpperCase()
    }
  },
  props: ['data', 'asset'],
}
</script>


