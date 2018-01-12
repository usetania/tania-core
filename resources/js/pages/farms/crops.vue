<template lang="pug">
  .crops.col
    .wrapper-md
      .row
        .col-sm-6.col-xs-12
          h1.m-t.font-thin.h3.text-black Crops
        .col-sm-6
          ul.nav.nav-tabs.h4.pull-right
            li(role="presentation" class="active"): a(href="#") Batch
            li: a(href="#") Schedules
    .wrapper
      .panel.no-border
        .panel-heading.wrapper
          span.h4.text-lt All Growing Batches on This Farm
          router-link.btn.btn-sm.btn-primary.btn-addon.pull-right(:to="{name: 'FarmCropsCreate'}")
            i.fa.fa-plus
            | Add a New Batch
        table.table.m-b
          thead
            tr
              th(style="witdh: 18%") Crop Variety
              th(style="width: 17%") Batch ID
              th(style="width: 10%") Creation Date
              th(style="width: 10%") Activity Type
              th(style="width: 10%") Initial Area
              th(style="width: 10%") Current Area
          tbody
            tr(v-for="crop in crops")
              td {{ crop.crop_type }}
              td {{ crop.batch_id }}
              td {{ crop.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY HH:mm:ss') }}
              td
                span(v-for="activity in crop.activity_type") {{ activity.name }}
              td
              td {{ crop.initial_area }}
              td
                span(v-for="area in crop.current_area") {{ area.name }}
              td
</template>