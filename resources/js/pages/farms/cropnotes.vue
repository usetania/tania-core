<template lang="pug">
  .crop-detail(v-if="loading === false")
    .col
      .row.wrapper-md
        .col-xs-8.col-xs-offset-2
          .m-t.m-b
            a.h5.text-lt.m-b(href="crop.html")
              i.fas.fa-long-arrow-alt-left.m-r
              | Back to Crop Batches
          ul.nav.nav-tabs.h4
            li
              router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }") Basic Info
            li.active
              router-link(:to="{ name: 'FarmCropNotes', params: { id: crop.uid } }")  Tasks & Notes
          .panel
            .panel-heading
              .row.m-t
                .col-sm-8
                  .h4.font-bold Tasks
                .col-sm-4.text-right
                  a.btn.btn-sm.btn-primary.btn-addon(href="crop-task-add.html")
                    i.fas.fa-plus
                    | Add Task
            .panel-body.bg-white-only
              .row
                .col-sm-12
                  // TASKS
                  ul.list-group.no-bg.no-borders.pull-in
                    li.list-group-item
                      .row
                        .col-xs-9
                          .pull-left.m-r
                            .checkbox
                              label.i-checks
                                input(type="checkbox" name="")
                                i
                          .clear
                            div
                              | Apply 
                              u Trifecta Plus
                              |  to 
                              span.identifier-sm tom-ail-cra-3nov
                              |  on 
                              span.areatag-sm Frontyard Garden
                            p.small
                              span#moreless1.hide
                                | Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                            p
                              a(href="" ui-toggle-class="show" target="#moreless1")
                                small.text
                                  | Read Details
                                  i.fa.fa-angle-down.m-l
                                small.text-active
                                  | Close Details
                                  i.fa.fa-angle-up.m-l
                            small.text-muted Due Date: 20/12/2017
                            span.status.status-normal NORMAL
                        .col-xs-2
                          span.label.label-nutrient NUTRIENT
                        .col-xs-1.text-right
                          a.h3(href="#")
                            i.fas.fa-edit
                    li.list-group-item
                      .row
                        .col-xs-9
                          .pull-left.m-r
                            .checkbox
                              label.i-checks
                                input(type="checkbox" name="")
                                i
                          .clear
                            div
                              span.identifier-sm tom-ail-cra-3nov
                              |  on 
                              span.areatag-sm Frontyard Garden
                              i.fas.fa-long-arrow-alt-right
                              |  Prune all suckers and dead leaves
                            small.text-muted Due 24/12/2017
                            span.status.status-normal NORMAL
                        .col-xs-2
                          span.label.label-sanitation SANITATION
                        .col-xs-1.text-right
                          a.h3(href="#")
                            i.fas.fa-edit
                    li.list-group-item
                      .row
                        .col-xs-9
                          .pull-left.m-r
                            .checkbox
                              label.i-checks
                                input(type="checkbox" name="")
                                i
                          .clear
                            div
                              span.identifier-sm tom-ail-cra-3nov
                              |  on 
                              span.areatag-sm Florania
                              i.fas.fa-long-arrow-alt-right
                              |  Move the rest of seedlings to growing area
                            p.small
                              span#moreless2.hide
                                | Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                            p
                              a(href="" ui-toggle-class="show" target="#moreless2")
                                small.text
                                  | Read Details
                                  i.fa.fa-angle-down.m-l
                                small.text-active
                                  | Close Details
                                  i.fa.fa-angle-up.m-l
                            small.text-muted Due 27/12/2017
                            span.status.status-urgent URGENT
                        .col-xs-2
                          span.label.label-general GENERAL
                        .col-xs-1.text-right
                          a.h3(href="#")
                            i.fas.fa-edit
                    li.list-group-item
                      .row
                        .col-xs-9
                          .pull-left.m-r
                            .checkbox
                              label.i-checks
                                input(type="checkbox" name="")
                                i
                          .clear
                            div
                              | Apply 
                              u Neem Oil
                              |  to 
                              span.identifier-sm tom-ail-cra-3nov
                              |  on 
                              span.areatag-sm Florania
                            small.text-muted Due 27/12/2017
                            span.status.status-normal NORMAL
                        .col-xs-2
                          span.label.label-pestcontrol PEST CONTROL
                        .col-xs-1.text-right
                          a.h3(href="#")
                            i.fas.fa-edit
                .col-sm-12
                  // QTY & TASK
                  .h4.font-bold.m-b Notes
                  .row
                    form(@submit.prevent="validateBeforeSubmit")
                      .col-xs-10
                        textarea.form-control.m-b#content(placeholder="Leave a note here..." rows="2" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('note.content') }" v-model="note.content" name="note.content")
                        span.help-block.text-danger(v-show="errors.has('note.content')") {{ errors.first('crop.container_cell') }}
                      .col-xs-2
                        button.btn.btn-success.pull-right.m-b(type="submit") Add Notes
                  ul.list-group.list-group-lg.no-bg.auto
                    li.list-group-item.row(v-for="cropNote in crop.notes")
                      .col-xs-10
                        .pull-left.m-r
                          i.fa.fa-file.block.m-b.m-t
                        span {{ cropNote.content }}
                        small.text-muted.clear.text-ellipsis {{ cropNote.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                      .col-xs-2
                        button.btn.btn-xs.btn-default.pull-right(v-on:click="deleteNote(cropNote.uid)")
                          i.fa.fa-trash
</template>
<script>
import { FindContainer } from '@/stores/helpers/farms/crop'
import { mapActions } from 'vuex'
import { StubCrop, StubNote } from '@/stores/stubs'
import Modal from '@/components/modal'
export default {
  name: 'FarmCropNotes',
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
      'createCropNotes',
      'deleteCropNote'
    ]),
    getCropContainer(key, count) {
      return FindContainer(key).label + ((count != 1)? 's':'')
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

