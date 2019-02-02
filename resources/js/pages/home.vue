<template lang="pug">
  .dashboard.col
    .wrapper-md
      .m-t.font-thin.h3.text-black Dashboard
    .wrapper-md
      .row
        // At A Glance
        .col-sm-4
          .card.ataglance
            .card-header
              .m-sm
                span.h4.text-lt At A Glance
            .card-body
              .row
                .col-md-6
                  h5.text-lt
                    i
                      font-awesome-icon(icon="dice-four")
                    |
                    |
                    router-link(:to="{ name: 'FarmAreas' }") {{ areas.length }} Areas

                .col-md-6
                  h5.text-lt
                    i
                      font-awesome-icon(icon="leaf")
                    |
                    |
                    router-link(:to="{ name: 'FarmCrops' }") {{ cropInformation.total_plant_variety }} Varieties
      
                .col-md-6
                  h5.text-lt
                    i
                      font-awesome-icon(icon="clipboard")
                    |
                    |
                    router-link(:to="{ name: 'Task' }") {{ tasksLength }} Tasks
            .card-footer.bg-light.lter.wrapper
              small.text-muted
                | You are using Tania 2.0.
                // There's a new version recently released. <a href="#">Take a look!</a>
        .col-sm-8
          // CROPS STATUS
          .card
            .card-header
              .m-sm
                span.float-right.text-primary
                  router-link(:to="{ name: 'FarmCrops' }")
                    | See all Crops 
                    i
                      font-awesome-icon(icon="angle-double-right")
                span.h4.text-lt Crops
            .card-body
              table.table
                thead
                  tr
                    th Crop Variety
                    th Batches Qty
                tbody
                  tr(v-if="crops.length == 0")
                    td(colspan="2") No Task Created
                  tr(v-for="crop in crops")
                    td: router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") {{ crop.inventory.name }}
                    td {{ crop.container.quantity }}
              Pagination(:pages="cropPages" @reload="getCrops")
      .row
        // TASK LIST
        .col-sm-12
          .card
            .card-header
              .m-sm
                span.float-right.text-primary
                  router-link(:to="{ name: 'Task' }")
                    | See all Tasks
                    i
                      font-awesome-icon(icon="angle-double-right")
                span.h4.text-lt Tasks
            .card-body
              TasksList(:domain="'HOME'")
              Pagination(:pages="taskPages" @reload="getTasks")
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import TaskLabel from './farms/tasks/task-label.vue'
import TasksList from './farms/tasks/task-list.vue'
import Pagination from '@/components/pagination.vue'
export default {
  name: 'Home',
  components: {
    Pagination,
    TaskLabel,
    TasksList
  },
  computed : {
    ...mapGetters({
      areas: 'getAllAreas',
      crops: 'getAllCrops',
      cropInformation: 'getInformation',
      cropPages: 'getCropsNumberOfPages',
      tasks: 'getTasks',
      taskPages: 'getTasksNumberOfPages',
      tasksLength: 'getNumberOfTasks',
    })
  },
  methods: {
    ...mapActions([
      'fetchAreas',
      'fetchCrops',
      'fetchTasks',
      'getInformation',
    ]),
    getCrops () {
      let pageId = 1
      if (typeof this.$route.query.page != "undefined") {
        pageId = parseInt(this.$route.query.page)
      }
      this.fetchCrops({ pageId : pageId, status : 'ACTIVE' })
    },
    getTasks () {
      let pageId = 1
      if (typeof this.$route.query.page != "undefined") {
        pageId = parseInt(this.$route.query.page)
      }
      this.fetchTasks({ pageId : pageId })
    },
  },
  mounted () {
    this.fetchAreas()
    this.getCrops()
    this.getTasks()
    this.getInformation()
  },
}
</script>

