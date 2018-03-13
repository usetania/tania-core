<template lang="pug">
  .task-list(v-if="loading === false")
    table.table.m-b-none(v-if="domain == 'AREA' || domain == 'RESERVOIR'")
      thead
        tr
          th 
          th Items 
          th Category
          th(v-if="domain != 'AREA' && domain != 'RESERVOIR'")
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
              MoreDetail(:data="task" :description="task.description")
              small.text-muted(v-if="task.due_date") Due date: {{ task.due_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                TaskLabel(:type="'PRIORITY'" :task="task")
                span.text-danger(v-if="task.is_due == true") Overdue!
                span.text-success(v-if="isToday(task.due_date)") Today
          td
            TaskLabel(:type="'CATEGORY'" :task="task")
          td(v-if="domain != 'AREA' && domain != 'RESERVOIR'")
            a.h3(style="cursor: pointer;" @click="openModal(task)")
              i.fas.fa-edit
    div(v-else)
      p(v-if="tasks.length == 0") No Task Created
      li.list-group-item.clearfix(v-for="task in tasks")
        .row
          .col-sm-1
            .checkbox
              label.i-checks
                input(type="checkbox" v-on:change="setTaskStatus(task.uid, task.status)" :checked="isCompleted(task.status)")
                i
          .col-sm-8
            span.h4.text-dark(v-if="task.category == 'PESTCONTROL' || task.category == 'NUTRIENT'")
              | Apply 
              u(v-if="task.domain_details.material") {{ task.domain_details.material.material_name }}
              |  to 
              span.identifier-sm(v-if="task.domain_details.crop") {{ task.domain_details.crop.crop_batch_id }}
              |  on 
              span.areatag-sm(v-if="task.domain_details.area") {{ task.domain_details.area.area_name }}
            span.h4.text-dark(v-else-if="task.category == 'AREA'")
              span.areatag-sm(v-if="task.domain_details.area") {{ task.domain_details.area.area_name }}
              i.fas.fa-long-arrow-alt-right
              |  {{ task.title }}
            span.h4.text-dark(v-else-if="task.category == 'RESERVOIR'")
              u(v-if="task.domain_details.reservoir") {{ task.domain_details.reservoir.reservoir_name }}
              i.fas.fa-long-arrow-alt-right
              |  {{ task.title }}
            span.h4.text-dark(v-else-if="task.category == 'CROP'")
              span.identifier-sm(v-if="task.domain_details.crop") {{ task.domain_details.crop.crop_batch_id }}
              |  on 
              span.areatag-sm(v-if="task.domain_details.area") {{ task.domain_details.area.area_name }}
              i.fas.fa-long-arrow-alt-right
              |  {{ task.title }}
            span.h4.text-dark(v-else-if="task.category == 'SAFETY' || task.category == 'SANITATION'")
              span.areatag-sm(v-if="task.domain_details.area") {{ task.domain_details.area.area_name }}
              i.fas.fa-long-arrow-alt-right
              |  {{ task.title }}
            span.h4.text-dark(v-else) {{ task.title }}
            MoreDetail(:data="task" :description="task.description")
            div
              small.text-muted Due date: {{ task.due_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
              .status.status-urgent(v-if="task.priority == 'URGENT'") URGENT
              span.text-danger(v-if="task.is_due == true") Overdue!
          .col-sm-2
            TaskLabel(:type="'CATEGORY'" :task="task")
          .col-sm-1.text-right
            a.h3(v-if="!isCompleted(task.status)" style="cursor: pointer;" @click="openModal(task)")
              i.fas.fa-edit
</template>

<script>
import { AddClicked } from '@/stores/helpers/farms/crop'
import { mapActions, mapGetters } from 'vuex'
import moment from 'moment-timezone'
import MoreDetail from '@/components/more-detail'
import TaskLabel from './task-label'
export default {
  name: 'TasksList',
  components: {
    MoreDetail,
    TaskLabel
  },
  created () {
    this.getTasks() 
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
      'getTasksByCategoryAndPriorityAndStatus',
      'getTasks',
      'fetchTasks',
      'isToday',
      'isCompleted',
      'setTaskCompleted',
      'setTaskDue',
      'setTaskStatus',
    ]),
    getTasks () {
      if (this.domain) {
        this.getTasksByDomainAndAssetId({ domain: this.domain, assetId: this.asset_id })
          .then(({ data }) => {
            this.loading = false
            this.tasks = AddClicked(data)
          })
          .catch(error => console.log(error))
      } else if (this.category != '' || this.priority != '' || this.status != '') {
        let status = (this.status == 'INCOMPLETE') ? '' : this.status
        this.getTasksByCategoryAndPriorityAndStatus({ category: this.category, priority: this.priority, status: status })
          .then(({ data }) => {
            this.loading = false
            this.tasks = AddClicked(data)
          })
          .catch(error => console.log(error))
      } else {
        this.fetchTasks()
          .then(({ data }) => {
            this.loading = false
            this.tasks = AddClicked(data)
          })
          .catch(error => console.log(error))
      }
    },
    isCompleted (status) {
      return (status == "COMPLETED") ? true : false
    },
    isToday (date) {
      return moment(date).tz('Asia/Jakarta').isSame(moment(), 'day')
    },
    openModal (data) {
      this.data = data
      this.$emit('openModal', this.data)
    },
    setTaskStatus (task_id, status) {
      this.setTaskCompleted(task_id)
        .then(() => { this.getTasks() })
        .catch(({ data }) => this.message = data)
    },
  },
  mounted(){
    this.$watch('reload', reload => {
      this.getTasks()
    }, {})
    this.$watch('category', category => {
      this.getTasks()
    }, {})
    this.$watch('priority', priority => {
      this.getTasks()
    }, {})
    this.$watch('status', priority => {
      this.getTasks()
    }, {})
  },
  props: ['asset_id', 'category', 'domain', 'priority', 'reload', 'status'],
}
</script>
