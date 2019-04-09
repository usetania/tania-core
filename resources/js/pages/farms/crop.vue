<template lang="pug">
.container-fluid.bottom-space
  .crop-detail(v-if="loading === false")
    modal(v-if="showMoveCropModal" @close="closeModal")
      moveCropTask(:crop="crop")
    modal(v-if="showDumpCropModal" @close="closeModal")
      dumpCropTask(:crop="crop")
    modal(v-if="showHarvestCropModal" @close="closeModal")
      harvestCropTask(:crop="crop")
    modal(v-if="showUploadCropModal" @close="closeModal")
      uploadCropTask(:crop="crop")

    .row
      .col
        router-link(:to="{name: 'FarmCrops'}")
          i.fa.fa-long-arrow-alt-left
          translate Back to Crop Batches

        b-nav
          b-nav-item
            router-link(:to="{ name: 'FarmCrop', params: { id: crop.uid } }")
              translate Basic Info
          b-nav-item
            router-link(:to="{ name: 'FarmCropNotes', params: { id: crop.uid } }")
              translate Tasks &amp; Notes

        b-card
          .panel-heading
            .row
              .col-xs-12.col-sm-12.col-md-6
                h3 {{ crop.inventory.name }}
                .identifier {{ crop.batch_id }}

                small.text-muted.m-t.clear
                  | {{ crop.area_status.seeding }}
                  |
                  translate Seeding
                  | ,
                  | {{ crop.area_status.growing }}
                  |
                  translate Growing
                  | ,
                  | {{ crop.area_status.dumped }}
                  |
                  translate Dumped

              .col-xs-12.col-sm-12.col-md-6
                .row
                  .col-md-6
                    button.btn.btn-success.btn-block(style="cursor: pointer;" @click="showHarvestCropModal = true")
                      i.fa.fa-cut
                      translate Harvest
                  .col-md-6
                    button.btn.btn-danger.btn-block(style="cursor: pointer;" @click="showDumpCropModal = true")
                      i.fa.fa-trash
                      translate Dump
                .row
                  .col-md-6
                    button.btn.btn-primary.btn-block(style="cursor: pointer;" @click="showMoveCropModal = true")
                      i.fa.fa-exchange-alt
                      translate Move
                  .col-md-6
                    button.btn.btn-default.btn-block(style="cursor: pointer;" @click="showUploadCropModal = true")
                      i.fa.fa-camera
                      translate Take Picture

          .row
            // STATUS
            .col
              .row
                .col-sm-6
                  small.text-muted
                    translate Seeding Date
                  h4 {{ crop.initial_area.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                  small.text-muted
                    translate Last Watering
                  h4(v-if="crop.initial_area.last_watered")
                    | {{ crop.initial_area.last_watered | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                    |
                    translate at
                    |
                    | {{ crop.initial_area.last_watered | moment('timezone', 'Asia/Jakarta').format('HH:mm') }}
                    |
                    |
                    span.areatag {{ crop.initial_area.name }}
                  h4(v-else)
                    | -
                  h4(v-for="area in crop.moved_area" v-if="area.last_watered")
                    | {{ area.last_watered | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                    |
                    translate at
                    |
                    | {{ area.last_watered | moment('timezone', 'Asia/Jakarta').format('HH:mm') }}
                    |
                    |
                    span.areatag {{ area.name }}
                .col-sm-6
                  small.text-muted
                    translate Initial Planting
                  h4.m-b
                    | {{ crop.initial_area.initial_quantity }} {{ getCropContainer(crop.container.type, crop.container.quantity) }}
                    |
                    translate on
                    |
                    |
                    span.areatag {{ crop.initial_area.name }}
                  small.text-muted
                    translate Current Quantity
                  h4(v-if="crop.initial_area.current_quantity > 0")
                    | {{ crop.initial_area.current_quantity }} {{ getCropContainer(crop.container.type, crop.container.quantity) }}
                    |
                    translate on
                    |
                    |
                    span.areatag {{ crop.initial_area.name }}
                  h4(v-for="area in crop.moved_area" v-if="area.current_quantity > 0")
                    | {{ area.current_quantity }} {{ getCropContainer(crop.container.type, crop.container.quantity) }}
                    |
                    translate on
                    |
                    |
                    span.areatag {{ area.name }}

          .row
            // ACTIVITY FEEDS
            .col-sm-12
              .cropactivity
                // ACTIVITY
                h4
                  translate Activity
                b-list-group
                  b-list-group-item(v-for="(activity, indexAct) in activities" :key="indexAct")
                    // MOVE
                    .row(v-if="activity.activity_type.code == 'MOVE'")
                      .col
                        .meta-wrapper
                          i.fa.fa-exchange-alt
                          |
                          |
                          translate Moved
                          |
                          |
                          u {{ activity.activity_type.quantity }} {{ getCropContainer(activity.container_type, activity.activity_type.quantity) }}
                          |
                          |
                          translate from
                          |
                          |
                          span.areatag-sm {{ activity.activity_type.source_area_name }}
                          |
                          |
                          translate to
                          |
                          |
                          span.areatag-sm {{ activity.activity_type.destination_area_name }}
                        Timestamp(:timestamp="activity.created_date")
                    // SEEDING
                    .row(v-if="activity.activity_type.code == 'SEED'")
                      .col
                        .meta-wrapper
                          i.fa.fa-utensil-spoon
                          |
                          |
                          translate Seeded
                          |
                          |
                          u {{ activity.activity_type.quantity }} {{ getCropContainer(activity.container_type, activity.activity_type.quantity) }}
                          |
                          |
                          translate of
                          |
                          |
                          span.identifier-sm {{ activity.batch_id }}
                          |
                          |
                          translate on
                          |
                          |
                          span.areatag-sm {{ activity.activity_type.area_name }}
                        Timestamp(:timestamp="activity.created_date")
                    // DUMP
                    .row(v-if="activity.activity_type.code == 'DUMP'")
                      .col
                        .meta-wrapper
                          i.fa.fa-trash
                          |
                          |
                          translate Dumped
                          |
                          |
                          u {{ activity.activity_type.quantity }} {{ getCropContainer(activity.container_type, activity.activity_type.quantity) }}
                          |
                          |
                          translate on
                          |
                          |
                          span.areatag-sm {{ activity.activity_type.source_area_name }}
                        MoreDetail(:data="activity" :description="activity.description")
                        Timestamp(:timestamp="activity.created_date")
                    // PHOTO
                    .row(v-if="activity.activity_type.code == 'PHOTO'")
                      .col
                        i.fa.fa-camera
                        |
                        |
                        MoreDetail(:data="activity" :description="activity.activity_type.description")
                        Timestamp(:timestamp="activity.created_date")
                        br
                        b-img(
                          :src="`/api/farms/crops/${crop.uid}/photos/${activity.activity_type.uid}`"
                          fluid
                        )
                    // HARVEST
                    .row(v-if="activity.activity_type.code == 'HARVEST'")
                      .col
                        .meta-wrapper
                          i.fas.fa-cut
                          |
                          |
                          translate Harvested
                          |
                          |
                          u(v-if="activity.activity_type.type == 'PARTIAL'")
                            translate Partial
                          u(v-if="activity.activity_type.type == 'ALL'")
                            translate All
                          |
                          |
                          translate  of
                          |
                          |
                          b {{ activity.activity_type.produced_gram_quantity }}
                          |
                          |
                          translate Grams
                          |
                          |
                          translate on
                          |
                          |
                          span.areatag-sm {{ activity.activity_type.source_area_name }}
                        Timestamp(:timestamp="activity.created_date")
                    // WATER
                    .row(v-if="activity.activity_type.code == 'WATER'")
                      .col
                        .meta-wrapper
                          i.fa.fa-tint
                          |
                          |
                          translate Area
                          |
                          |
                          u: router-link(:to="{ name: 'FarmArea', params: { id: activity.activity_type.area_id } }") {{ activity.activity_type.area_name }}
                          |
                          |
                          translate Watered
                        Timestamp(:timestamp="activity.created_date")
                    // PESTICIDE
                    .row(v-if="activity.activity_type.code == 'TASK_PEST_CONTROL'")
                      .col-xs-1.text-center
                        i.fa.fa-bug.block.m-b.m-t
                      .col-xs-11
                        div
                          translate Apply
                          u {{ activity.activity_type.material_name }}
                          translate to
                          span.identifier-sm {{ crop.batch_id }}
                          translate on
                          span.areatag-sm {{ activity.activity_type.area_name }}
                        Timestamp(:timestamp="activity.created_date")
                    //SAFETY
                    .row(v-if="activity.activity_type.code == 'TASK_SAFETY'")
                      .col-xs-1.text-center
                        i.fas.fa-shield-alt.block.m-b.m-t
                      .col-xs-11
                        div
                          span.identifier-sm {{ crop.batch_id }}
                          translate on
                          span.areatag-sm {{ activity.activity_type.area_name }}
                          i.fas.fa-long-arrow-alt-right
                          |  {{ activity.activity_type.title }}
                        MoreDetail(:data="activity" :description="activity.activity_type.description")
                        Timestamp(:timestamp="activity.created_date")
                    .row(v-if="activity.activity_type.code == 'TASK_CROP'")
                      .col-xs-1.text-center
                        i.fas.fa-leaf.block.m-b.m-t
                      .col-xs-11
                        div
                          span.identifier-sm {{ crop.batch_id }}
                          translate on
                          span.areatag-sm {{ activity.activity_type.area_name }}
                          i.fas.fa-long-arrow-alt-right
                          |  {{ activity.activity_type.title }}
                        MoreDetail(:data="activity" :description="activity.activity_type.description")
                        Timestamp(:timestamp="activity.created_date")
                    .row(v-if="activity.activity_type.code == 'TASK_NUTRIENT'")
                      .col-xs-1.text-center
                        i.fa.fa-flask.block.m-b.m-t
                      .col-xs-11
                        div
                          translate Apply
                          u {{ activity.activity_type.material_name }}
                          translate on
                          span.areatag-sm {{ activity.activity_type.area_name }}
                        Timestamp(:timestamp="activity.created_date")
                    .row(v-if="activity.activity_type.code == 'TASK_SANITATION'")
                      .col-xs-1.text-center
                        i.fas.fa-medkit.block.m-b.m-t
                      .col-xs-11
                        div
                          span.identifier-sm {{ crop.batch_id }}
                          translate on
                          span.areatag-sm {{ activity.activity_type.area_name }}
                          i.fas.fa-long-arrow-alt-right
                          |  {{ activity.activity_type.title }}
                        MoreDetail(:data="activity" :description="activity.activity_type.description")
                        Timestamp(:timestamp="activity.created_date")
</template>

<script>
import { mapActions } from 'vuex';
import { FindContainer, AddClicked } from '../../stores/helpers/farms/crop';
import { StubCrop, StubNote } from '../../stores/stubs';
import Modal from '../../components/modal.vue';
import MoreDetail from '../../components/more-detail.vue';
import Timestamp from '../../components/timestamp.vue';

export default {
  name: 'FarmCrop',
  components: {
    moveCropTask: () => import('./activities/move-crop-task.vue'),
    dumpCropTask: () => import('./activities/dump-crop-task.vue'),
    harvestCropTask: () => import('./activities/harvest-crop-task.vue'),
    uploadCropTask: () => import('./activities/upload-crop-task.vue'),
    MoreDetail,
    Modal,
    Timestamp,
  },
  data() {
    return {
      loading: true,
      crop: Object.assign({}, StubCrop),
      note: Object.assign({}, StubNote),
      cropNotes: [],
      showMoveCropModal: false,
      showDumpCropModal: false,
      showHarvestCropModal: false,
      showUploadCropModal: false,
      activities: [],
    };
  },
  created() {
    this.loadActivities(this.$route.params.id)
  },
  methods: {
    ...mapActions([
      'fetchActivities',
      'getCropByUid',
    ]),
    closeModal() {
      this.showMoveCropModal = false;
      this.showDumpCropModal = false;
      this.showHarvestCropModal = false;
      this.showUploadCropModal = false;
      this.loadActivities(this.crop.uid);
    },
    getCropContainer(key, count) {
      return FindContainer(key).label + ((count !== 1) ? 's' : '');
    },
    loadActivities(cropId) {
      this.getCropByUid(cropId)
        .then(({ data }) => {
          this.crop = data;
          this.fetchActivities(cropId)
            .then(({ data }) => {
              this.loading = false;
              this.activities = AddClicked(data);
            })
            .catch(error => console.log(error));
        })
        .catch(error => console.log(error));
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
.bottom-space {
  padding-bottom: 60px;
}
</style>
