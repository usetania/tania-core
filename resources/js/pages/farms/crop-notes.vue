<template lang="pug">
  .crop-detail.col(v-if="loading === false")
    modal(v-if="showTaskModal" @close="closeModal")
      cropTask(:crop="crop" :data="data")
    .row.wrapper-md
      .col-xs-8.col-xs-offset-2
        .m-t.m-b
          a.h5.text-lt.m-b(href="#/crops")
            i.fa.fa-long-arrow-alt-left.m-r
            translate Back to Crop Batches
        ul.nav.nav-tabs.h4
          li
            router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }")
              translate Basic Info
          li.active
            router-link(:to="{ name: 'FarmCropNotes', params: { id: crop.uid } }")
              translate Tasks & Notes
        .panel
          .panel-heading
            .row.m-t
              .col-sm-8
                .h4.font-bold
                  translate Tasks
              .col-sm-4.text-right
                a.btn.btn-sm.btn-primary.btn-addon(style="cursor: pointer;" @click="openModal()")
                  i.fa.fa-plus
                  translate Add Task
          .panel-body.bg-white-only
            .row
              TasksList(:domain="'CROP'" :asset_id="crop.uid" :reload="reload"  @openModal="openModal")
              .col-sm-12
                .h4.font-bold.m-b
                  translate Notes
                .row
                  form(@submit.prevent="validateBeforeSubmit")
                    .col-xs-10
                      textarea.form-control.m-b#content(placeholder="Leave a note here..." rows="2" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('note.content') }" v-model="note.content" name="note.content")
                      span.help-block.text-danger(v-show="errors.has('note.content')") {{ errors.first('crop.container_cell') }}
                    .col-xs-2
                      button.btn.btn-success.pull-right.m-b(type="submit")
                        translate Add Notes
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
import { FindContainer } from '../../stores/helpers/farms/crop'
import { mapActions } from 'vuex'
import { StubCrop, StubNote } from '../../stores/stubs'
import Modal from '../../components/modal.vue'
export default {
  name: 'FarmCropNotes',
  components: {
    cropTask: () => import('./tasks/crop-task-form.vue'),
    TasksList: () => import('./tasks/task-list.vue'),
    Modal
  },
  data () {
    return {
      crop: Object.assign({}, StubCrop),
      cropNotes: [],
      data: {},
      loading: true,
      note: Object.assign({}, StubNote),
      reload: false,
      showTaskModal: false,
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
      'createCropNotes',
      'deleteCropNote',
      'getCropByUid',
    ]),
    closeModal () {
      this.showTaskModal = false
      this.reload = !this.reload
    },
    create () {
      this.note.obj_uid = this.$route.params.id
      this.createCropNotes(this.note)
        .then(data => this.crop = data)
        .catch(() => this.$toasted.error('Error in note submission'))
    },
    deleteNote(note_uid) {
      this.note.obj_uid = this.$route.params.id
      this.note.uid = note_uid
      this.deleteCropNote(this.note)
        .then(data => this.crop = data)
        .catch(() => this.$toasted.error('Error in note deletion'))
    },
    getCropContainer(key, count) {
      return FindContainer(key).label + ((count != 1)? 's':'')
    },
    openModal(data) {
      if (data) {
        this.data = data
      } else {
        this.data = {}
      }
      this.showTaskModal = true
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

