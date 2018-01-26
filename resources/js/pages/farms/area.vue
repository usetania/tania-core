<template lang="pug">
  .area-detail.col(v-if="loading === false")
    .wrapper-md
      h1.m-n.font-thin.h3.text-black {{ area.name }}
      small.text-muted {{ getType(area.type).label }}
    .wrapper-md
      .row
        .col-md-4.col-xs-12
          .panel.basicinfo
            .panel-heading
              span.h4.text-lt Basic info
            .item
              img.img-full(src="../../../images/germination.jpg")
            .list-group.no-radius.alt
              .list-group-item
                span.col-sm-7.text-muted.point Area Size {{ getSizeUnit(area.size.unit.symbol).label }}
                span {{ area.size.value }}
              .list-group-item
                span.col-sm-7.text-muted.point Location
                span.text-truncate {{ getLocation(area.location).label }}
              .list-group-item
                span.col-sm-7.text-muted.point Batches
                span {{ area.total_crop_batch }}
              .list-group-item
                span.col-sm-7.text-muted.point Crop Variety
                span {{ area.total_variety }}
              .list-group-item
                span.col-sm-7.text-muted.point Connected Device
                span 5
              .list-group-item
                span.col-sm-7.text-muted.point Reservoir
                span {{ area.reservoir.name }}
        .col-md-8
          .panel
            .panel-heading
              span.pull-right
                i.fa.fa-cog
              span.h4.text-lt Current status
            .panel-body
      //- Ending row

      //- Starting row
      .panel
        .panel-heading
          span.h4.text-lt Current status
        table.table.m-b-none
          thead
            tr
              th Crop Variety
              th Batch ID
              th Seeding Date
              th Days Since Seeding
              th Quantity
          tbody
            tr
              td.text-lt Rosemary Primed
              td
                span.identified ros-pri-1nov
              td 01/11/2017
              td 32
              td 42 Pots
      //- Ending row

      //- Starting row
      .row
        .col-sm-6.col-xs-12
          .panel
            .panel-heading
              span.h4.text-lt Tasks
            table.table.m-b-none
              thead
                tr
                  th Status
                  th Description
              tbody
                tr
                  td
                    span.label.label-danger URGENT
                  td
                    a
                      div Fumigating with Rentokil
                      small.text-muted 01/01.2018
                tr
                  td
                    span.label.label-info NORMAL
                  td
                    a
                      div Fumigating with Rentokil
                      small.text-muted 01/01.2018
        .col-sm-6.col-xs-12
          .panel
            .panel-heading
              span.h4.text-lt Notes
            .panel-body
              .input-group
                input.form-control.input-sm(type="text" placeholder="Create a note")
                span.input-group-btn
                  button.btn.btn-sm.btn-success(type="submit")
                    i.fa.fa-send
            ul.list-group.list-group-lg.no-bg.auto
              li.list-group-item.row
                .col-sm-9
                  span The screen must be checked every month to make sure there are no rips or holes or unfastened
                  small.text-muted.clear.text-ellipsis 19/11/2017
                .col-sm-3
                  button.btn.btn-xs.btn-default.pull-right
                    i.fa.fa-trash
      //- Ending row
</template>

<script>
import { FindAreaType, FindAreaSizeUnit, FindAreaLocation } from '@/stores/helpers/farms/area'
import { StubArea } from '@/stores/stubs'
import { mapActions } from 'vuex'
export default {
  name: 'Area',
  data () {
    return {
      loading: true,
      area: Object.assign({}, StubArea)
    }
  },
  created () {
    this.getAreaByUid(this.$route.params.id)
      .then(({ data }) =>  {
        this.loading = false
        this.area = data
      })
      .catch(error => console.log(error))
  },
  methods: {
    ...mapActions([
      'getAreaByUid'
    ]),
    getType(key) {
      return FindAreaType(key)
    },
    getSizeUnit(key) {
      return FindAreaSizeUnit(key)
    },
    getLocation(key) {
      return FindAreaLocation(key)
    }
  }
}
</script>

