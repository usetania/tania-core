<template lang="pug">
  .hbox
    .col
      modal(v-if="showModal" @close="showModal = false")
        TaskCreate(:asset="'General'" :data="data")
      .vbox
        .row-row
          .cell
            .cell-inner
              .wrapper-md
                .row
                  .col-sm-12
                    .panel
                      .panel-heading.wrapper-md
                        .row
                          .col-sm-8
                            .h3.text-lt Tasks
                          .col-sm-4.text-right
                            a.btn.btn-sm.btn-primary.btn-addon(style="cursor: pointer;" @click="showModal = true")
                              i.fas.fa-plus
                              | Add Task
                      .panel-body.bg-white-only
                        .row
                        TasksList(:category="selected_category" :priority="selected_priority" :status="status")
    .col.w-lg.bg-light.lter.b-l.bg-auto.no-border-xs
      .wrapper
        .form-group
          div
            label.control-label Category
          div
            select.form-control(:class="{'input': true, 'text-danger': errors.has('category') }" name="category" @change="categoryChange($event.target.value)" v-model="selected_category")
              option(value="") All
              option(value="AREA") Area
              option(value="RESERVOIR") Reservoir
              option(value="CROP") Crop
              option(value="GENERAL") General
              option(v-for="category in options.taskCategories" :value="category.key") {{ category.label }}
        .form-group
          div
            label.control-label Priority
          div
            select.form-control(@change="priorityChange($event.target.value)" v-model="selected_priority")
              option(value="") All
              option(value="URGENT") Urgent
              option(value="NORMAL") Normal
      .wrapper
        ul.list-group.no-bg.no-borders.pull-in
          li.list-group-item(v-bind:class="{ active: isActive('COMPLETED') }")
            .wrapper-xs
              .h4.text-info
                a(style="cursor: pointer;" @click="statusSelected('COMPLETED')") Completed
          li.list-group-item(v-bind:class="{ active: isActive('INCOMPLETE') }")
            .wrapper-xs
              .h4.text-muted
                a(style="cursor: pointer;" @click="statusSelected('INCOMPLETE')") Incomplete
          li.list-group-item(v-bind:class="{ active: isActive('OVERDUE') }")
            .wrapper-xs
              .h5.text-danger
                a(style="cursor: pointer;" @click="statusSelected('OVERDUE')") Overdue
          li.list-group-item(v-bind:class="{ active: isActive('TODAY') }")
            .wrapper-xs
              .h5.text-success
                a(style="cursor: pointer;" @click="statusSelected('TODAY')") Today
          li.list-group-item(v-bind:class="{ active: isActive('THISWEEK') }")
            .wrapper-xs
              .h5.text-lt
                a(style="cursor: pointer;" @click="statusSelected('THISWEEK')") This Week
          li.list-group-item(v-bind:class="{ active: isActive('THISMONTH') }")
            .wrapper-xs
              .h5.text-lt
                a(style="cursor: pointer;" @click="statusSelected('THISMONTH')") This Month
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal'
import TasksList from '@/pages/farms/tasks/task-list.vue'
import TaskCreate from '@/pages/farms/tasks/task-create.vue'
import { TaskDomainCategories } from '@/stores/helpers/farms/task'

export default {
  name: 'Tasks',
  components: {
    Modal,
    TaskCreate,
    TasksList,
  },
  computed: {
    ...mapGetters({
    })
  },
  created () {
  },
  data () {
    return {
      data: {},
      options: {
        taskCategories: Array.from(TaskDomainCategories),
      },
      selected_category: '',
      selected_priority: '',
      showModal: false,
      status: 'INCOMPLETE',
    }
  },
  methods: {
    ...mapActions([
    ]),
    closeModal () {
      this.showModal = false
    },
    categoryChange (type) {
      this.selected_category = type
    },
    priorityChange (type) {
      this.selected_priority = type
    },
    statusSelected (status) {
      this.status = status
    },
    isActive (status) {
      return this.status == status
    }
  }
}
</script>
