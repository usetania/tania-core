<template lang="pug">
  .area-detail.col(v-if="loading === false")
    modal(v-if="showModal" @close="closeModal")
      FarmAreaTaskForm(:data="area" :asset="'Area'")
    modal(v-if="showWaterTaskModal" @close="closeModal")
      WaterTaskModal(:crops="areaCrops" :area="area")
    .wrapper-md
      .pull-right
        a#addTaskForm.btn.btn-sm.btn-addon.btn-primary.m-r(style="cursor: pointer;" @click="openModal()")
          i.fas.fa-plus
          | Add Task
        a#waterAreaForm.btn.btn-sm.btn-addon.btn-info(v-if="areaCrops.length > 0" style="cursor: pointer;" @click="showWaterTaskModal = true")
          i.fas.fa-tint
          | Watering
      h1.m-n.font-thin.h3.text-primary {{ area.name }}
      small.text-muted {{ getType(area.type).label }}
    .wrapper-md
      .row
        .col-md-6.col-xs-12
          .panel.basicinfo
            .panel-heading
              span.h4.text-lt Basic info
            .item
              img.img-full(v-if="area.photo.filename.length > 0" :src="'/api/farms/' + farm.uid + '/areas/' + area.uid + '/photos'")
              img.img-full(v-else src="../../../images/no-img.png")
            .list-group.no-radius.alt
              .list-group-item
                span.col-sm-7.text-muted.point Area Size {{ getSizeUnit(area.size.unit.symbol).label }}
                span {{ area.size.value }}
              .list-group-item
                span.col-sm-7.text-muted.point Location
                span {{ getLocation(area.location.code).label }}
              .list-group-item
                span.col-sm-7.text-muted.point Batches
                span {{ area.total_crop_batch }}
              .list-group-item
                span.col-sm-7.text-muted.point Crop Variety
                span {{ area.total_variety }}
              .list-group-item
                span.col-sm-7.text-muted.point Reservoir
                span {{ area.reservoir.name }}
        .col-md-6.col-xs-12
          .panel
            .panel-heading
              span.h4.text-lt Notes
            .panel-body
              form(@submit.prevent="validateBeforeSubmit")
                .input-group
                  input#content.form-control.input-sm(type="text" placeholder="Create a note" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('note.content') }" v-model="note.content" name="note.content")
                  span.input-group-btn
                    button.btn.btn-sm.btn-success(type="submit")
                      i.fas.fa-paper-plane
                  span.help-block.text-danger(v-show="errors.has('note.content')") {{ errors.first('note.content') }}
            ul.list-group.list-group-lg.no-bg.auto
              li.list-group-item.row(v-for="areaNote in area.notes")
                .col-sm-9
                  span {{ areaNote.content }}
                  small.text-muted.clear.text-ellipsis {{ areaNote.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                .col-sm-3
                  button.btn.btn-xs.btn-default.pull-right(v-on:click="deleteNote(areaNote.uid)")
                    i.fas.fa-trash
      //- Ending row

      //- Starting row
      .panel
        .panel-heading
          span.h4.text-lt Current status
        FarmCropsListing(:crops="areaCrops" :domain="'AREA'")
      //- Ending row

      //- Starting row
      .row
        .task-list.col-xs-12
          .panel
            .panel-heading
              span.h4.text-lt Tasks
            TasksList(:domain="'AREA'" :asset_id="area.uid" :reload="reload")
</template>

<script>
import { FindAreaType, FindAreaSizeUnit, FindAreaLocation } from '@/stores/helpers/farms/area'
import { StubArea, StubNote } from '@/stores/stubs'
import { mapActions, mapGetters } from 'vuex'
import Modal from '@/components/modal'
import moment from 'moment-timezone'
export default {
  name: 'Area',
  components: {
    FarmAreaTaskForm: () => import('./tasks/task-form.vue'),
    FarmCropsListing: () => import('./crops-listing.vue'),
    TasksList: () => import('./tasks/task-list.vue'),
    WaterTaskModal: () => import('./tasks/water-task.vue'),
    Modal
  },
  computed: {
    ...mapGetters({
      farm: 'getCurrentFarm'
    })
  },
  created () {
    this.getAreaByUid(this.$route.params.id)
      .then(({ data }) =>  {
        this.area = data
        this.loading = false
        this.fetchAreaCrops(this.area.uid)
          .then(({ data }) =>  {
            this.areaCrops = data
          })
          .catch(error => console.log(error))
      })
      .catch(error => console.log(error))
  },
  data () {
    return {
      area: Object.assign({}, StubArea),
      areaNotes: [],
      areaCrops: [],
      loading: true,
      note: Object.assign({}, StubNote),
      reload: false,
      showModal: false,
      showWaterTaskModal: false,
    }
  },
  methods: {
    ...mapActions([
      'createAreaNotes',
      'deleteAreaNote',
      'fetchAreaCrops',
      'getAreaByUid',
    ]),
    closeModal () {
      this.showModal = false
      this.showWaterTaskModal = false
      this.reload = !this.reload
    },
    create () {
      this.note.obj_uid = this.$route.params.id
      this.createAreaNotes(this.note)
        .then(data => {
          this.area = data
          this.note.content = ''
          this.$nextTick(() => this.$validator.reset())
        })
        .catch(({ data }) => this.message = data)
    },
    deleteNote (note_uid) {
      this.note.obj_uid = this.$route.params.id
      this.note.uid = note_uid
      this.deleteAreaNote(this.note)
        .then(data => this.area = data)
        .catch(({ data }) => this.message = data)
    },
    getLocation (key) {
      return FindAreaLocation(key)
    },
    getSizeUnit (key) {
      return FindAreaSizeUnit(key)
    },
    getType (key) {
      return FindAreaType(key)
    },
    openModal() {
      this.data = {}
      this.showModal = true
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
