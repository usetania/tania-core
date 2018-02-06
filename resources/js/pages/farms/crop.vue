<template lang="pug">
  .crop-detail(v-if="loading === false")
    modal(v-if="showMoveCropModal" @close="showMoveCropModal = false")
      moveCropTask(:crop="crop")
    modal(v-if="showDumpCropModal" @close="showDumpCropModal = false")
      dumpCropTask(:crop="crop")
    modal(v-if="showHarvestCropModal" @close="showHarvestCropModal = false")
      harvestCropTask(:crop="crop")
    modal(v-if="showUploadCropModal" @close="showUploadCropModal = false")
      uploadCropTask(:crop="crop")
    modal(v-if="showFertilizerCropModal" @close="showFertilizerCropModal = false")
      cropTask(:crop="crop" :isfertilizer="true")
    modal(v-if="showPesticideCropModal" @close="showPesticideCropModal = false")
      cropTask(:crop="crop" :ispesticide="true")
    modal(v-if="showPruneCropModal" @close="showPruneCropModal = false")
      cropTask(:crop="crop" :isprune="true")
    modal(v-if="showOtherCropModal" @close="showOtherCropModal = false")
      cropTask(:crop="crop" :isother="true")
    .col
      .row.wrapper-md
        .col-xs-8.col-xs-offset-2
          .m-t.m-b
            a.h5.text-lt.m-b(href="#/crops")
              i.fa.fa-long-arrow-alt-left.m-r
              | Back to Crop Batches
          ul.nav.nav-tabs.h4
            li.active
              router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") Basic Info
            li
              router-link(:to="{ name: 'FarmCropNotes', params: { id: crop.uid } }")  Tasks & Notes
          .panel
            .panel-heading.b-b.b-light.wrapper
              .row
                .col-sm-7.m-t.m-b
                  .h3.text-lt.m-b {{ crop.inventory.name }}
                  .identifier {{ crop.batch_id }}
                  small.text-muted.m-t.clear {{ crop.activity_type.total_seeding }} Seeding, {{ crop.activity_type.total_growing }} Growing, 5 Dumped
                .col-sm-5.m-t.m-b
                  .row
                    .col-sm-6.m-b
                      button.btn.btn-success.btn-block(style="cursor: pointer;" @click="showHarvestCropModal = true")
                        i.fa.fa-cut.m-r
                        | Harvest
                    .col-sm-6.m-b
                      button.btn.btn-danger.btn-block(style="cursor: pointer;" @click="showDumpCropModal = true")
                        i.fa.fa-trash.m-r
                        | Dump
                  .row
                    .col-sm-6
                      button.btn.btn-primary.btn-block(style="cursor: pointer;" @click="showMoveCropModal = true")
                        i.fa.fa-exchange.m-r
                        | Move
                    .col-sm-6
                      button.btn.btn-default.btn-block(style="cursor: pointer;" @click="showUploadCropModal = true")
                        i.fa.fa-camera.m-r
                        | Take Picture
            .panel-body.bg-white-only.b-light
              .row
                // STATUS
                .col-sm-12
                  .hbox.bg-light.lter.wrapper-md
                    .row
                      .col-sm-6
                        small.text-muted Seeding Date
                        .h4.m-b {{ crop.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                        small.text-muted Last Watering
                        .h4.m-b
                          | 05/12/2017 at 16:42 on 
                          span.areatag Frontyard Garden
                      .col-sm-6
                        small.text-muted Initial Planting
                        .h4.m-b
                          | {{ crop.initial_area.initial_quantity }} {{ getCropContainer(crop.container.type.code, crop.container.quantity) }} on 
                          span.areatag {{ crop.initial_area.area.name }}
                        small.text-muted Current Quantity
                        .h4.m-b
                          | 12 Pots on 
                          span.areatag Florania
                        .h4
                          | 13 Pots on 
                          span.areatag Frontyard Garden
              .row.m-t
                // ACTIVITY FEEDS
                .col-sm-12
                  .cropactivity
                    // ACTIVITY
                    .h4.font-bold.m-b.clearfix Activity
                    ul.list-group.no-bg.no-borders.pull-in
                      li.list-group-item
                        .row
                          .col-xs-1.text-center
                            i.fa.fa-medkit.block.m-b.m-t
                          .col-xs-11
                            div
                              span.areatag-sm Frontyard Garden
                              i.fa.fa-long-arrow-alt-right
                              |  Prune leaves from bottom part to center, and bend the stems.
                            p.small
                              span#moreless4.hide
                                | Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                            div
                              a(href="" ui-toggle-class="show" target="#moreless4")
                                small.text
                                  | Read details 
                                  i.fa.fa-angle-down
                                small.text-active
                                  | Close details 
                                  i.fa.fa-angle-up
                            small.text-muted 02/02/2018 at 16:21
</template>
<script>
import { FindContainer } from '@/stores/helpers/farms/crop'
import { mapActions } from 'vuex'
import { StubCrop, StubNote } from '@/stores/stubs'
import Modal from '@/components/modal'
export default {
  name: 'FarmCrop',
  components: {
    moveCropTask: () => import('./activities/move-crop-task.vue'),
    dumpCropTask: () => import('./activities/dump-crop-task.vue'),
    harvestCropTask: () => import('./activities/harvest-crop-task.vue'),
    uploadCropTask: () => import('./activities/upload-crop-task.vue'),
    cropTask: () => import('./tasks/crop-task.vue'),
    Modal
  },
  data () {
    return {
      loading: true,
      crop: Object.assign({}, StubCrop),
      note: Object.assign({}, StubNote),
      cropNotes: [],
      showMoveCropModal: false,
      showDumpCropModal: false,
      showHarvestCropModal: false,
      showUploadCropModal: false,
      showFertilizerCropModal: false,
      showPesticideCropModal: false,
      showPruneCropModal: false,
      showOtherCropModal: false,
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
      'getCropByUid',
    ]),
    getCropContainer(key, count) {
      return FindContainer(key).label + ((count != 1)? 's':'')
    },
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
  }
}
</script>

