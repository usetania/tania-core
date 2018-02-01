<template lang="pug">
  .fertilizer-crop-task
    .modal-header
      span.h4.font-bold Fertilize 
        span.identifier {{ crop.batch_id }}
      span.pull-right.text-muted(style="cursor: pointer;" @click="$parent.$emit('close')")
        i.fa.fa-close
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
                  input#priority(type="radio" name="priority" value="yes" checked="checked" v-model="task.priority" v-validate="'required'")
                  i
                  | Yes
              .radio
                label.i-checks.i-checks-sm
                  input#priority(type="radio" name="priority" value="no" v-model="task.priority" v-validate="'required'")
                  i
                  | No
        .row
          .col-xs-6
            .form-group
              label(for="asset_id") 
                | Select area to do your task
              select.form-control#asset_id(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('asset_id') }" v-model="task.asset_id" name="asset_id")
                option(value="") Please select area
                option(v-for="area in areas" :value="area.uid") {{ area.name }}
              span.help-block.text-danger(v-show="errors.has('asset_id')") {{ errors.first('asset_id') }}
          .col-xs-6
            .form-group
              label(for="category") 
                | Task Category
              select.form-control#category(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('category') }" v-model="task.category" name="category")
                option(selected="selected" value="Nutrient") Nutrient
              span.help-block.text-danger(v-show="errors.has('category')") {{ errors.first('category') }}
        .form-group
          label(for="fertilizer") 
            | Select fertilizer you are going to use
          select.form-control#fertilizer(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('fertilizer') }" v-model="task.fertilizer" name="fertilizer")
            option(value="") Please select fertilizer
            option(v-for="fertilizer in fertilizers" :value="fertilizer.uid") {{ fertilizer.name }}
          span.help-block.text-danger(v-show="errors.has('fertilizer')") {{ errors.first('fertilizer') }}
        .form-group
          label(for="title") Title
          input.form-control#title(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('title') }" v-model="task.title" name="title")
          span.help-block.text-danger(v-show="errors.has('title')") {{ errors.first('title') }}
        .form-group
          label(for="description") Description
          textarea.form-control#description(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('description') }" v-model="task.description" name="description" rows="3")
          span.help-block.text-danger(v-show="errors.has('description')") {{ errors.first('description') }}
        .form-group
          .text-center.m-t
            button.btn.btn-primary(type="submit")
              i.fa.fa-check
              |  OK
</template>


<script>
import { mapGetters, mapActions } from 'vuex'
import { StubTask } from '@/stores/stubs'
import Datepicker from 'vuejs-datepicker';
export default {
  name: "FertilizerCropTask",
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
      fertilizers: [],
      task: Object.assign({}, StubTask),
    }
  },
  props: ['crop'],
  mounted () {
    this.fetchAreas()
  },
  created () {
  },
  methods: {
    ...mapActions([
      'fetchAreas',
      'openPicker',
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
    openPicker () {
      this.$refs.openCal.showCalendar()
    },
  }
}
</script>
