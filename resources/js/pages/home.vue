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
                    router-link(:to="{ name: 'FarmAreas' }") 4 Areas
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-leaf
                    router-link(:to="{ name: 'FarmCrops' }") 9 Varieties
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-clipboard
                    a(href="#") 7 Tasks
                li.h4
                  .col-md-6.col-xs-6
                    i.fa.fa-gamepad
                    a(href="#") 4 Devices
            .panel-footer.bg-light.lter.wrapper.no-border
              small.text-muted
                | You are using Tania 1.5 right now.
                // There's a new version recently released. <a href="#">Take a look!</a>
        // Current Status
        .col-sm-8
          .panel.curstat.no-border
            .panel-heading
              .m-sm
                span.pull-right
                  i.fa.fa-cog
                span.h4.text-lt Area Current Conditions
            .panel-body.m-b
              .row.m-b
                .col-sm-12.col-xs-12
                  select.form-control
                    option All Area
                    option Raise Bed A
                    option Raise Bed B
              .row.m-t
                // Sensor 1
                .col-sm-4.col-xs-12
                  .hbox.bg-white.m-t.wrapper.b.b-light
                    .h4.text-lt.m-b
                      a(href="#") Xiaomi Flora - AR142
                    .row.b-t.b-light
                      .col-xs-6.m-t
                        small.text-muted Temperature
                        .h3
                          | 28.2 
                          small ℃
                      .col-xs-6.m-t
                        small.text-muted Humidity
                        .h3 62%
                    .row.m-t
                      .col-xs-6
                        small.text-muted Soil Fertility
                        .h3 385
                      .col-xs-6
                        small.text-muted Soil Moisture
                        .h3 43%
                    .row.m-t
                      .col-xs-12
                        small.text-muted Light Intensity
                        .h3 38721
                // Sensor 2
                .col-sm-4.col-xs-12
                  .hbox.bg-white.m-t.wrapper.b.b-light
                    .h4.text-lt.m-b Xiaomi Flora - AR143
                    .row.b-t.b-light
                      .col-xs-6.m-t
                        small.text-muted Temperature
                        .h3
                          | 28.7 
                          small ℃
                      .col-xs-6.m-t
                        small.text-muted Humidity
                        .h3 67%
                    .row.m-t
                      .col-xs-6
                        small.text-muted Soil Fertility
                        .h3 427
                      .col-xs-6
                        small.text-muted Soil Moisture
                        .h3 47%
                    .row.m-t
                      .col-xs-12
                        small.text-muted Light Intensity
                        .h3 32859
                // Sensor 3
                .col-sm-4.col-xs-12
                  .hbox.bg-white.m-t.wrapper.b.b-light
                    .h4.text-lt.m-b Xiaomi Flora - AR147
                    .row.b-t.b-light
                      .col-xs-6.m-t
                        small.text-muted Temperature
                        .h3
                          | 28.7 
                          small ℃
                      .col-xs-6.m-t
                        small.text-muted Humidity
                        .h3 67%
                    .row.m-t
                      .col-xs-6
                        small.text-muted Soil Fertility
                        .h3 427
                      .col-xs-6
                        small.text-muted Soil Moisture
                        .h3 47%
                    .row.m-t
                      .col-xs-12
                        small.text-muted Light Intensity
                        .h3 32859
            .panel-footer.m-t.bg-light.lter.no-border.wrapper
              small.text-muted Last sync: 08.15 AM (2 hours 45 minutes ago)
      // Crops Status
      .row
        .col-sm-6
          // CROPS STATUS
          .panel.no-border(style="min-height: 200px;")
            .panel-heading
              .m-b.m-t
                span.pull-right.text-primary
                  a(href="crop.html")
                    | See all Crops 
                    i.fa.fa-angle-double-right
                span.h4.text-lt Crops
            table.table.m-b-none
              thead
                tr
                  th(style="width: 70%;") Crop Variety
                  th Batches Qty
              tbody
                tr
                  td Bayam Lu Hsien
                  td 7
                tr
                  td Brokoli Green Super
                  td 5
                tr
                  td Kale Nero Toscana
                  td 2
                tr
                  td Lettuce Cos Lobjoits
                  td 7
                tr
                  td Pak Choi
                  td 7
            .panel-footer
              .text-center
                ul.pagination.pagination-sm.m-t-none.m-b-none
                  li
                    a(href="")
                      i.fa.fa-chevron-left
                  li
                    a(href="") 1
                  li
                    a(href="") 2
                  li
                    a(href="") 3
                  li
                    a(href="") 4
                  li
                    a(href="") 5
                  li
                    a(href="")
                      i.fa.fa-chevron-right
        .col-sm-6
          // TASK LIST
          .panel.no-border(style="min-height: 200px;")
            .panel-heading
              .m-b.m-t
                span.pull-right.text-primary
                  a(href="task.html")
                    | See all Tasks 
                    i.fa.fa-angle-double-right
                span.h4.text-lt Tasks
            table.table.m-b-none
              thead
                tr
                  th(style="width: 60%;") Description
                  th(style="width: 20%") Category
                  th Status
              tbody
                tr
                  td
                    a(href="task-detail.html")
                      div Fumigating with Rentokil
                      small.text-muted Due date: 04/01/2018
                      .text-danger Overdue!
                  td
                    span.label.label-area AREA
                  td
                    span.label.label-danger URGENT
                tr
                  td
                    a(href="task-detail.html")
                      div Add new layer of composted soil up to 5cm high
                      small.text-muted 07/01/2018
                  td
                    span.label.label-area AREA
                  td
                    span.label.label-info NORMAL
                tr
                  td
                    a(href="task-detail.html")
                      div Fertilize the bed and clean up the rotten or dry leaves
                      small.text-muted 09/01/2018
                  td
                    span.label.label-area AREA
                  td
                    span.label.label-info NORMAL
            .panel-footer
              .text-center
                ul.pagination.pagination-sm.m-t-none.m-b-none
                  li
                    a(href="")
                      i.fa.fa-chevron-left
                  li
                    a(href="") 1
                  li
                    a(href="") 2
                  li
                    a(href="") 3
                  li
                    a(href="") 4
                  li
                    a(href="") 5
                  li
                    a(href="")
                      i.fa.fa-chevron-right
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
export default {
  name: 'Home',
  computed : {
    ...mapGetters({
      areas: 'getAllAreas'
    })
  },
  mounted () {
    this.fetchAreas()
  },
  methods: {
    ...mapActions([
      'fetchAreas'
    ])
  }
}
</script>

