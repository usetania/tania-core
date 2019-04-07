<template lang="pug">
  .container-fluid.bottom-space(v-if="loading === false")
    modal(v-if="showModal" @close="closeModal")
      FarmAreaTaskForm(:data="area" :asset="'Area'")
    modal(v-if="showWaterTaskModal" @close="closeModal")
      WaterTaskModal(:crops="areaCrops" :area="area")

    .row
      .col
        h3.title-page {{ area.name }}
          |
          |
          small.text-muted {{ getType(area.type).label }}

    .row
      .col
        a#addTaskForm.btn.btn-sm.btn-primary(style="cursor: pointer;" @click="openModal()")
          i.fas.fa-plus
          translate Add Task

        a#waterAreaForm.btn.btn-sm.btn-info(
          v-if="areaCrops.length > 0" style="cursor: pointer;"
          @click="showWaterTaskModal = true"
        )
          i.fas.fa-tint
          translate Watering

    .cards-wrapper
      .row
        .col-xs-12.col-sm-12.col-md-5
          b-card(
            :title="$gettext('Basic info')"
            :img-src="areaImage"
            class="card-ui"
          )

            b-list-group
              b-list-group-item
                span.col-sm-7.text-muted.point
                  translate Area Size
                  |
                  | {{ getSizeUnit(area.size.unit.symbol).label }}
                span {{ area.size.value }}
              b-list-group-item
                span.col-sm-7.text-muted.point
                  translate Location
                span {{ getLocation(area.location.code).label }}
              b-list-group-item
                span.col-sm-7.text-muted.point
                  translate Batches
                span {{ area.total_crop_batch }}
              b-list-group-item
                span.col-sm-7.text-muted.point
                  translate Crop Variety
                span {{ area.total_variety }}
              b-list-group-item
                span.col-sm-7.text-muted.point
                  translate Reservoir
                span {{ area.reservoir.name }}
        .col-xs-12.col-sm-12.col-md-7
          b-card(
            :text="$gettext('Current Status')"
            class="card-ui"
          )
            FarmCropsListing(:crops="areaCrops" :domain="'AREA'")
      //- Ending row

      //- Starting row
      .row
        .col-md-6.col-xs-12
          b-card(
            :title="$gettext('Notes')"
            class="card-ui"
          )
            b-form(@submit.prevent="validateBeforeSubmit")
              .input-group
                b-form-input#content.form-control.input-sm(type="text" placeholder="Create a note" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('note.content') }" v-model="note.content" name="note.content")
                b-input-group-append
                  b-button.btn.btn-sm.btn-success(type="submit")
                    i.fas.fa-paper-plane
                span.help-block.text-danger(v-show="errors.has('note.content')") {{ errors.first('note.content') }}
            b-list-group.list-notes
              b-list-group-item(v-for="areaNote in area.notes" :key="areaNote.uid")
                .row
                  .col-xs-8.col-sm-8.col-md-9.col-lg-10
                    span {{ areaNote.content }}
                      small.text-muted.clear.text-ellipsis
                        |
                        | {{ areaNote.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                  .col-xs-4.col-sm-4.col-md-3.col-lg-2
                    button.btn.btn-xs.btn-default.text-center(v-on:click="deleteNote(areaNote.uid)")
                      i.fas.fa-trash
        .col-md-6.col-xs-12.task-list
          b-card(
            :title="$gettext('Tasks')"
            class="card-ui"
          )
            TasksList(:domain="'AREA'" :asset_id="area.uid" :reload="reload")
      //- Ending row
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { FindAreaType, FindAreaSizeUnit, FindAreaLocation } from '../../stores/helpers/farms/area';
import { StubArea, StubNote } from '../../stores/stubs';
import Modal from '../../components/modal.vue';
import noImage from '../../../images/no-img.png';

export default {
  name: 'Area',
  components: {
    FarmAreaTaskForm: () => import('./tasks/task-form.vue'),
    FarmCropsListing: () => import('./crops-listing.vue'),
    TasksList: () => import('./tasks/task-list.vue'),
    WaterTaskModal: () => import('./tasks/water-task.vue'),
    Modal,
  },
  data() {
    return {
      area: Object.assign({}, StubArea),
      areaNotes: [],
      areaCrops: [],
      loading: true,
      note: Object.assign({}, StubNote),
      reload: false,
      showModal: false,
      showWaterTaskModal: false,
    };
  },
  computed: {
    ...mapGetters({
      farm: 'getCurrentFarm',
    }),
    areaImage() {
      let image = noImage;
      if (this.area.photo.filename.length > 0) {
        image = `/api/farms/${this.farm.uid}/areas/${this.area.uid}/photos`;
      }
      return image;
    },
  },
  created() {
    this.getAreaByUid(this.$route.params.id)
      .then(({ data }) => {
        this.area = data;
        this.loading = false;

        this.fetchAreaCrops(this.area.uid)
          .then(({ data }) => {
            this.areaCrops = data;
          })
          .catch(error => error);
      })
      .catch(error => error);
  },
  methods: {
    ...mapActions([
      'createAreaNotes',
      'deleteAreaNote',
      'fetchAreaCrops',
      'getAreaByUid',
    ]),
    closeModal() {
      this.showModal = false;
      this.showWaterTaskModal = false;
      this.reload = !this.reload;
    },
    create() {
      this.note.obj_uid = this.$route.params.id;
      this.createAreaNotes(this.note)
        .then((data) => {
          this.area = data;
          this.note.content = '';
          this.$nextTick(() => this.$validator.reset());
        })
        .catch(({ data }) => {
          this.message = data;
        });
    },
    deleteNote(noteUid) {
      this.note.obj_uid = this.$route.params.id;
      this.note.uid = noteUid;
      this.deleteAreaNote(this.note)
        .then((data) => {
          this.area = data;
        })
        .catch(({ data }) => {
          this.message = data;
        });
    },
    getLocation(key) {
      return FindAreaLocation(key);
    },
    getSizeUnit(key) {
      return FindAreaSizeUnit(key);
    },
    getType(key) {
      return FindAreaType(key);
    },
    openModal() {
      this.data = {};
      this.showModal = true;
    },
    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.create();
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
i.fas.fa-plus,
i.fas.fa-tint {
  text-align: left;
  width: 30px;
}

#addTaskForm {
  margin-right: 20px;
}

h3.title-page {
  margin: 20px 0 30px 0;
}

.cards-wrapper {
  margin-top: 20px;

  .card-ui {
    margin-bottom: 20px;

    i {
      width: 30px;
    }
  }
}
.bottom-space {
  padding-bottom: 60px;
}
.list-notes {
  margin-top: 20px;
}
</style>
