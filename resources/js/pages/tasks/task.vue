<template lang="pug">
  .col
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
                          a.btn.btn-sm.btn-primary.btn-addon(href="task-add.html")
                            i.fas.fa-plus
                            | Add Task
                    .panel-body.bg-white-only
                      .row
                      TasksList
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
