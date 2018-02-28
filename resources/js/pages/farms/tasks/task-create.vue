<template lang="pug">
  .tasks-create
    .modal-header
      h4.font-bold(v-if="asset != 'General'")
        | {{ asset }}: Add New Task on 
        span.areatag {{ data.name }}
      h4.font-bold(v-else) Add New Task
    .modal-body
      form(@submit.prevent="validateBeforeSubmit")
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
              label(for="priority") Is this task urgent?
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
        .form-group
          label(for="category") 
            | Task Category
          select.form-control#category(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('category') }" v-model="task.category" name="category")
            option(value="") Please select category
            option(v-if="asset == 'Area'" value="AREA") Area
            option(v-if="asset == 'Reservoir'" value="RESERVOIR") Reservoir
            option(v-if="asset == 'General'" value="GENERAL") General
            option(v-for="category in options.taskCategories" :value="category.key") {{ category.label }}
          span.help-block.text-danger(v-show="errors.has('category')") {{ errors.first('category') }}
        .form-group
          label(for="title") Title
          input.form-control#title(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('title') }" v-model="task.title" name="title")
          span.help-block.text-danger(v-show="errors.has('title')") {{ errors.first('title') }}
        .form-group
          label(for="description") Description
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
  name: "FarmTasksCreate",
  components: {
      Datepicker
  },
  data () {
    return {
      task: Object.assign({}, StubTask),
      options: {
        taskCategories: Array.from(TaskDomainCategories),
      }
    }
  },
  methods: {
    ...mapActions([
      'openPicker',
      'createTask',
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
      console.log(this.data.uid)
      if (typeof this.data.uid != "undefined") {
        this.task.asset_id = this.data.uid
      }
      this.task.domain = this.asset.toUpperCase()
      this.createTask(this.task)
        .then(this.$parent.$emit('close'))
        .catch(({ data }) => this.message = data)
    },
  },
  mounted () {
    console.log(this.asset)
  },
  props: ['data', 'asset'],
}
</script>


