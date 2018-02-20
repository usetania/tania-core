<template lang="pug">
  .dashboard.col
    .wrapper-md
      .m-t.font-thin.h3.text-black Dashboard
    .wrapper-md
      .row
        // At A Glance
        .col-sm-4
          .panel.ataglance.no-border
            .panel-heading
              .m-sm
                span.h4.text-lt At A Glance
            .panel-body
              ul
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-th-large
                    router-link(:to="{ name: 'FarmAreas' }") {{ areas.length }} Areas
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-leaf
                    router-link(:to="{ name: 'FarmCrops' }") {{ crops.length }} Varieties
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-clipboard
                    router-link(:to="{ name: 'Task' }") {{ tasks.length }} Tasks
            .panel-footer.bg-light.lter.wrapper.no-border
              small.text-muted
                | You are using Tania 1.5 right now.
                // There's a new version recently released. <a href="#">Take a look!</a>
        .col-sm-8
          // CROPS STATUS
          .panel.no-border
            .panel-heading
              .m-b.m-t
                span.pull-right.text-primary
                  router-link(:to="{ name: 'FarmCrops' }")
                    | See all Crops 
                    i.fa.fa-angle-double-right
                span.h4.text-lt Crops
            table.table.m-b-none
              thead
                tr
                  th Crop Variety
                  th Batches Qty
              tbody
                tr(v-for="crop in crops")
                  td: router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") {{ crop.inventory.name }}
                  td {{ crop.container.quantity }}
      .row
        // TASK LIST
        .panel.no-border
          .panel-heading
            .m-b.m-t
              span.pull-right.text-primary
                router-link(:to="{ name: 'Task' }")
                  | See all Tasks 
                  i.fa.fa-angle-double-right
              span.h4.text-lt Tasks
          table.table.m-b-none
            thead
              tr
                th Description
                th Category
                th Status
            tbody
              tr(v-if="tasks.length == 0")
                td(colspan="3") No Task Created
              tr(v-for="task in tasks")
                td
                  a(href="#")
                    div {{ task.title }}
                    small.text-muted Due date: {{ task.due_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                    .text-danger(v-if="task.is_due == true") Overdue!
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
                td
                  span.status.status-urgent(v-if="task.priority == 'URGENT'") URGENT
                  span.status.status-normal(v-if="task.priority == 'NORMAL'") NORMAL
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'Home',
  computed : {
    ...mapGetters({
      areas: 'getAllAreas',
      crops: 'getAllCrops',
      tasks: 'getAllTasks',
    })
  },
  mounted () {
    this.fetchAreas()
    this.fetchCrops()
    this.fetchTasks()
  },
  methods: {
    ...mapActions([
      'fetchAreas',
      'fetchCrops',
      'fetchTasks',
    ])
  }
}
</script>

