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
                            .h3.text-lt Incomplete Tasks
                          .col-sm-4.text-right
                            a.btn.btn-sm.btn-primary.btn-addon(style="cursor: pointer;" @click="showModal = true")
                              i.fas.fa-plus
                              | Add Task
                      .panel-body.bg-white-only
                        .row
                        TasksList
    .col.w-lg.bg-light.lter.b-l.bg-auto.no-border-xs
      .wrapper
        .form-group
          div
            label.control-label Category
          div
            select.form-control
              option Area
              option Crop
              option Finance
              option General
              option Inventory
              option Nutrient
              option Pest Control
              option Reservoir
              option Safety
              option Sanitation
        .form-group
          div
            label.control-label Priority
          div
            select.form-control
              option Urgent
              option Normal
      .wrapper
        ul.list-group.no-bg.no-borders.pull-in
          li.list-group-item
            .wrapper-xs
              .h4.text-info
                a(href="task-completed.html") Completed
          li.list-group-item.active
            .wrapper-xs
              .h4.text-muted
                a(href="task.html") Incomplete
          li.list-group-item.m-l
            .wrapper-xs
              .h5.text-danger
                a(href="") Overdue
          li.list-group-item.m-l
            .wrapper-xs
              .h5.text-success
                a(href="") Today
          li.list-group-item.m-l
            .wrapper-xs
              .h5.text-lt
                a(href="") This Week
          li.list-group-item.m-l
            .wrapper-xs
              .h5.text-lt
                a(href="") This Month
      .wrapper-md
      .wrapper-md
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal'
import TasksList from '@/pages/farms/tasks/task-list.vue'
import TaskCreate from '@/pages/farms/tasks/task-create.vue'

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
      showModal: false,
    }
  },
  methods: {
    ...mapActions([
    ]),
    closeModal () {
      this.showModal = false
    },
  }
}
</script>
