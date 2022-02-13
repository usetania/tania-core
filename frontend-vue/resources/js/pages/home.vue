<template lang="pug">
.container-fluid.bottom-space
  .row
    .col
      h3.title-page
        translate Dashboard

  .row
    // At A Glance
    .col-xs-12.col-sm-12.col-md-6
      b-card(
        :title="$gettext('At A Glance')"
        :footer="$gettext('You are using Tania 1.7 right now')"
        class="card-ui"
      )
        .row
          .col-xs-12.col-sm-12
            i.fa.fa-grip-horizontal
            router-link(:to="{ name: 'FarmAreas' }")
              | {{ areas.length }}
              |
              translate Areas

          .col-xs-12.col-sm-12
            i.fa.fa-leaf
            router-link(:to="{ name: 'FarmCrops' }")
              | {{ cropInformation.total_plant_variety }}
              |
              translate Varieties

          .col-xs-12.col-sm-12
            i.fa.fa-clipboard
            router-link(:to="{ name: 'Task' }")
              | {{ tasksLength }}
              |
              translate Tasks

    // CROPS STATUS
    .col-xs-12.col-sm-12.col-md-6
      b-card(
        :title="$gettext('What is On Production')"
        class="card-ui"
      )
        router-link(:to="{ name: 'FarmCrops' }" class="see-all")
          translate See all Crops

        table.table
          thead
            tr
              th
                translate Varieties
              th
                translate Qty
          tbody
            tr(v-if="crops.length == 0")
              td(colspan="2")
                translate No crops on production.
                |
                |
                router-link(:to="{ name: 'FarmCrops' }") Add your first crops here.
            tr(v-for="crop in crops")
              td
                router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }")
                  | {{ crop.inventory.name }}
              td {{ crop.container.quantity }}
        Pagination(:pages="cropPages" @reload="getCrops")

  // TASK LIST
  .row
    .col-xs-12.col-sm-12.col-md-12
      b-card(
        :title="$gettext('Tasks')"
        class="card-ui"
      )
        router-link(:to="{ name: 'Task' }" class="see-all")
          translate See all Tasks
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

<style lang="scss" scoped>
h3.title-page {
  margin: 20px 0 30px 0;
}

.card-ui {
  margin-bottom: 20px;

  i {
    width: 30px;
  }

  .see-all {
    display: block;
    margin-bottom: 15px;
  }
}

.bottom-space {
  padding-bottom: 60px;
}
</style>
