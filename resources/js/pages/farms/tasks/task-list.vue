<template lang="pug">
  table.table.m-b-none(v-if="loading === false")
    thead
      tr
        th 
        th Items
        th Category
        th(v-if="domain == 'CROP'")
    tbody
      tr(v-if="tasks.length == 0")
        td(colspan="3") No Task Created
      tr(v-for="task in tasks")
        td
          .checkbox
            label.i-checks
              input(type="checkbox" v-on:change="setTaskStatus(task.uid, task.status)" :checked="isCompleted(task.status)")
              i
        td
          a(href="#")
            div {{ task.title }} : {{ task.description }}
            small.text-muted Due date: {{ task.due_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
            span.status.status-urgent(v-if="task.priority == 'URGENT'") URGENT
            span.status.status-normal(v-if="task.priority == 'NORMAL'") NORMAL
            span.text-danger(v-if="task.is_due == true") Overdue!
            span.text-success(v-if="isToday(task.due_date)") Today
        td
          span.label.label-pestcontrol(v-if="task.category == 'PESTCONTROL'") PEST CONTROL
          span.label.label-sanitation(v-if="task.category == 'SANITATION'") SANITATION
          span.label.label-area(v-if="task.category == 'AREA'") AREA
          span.label.label-safety(v-if="task.category == 'SAFETY'") SAFETY
          span.label.label-finance(v-if="task.category == 'FINANCE'") FINANCE
          span.label.label-crop(v-if="task.category == 'CROP'") CROP
          span.label.label-water(v-if="task.category == 'WATER'") WATER
          span.label.label-inventory(v-if="task.category == 'INVENTORY'") INVENTORY
          span.label.label-general(v-if="task.category == 'GENERAL'") GENERAL
          span.label.label-nutrient(v-if="task.category == 'NUTRIENT'") NUTRIENT
        td(v-if="domain == 'CROP'")
          a.h3(href="#")
            i.fas.fa-edit
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import moment from 'moment-timezone'
export default {
  name: 'TasksList',
  created () {
    this.fetchTasks() 
  },
  data () {
    return {
      loading: true,
      tasks: [],
    }
  },
  methods: {
    ...mapActions([
      'getTasksByDomainAndAssetId',
      'fetchTasks',
      'isToday',
      'isCompleted',
      'setTaskCompleted',
      'setTaskDue',
      'setTaskStatus',
    ]),
    fetchTasks () {
      this.getTasksByDomainAndAssetId({ domain: this.domain, assetId: this.asset_id })
        .then(({ data }) => {
          this.loading = false
          this.tasks = data
        })
        .catch(error => console.log(error))
    },
    isCompleted (status) {
      return (status == "COMPLETED") ? true : false
    },
    isToday (date) {
      return moment(date).tz('Asia/Jakarta').isSame(moment(), 'day')
    },
    setTaskStatus (task_id, status) {
      this.setTaskCompleted(task_id)
        .then(this.fetchTasks())
        .catch(({ data }) => this.message = data)
    },
  },
  mounted(){
    this.$watch('reload', reload => {
      this.fetchTasks()
    }, {})
  },
  props: ['domain', 'asset_id', 'reload'],
}
</script>
