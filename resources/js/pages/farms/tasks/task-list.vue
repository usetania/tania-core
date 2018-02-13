<template lang="pug">
  .task-list.col-sm-6.col-xs-12(v-if="loading === false")
    .panel
      .panel-heading
        span.h4.text-lt Tasks
      table.table.m-b-none
        thead
          tr
            th 
            th Items
            th Category
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
                div {{ task.title }}
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
