<template lang="pug">
//-
  .hbox
    .col
      .vbox
        .row-row
          .cell
            .cell-inner
              .wrapper-md
                .m-b
                  a.h5.text-lt.m-b(href="#/crops")
                    i.fa.fa-long-arrow-left
                    |  Back to Crop Batches
                .row
                  .panel
                    .panel-heading.b-b.b-light.wrapper
                      .row
                        .col-sm-7
                          .h3.m-b Tomato Ailsa Craig
                          .identifier tom-ail-cra-3nov
                        .col-sm-5
                          small.text-muted Activity Type
                          .h4.m-b.m-t 12 Seeding, 13 Growing
                    .panel-body.bg-light.lter.b-b.b-light.m-l-n-xxs.m-r-n-xxs
                      .row
                        .col-sm-6.h5
                          | Seeded on 
                          b 03/11/2017
                          |  at 
                          u Florania
                        .col-sm-6.text-right.h5
                          | Initial Qty: 
                          b 30 Pots
                    .panel-body.bg-white-only
                      .row.m-t
                        .col-sm-7
                          .m-b
                            .h4.font-bold.m-b Notes
                            textarea.form-control.m-b(placeholder="Leave a note here..." rows="2")
                            button.btn.btn-xs.btn-success.m-b Add Notes
                            ul.list-group.list-group-lg.no-bg.auto
                              li.list-group-item.row
                                .pull-left.m-r
                                  i.fa.fa-file.block.m-b.m-t
                                span If there's any wilted or damaged seedling, just dump it right away
                                small.text-muted.clear.text-ellipsis 05/11/2017
                          .cropactivity
                            .h4.font-bold.m-b.clearfix Activity
                            ul.list-group.no-bg.no-borders.pull-in
                              li.list-group-item
                                .pull-left.m-r
                                  i.fa.fa-bug.block.m-b.m-t
                                .clear
                                  div
                                    | Apply 
                                    u Neem Oil
                                    |  on 
                                    span.areatag-sm Frontyard Garden
                                  small.text-muted 05/12/2017 at 09:31
                        .col-sm-5
                          // QTY & TASK
                          .hbox.bg-light.lter.wrapper(style="min-height: 40px;")
                            small.text-muted Current Quantity
                            .h4.m-b.m-t
                              | 12 Plants at 
                              span.areatag-sm Florania
                            .h4
                              | 13 Plants at 
                              span.areatag-sm Frontyard Garden
                          .m-t
                            .h4.font-bold Tasks
                            ul.list-group.no-bg.no-borders.pull-in
                              li.list-group-item
                                .pull-left.m-r
                                  .checkbox
                                    label.i-checks
                                      input(type="checkbox" name="")
                                      i
                                .clear
                                  div
                                    | Fertilize with 
                                    u Trifecta Plus
                                    |  on 
                                    span.areatag-sm Frontyard Garden
                                  small.text-muted Due 20/12/2017
                              li.list-group-item
                                .pull-left.m-r
                                  .checkbox
                                    label.i-checks
                                      input(type="checkbox" name="")
                                      i
                                .clear
                                  div
                                    | Prune leaves on 
                                    span.areatag-sm Frontyard Garden
                                  small.text-muted Due 24/12/2017
    .col.w-md.bg-light.lter.b-l.bg-auto.no-border-xs
      .wrapper-md
        .m-b.text-md Activity
        .m-b
          button.btn.btn-addon.btn-default
            i.fa.fa-exchange
            | Move
        .m-b
          button.btn.btn-addon.btn-primary
            i.fa.fa-trash
            | Dump
        .m-b
          button.btn.btn-addon.btn-success
            i.fa.fa-scissors
            | Harvest
        .m-b
          button.btn.btn-addon.btn-dark
            i.fa.fa-camera
            | Upload Photo
      .wrapper-md
        .m-b.text-md Tasks
        .m-b
          button.btn.btn-addon.btn-warning
            i.fa.fa-flask
            | Fertilize
        .m-b
          button.btn.btn-addon.btn-danger
            i.fa.fa-bug
            | Pesticide
        .m-b
          button.btn.btn-addon.btn-dark
            i.fa.fa-scissors
            | Prune
        .m-b
          button.btn.btn-addon.btn-default
            i.fa.fa-tasks
            | Other Task
</template>
<script>
import { mapActions } from 'vuex'
import { StubCrop } from '@/stores/stubs'
export default {
  name: 'FarmCrop',
  data () {
    return {
      loading: true,
      crop: Object.assign({}, StubCrop)
    }
  },
  created () {
    this.getCropByUid(this.$route.params.id)
      .then(({ data }) =>  {
        this.loading = false
        this.crop = data
      })
      .catch(error => console.log(error))
  },
  methods: {
    ...mapActions([
      'getCropByUid'
    ])
  }
}
</script>

