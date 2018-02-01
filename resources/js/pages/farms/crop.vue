<template lang="pug">
  .crop-detail.hbox(v-if="loading === false")
    modal(v-if="showMoveCropModal" @close="showMoveCropModal = false")
      moveCropTask
    modal(v-if="showDumpCropModal" @close="showDumpCropModal = false")
      dumpCropTask
    modal(v-if="showHarvestCropModal" @close="showHarvestCropModal = false")
      harvestCropTask
    modal(v-if="showUploadCropModal" @close="showUploadCropModal = false")
      uploadCropTask
    modal(v-if="showFertilizerCropModal" @close="showFertilizerCropModal = false")
      fertilizerCropTask
    modal(v-if="showPesticideCropModal" @close="showPesticideCropModal = false")
      pesticideCropTask
    modal(v-if="showPruneCropModal" @close="showPruneCropModal = false")
      pruneCropTask
    modal(v-if="showOtherCropModal" @close="showOtherCropModal = false")
      otherCropTask
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
                            .h3.m-b {{ crop.inventory.variety }}
                            .identifier {{ crop.batch_id }}
                          .col-sm-5
                            small.text-muted Activity Type
                            .h4.m-b.m-t {{ crop.activity_type.total_seeding }} Seeding, {{ crop.activity_type.total_growing }} Growing
                      .panel-body.bg-light.lter.b-b.b-light.m-l-n-xxs.m-r-n-xxs
                        .row
                          .col-sm-6.h5
                            | Seeded on 
                            b {{ crop.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                            |  at 
                            u {{ crop.initial_area.area.name }}
                          .col-sm-6.text-right.h5
                            | Initial Qty: 
                            b {{ crop.container.quantity }} {{ getCropContainer(crop.container.type.code, crop.container.quantity) }}
                      .panel-body.bg-white-only
                        .row.m-t
                          .col-sm-7
                            .m-b
                              .h4.font-bold.m-b Notes
                              form(@submit.prevent="validateBeforeSubmit")
                                textarea.form-control.m-b#content(placeholder="Leave a note here..." rows="2" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('note.content') }" v-model="note.content" name="note.content")
                                span.help-block.text-danger(v-show="errors.has('note.content')") {{ errors.first('crop.container_cell') }}
                                button.btn.btn-xs.btn-success.m-b(type="submit") Add Note
                              ul.list-group.list-group-lg.no-bg.auto
                                li.list-group-item.row(v-for="cropNote in crop.notes")
                                  .col-sm-9
                                    .pull-left.m-r
                                      i.fa.fa-file.block.m-b.m-t
                                    span {{ cropNote.content }}
                                    small.text-muted.clear.text-ellipsis {{ cropNote.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                                  .col-sm-3
                                    button.btn.btn-xs.btn-default.pull-right(v-on:click="deleteNote(cropNote.uid)")
                                      i.fa.fa-trash
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
                                | {{ crop.initial_area.current_quantity }} Plants at 
                                span.areatag-sm {{ crop.initial_area.area.name }}
                              .h4(v-for="area in crop.moved_area")
                                | {{ area.current_quantity }} Plants at 
                                span.areatag-sm {{ area.area.name }}
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
            button.btn.btn-addon.btn-default(style="cursor: pointer;" @click="showMoveCropModal = true")
              i.fa.fa-exchange
              | Move
          .m-b
            button.btn.btn-addon.btn-primary(style="cursor: pointer;" @click="showDumpCropModal = true")
              i.fa.fa-trash
              | Dump
          .m-b
            button.btn.btn-addon.btn-success(style="cursor: pointer;" @click="showHarvestCropModal = true")
              i.fa.fa-scissors
              | Harvest
          .m-b
            button.btn.btn-addon.btn-dark(style="cursor: pointer;" @click="showUploadCropModal = true")
              i.fa.fa-camera
              | Upload Photo
        .wrapper-md
          .m-b.text-md Tasks
          .m-b
            button.btn.btn-addon.btn-warning(style="cursor: pointer;" @click="showFertilizerCropModal = true")
              i.fa.fa-flask
              | Fertilize
          .m-b
            button.btn.btn-addon.btn-danger(style="cursor: pointer;" @click="showPesticideCropModal = true")
              i.fa.fa-bug
              | Pesticide
          .m-b
            button.btn.btn-addon.btn-dark(style="cursor: pointer;" @click="showPruneCropModal = true")
              i.fa.fa-scissors
              | Prune
          .m-b
            button.btn.btn-addon.btn-default(style="cursor: pointer;" @click="showOtherCropModal = true")
              i.fa.fa-tasks
              | Other Task
</template>
<script>
import { FindCropContainer } from '@/stores/helpers/farms/crop'
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
    fertilizerCropTask: () => import('./tasks/fertilizer-crop-task.vue'),
    pesticideCropTask: () => import('./tasks/pesticide-crop-task.vue'),
    pruneCropTask: () => import('./tasks/prune-crop-task.vue'),
    otherCropTask: () => import('./tasks/other-crop-task.vue'),
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
      'createCropNotes',
      'deleteCropNote'
    ]),
    getCropContainer(key, count) {
      return FindCropContainer(key).label + ((count != 1)? 's':'')
    },
    deleteNote(note_uid) {
      this.note.obj_uid = this.$route.params.id
      this.note.uid = note_uid
      this.deleteCropNote(this.note)
        .then(data => this.crop = data)
        .catch(({ data }) => this.message = data)
    },
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
      this.note.obj_uid = this.$route.params.id
      this.createCropNotes(this.note)
        .then(data => this.crop = data)
        .catch(({ data }) => this.message = data)
    },
  }
}
</script>

