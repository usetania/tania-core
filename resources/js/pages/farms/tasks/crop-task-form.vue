<template lang="pug">
  .crop-task
    .modal-header
      span.h4.font-bold Crop: Add New Task for 
        span.identifier {{ batch_id }}
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
        .row
          .col-xs-6
            .form-group
              label(for="area_id") 
                | Select area to do your task
              select.form-control#area_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('area_id') }" v-model="task.area_id" name="area_id")
                option(value="") Please select area
                option(v-for="area in areas" :value="area.uid") {{ area.name }}
              span.help-block.text-danger(v-show="errors.has('area_id')") {{ errors.first('area_id') }}
          .col-xs-6
            .form-group
              label(for="category") 
                | Task Category
              select.form-control#category(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('category') }" v-model="task.category" name="category" @change="typeChanged($event.target.value)")
                option(value="") Please select category
                option(value="CROP") Crop
                option(value="NUTRIENT") Nutrient
                option(v-for="category in options.taskCategories" :value="category.key") {{ category.label }}
              span.help-block.text-danger(v-show="errors.has('category')") {{ errors.first('category') }}
        .form-group(v-if="isfertilizer")
          label(for="fertilizer") 
            | Select fertilizer you are going to use
          select.form-control#fertilizer(:class="{'input': true, 'text-danger': errors.has('fertilizer') }" v-model="task.fertilizer" name="fertilizer")
            option(value="") Please select fertilizer
            option(v-for="fertilizer in fertilizers" :value="fertilizer.uid") {{ fertilizer.name }}
          span.help-block.text-danger(v-show="errors.has('fertilizer')") {{ errors.first('fertilizer') }}
        .form-group(v-if="ispesticide")
          label(for="pesticide") 
            | Select pesticide you are going to use
          select.form-control#pesticide(:class="{'input': true, 'text-danger': errors.has('pesticide') }" v-model="task.pesticide" name="pesticide")
            option(value="") Please select pesticide
            option(v-for="pesticide in pesticides" :value="pesticide.uid") {{ pesticide.name }}
          span.help-block.text-danger(v-show="errors.has('pesticide')") {{ errors.first('pesticide') }}
        .form-group
          label(for="title") Title
          input.form-control#title(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('title') }" v-model="task.title" name="title")
          span.help-block.text-danger(v-show="errors.has('title')") {{ errors.first('title') }}
        .form-group
          label(for="description") Description
          textarea.form-control#description(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('description') }" v-model="task.description" name="description" rows="3")
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
import { mapGetters, mapActions } from 'vuex'
import { StubTask } from '@/stores/stubs'
import { TaskDomainCategories } from '@/stores/helpers/farms/task'
import Datepicker from 'vuejs-datepicker';
export default {
  name: "CropTask",
  components: {
      Datepicker
  },
  computed : {
    ...mapGetters({
      areas: 'getAllAreas',
    })
  },
  data () {
    return {
      batch_id: '',
      crop_id: '',
      isfertilizer: false,
      ispesticide: false,
      fertilizers: [],
      pesticides: [],
      task: Object.assign({}, StubTask),
      options: {
        taskCategories: Array.from(TaskDomainCategories),
      }
    }
  },
  methods: {
    ...mapActions([
      'fetchAreas',
      'getCropByUid',
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
    submit () {
      this.task.asset_id = this.crop_id
      this.task.domain = "CROP"
      this.submitTask(this.task)
        .then(this.$parent.$emit('close'))
        .catch(({ data }) => this.message = data)
    },
    openPicker () {
      this.$refs.openCal.showCalendar()
    },
    typeChanged (type) {
      this.isfertilizer = false
      this.ispesticide = false
      if (type == "CROP") {
        this.isfertilizer = true
      } else if (type == "PEST_CONTROL") {
        this.ispesticide = true
      }
    }
  },
  mounted () {
    this.fetchAreas()
    if (typeof this.data.uid != "undefined") {
      this.task.uid = this.data.uid
      this.task.due_date = this.data.due_date
      this.task.priority = this.data.priority
      this.task.category = this.data.category
      this.task.title = this.data.title
      this.task.description = this.data.description
      this.crop_id = this.data.asset_id
      this.getCropByUid(this.data.asset_id)
        .then(({ data }) =>  {
          this.batch_id = data.batch_id
          }).catch(error => console.log(error))
    } else {
      this.batch_id = this.crop.batch_id
      this.crop_id = this.crop.uid
    }
  },
  props: ['crop', 'data'],
}
</script>
