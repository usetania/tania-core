<template lang="pug">
  .dashboard.col
    .wrapper-md
      .m-t.font-thin.h3.text-black
        translate Dashboard
    .wrapper-md
      .row
        // At A Glance
        .col-sm-4
          .panel.ataglance.no-border
            .panel-heading
              .m-sm
                span.h4.text-lt
                  translate At A Glance
            .panel-body
              ul
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-th-large
                    router-link(:to="{ name: 'FarmAreas' }")
                      | {{ areas.length }}
                      |
                      translate Areas
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-leaf
                    router-link(:to="{ name: 'FarmCrops' }")
                      | {{ cropInformation.total_plant_variety }}
                      |
                      translate Varieties
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-clipboard
                    router-link(:to="{ name: 'Task' }")
                      | {{ tasksLength }}
                      |
                      translate Tasks
            .panel-footer.bg-light.lter.wrapper.no-border
              small.text-muted
                translate You are using Tania 1.7 right now.
                // There's a new version recently released. <a href="#">Take a look!</a>
        .col-sm-8
          // CROPS STATUS
          .panel.no-border
            .panel-heading
              .m-b.m-t
                span.pull-right.text-primary
                  router-link(:to="{ name: 'FarmCrops' }")
                    translate See all Crops
                span.h4.text-lt
                  translate What's On Production
            table.table.m-b-none
              thead
                tr
                  th
                    translate Crop Variety
                  th
                    translate Batches Qty
              tbody
                tr(v-if="crops.length == 0")
                  td(colspan="2")
                    translate No crops on production.
                    | &nbsp;
                    a(href="/#/crops") Add your first crops here.
                tr(v-for="crop in crops")
                  td
                    router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }")
                      | {{ crop.inventory.name }}
                  td {{ crop.container.quantity }}
            Pagination(:pages="cropPages" @reload="getCrops")
      .row
        // TASK LIST
        .col-sm-12
          .panel.no-border
            .panel-heading
              .m-b.m-t
                span.pull-right.text-primary
                  router-link(:to="{ name: 'Task' }")
                    translate See all Tasks
                span.h4.text-lt
                  translate Tasks
            TasksList(:domain="'HOME'")
            Pagination(:pages="taskPages" @reload="getTasks")
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import TaskLabel from './farms/tasks/task-label.vue';
import TasksList from './farms/tasks/task-list.vue';
import Pagination from '../components/pagination.vue';

export default {
  name: 'Home',
  components: {
    Pagination,
    TaskLabel,
    TasksList,
  },
  computed: {
    ...mapGetters({
      areas: 'getAllAreas',
      crops: 'getAllCrops',
      cropInformation: 'getInformation',
      cropPages: 'getCropsNumberOfPages',
      tasks: 'getTasks',
      taskPages: 'getTasksNumberOfPages',
      tasksLength: 'getNumberOfTasks',
    }),
  },
  mounted() {
    this.fetchAreas();
    this.getCrops();
    this.getTasks();
    this.getInformation();
  },
  methods: {
    ...mapActions([
      'fetchAreas',
      'fetchCrops',
      'fetchTasks',
      'getInformation',
    ]),
    getCrops() {
      let pageId = 1;
      if (typeof this.$route.query.page !== 'undefined') {
        pageId = parseInt(this.$route.query.page, 10);
      }
      this.fetchCrops({ pageId, status: 'ACTIVE' });
    },
    getTasks() {
      let pageId = 1;
      if (typeof this.$route.query.page !== 'undefined') {
        pageId = parseInt(this.$route.query.page, 10);
      }
      this.fetchTasks({ pageId });
    },
  },
};
</script>
