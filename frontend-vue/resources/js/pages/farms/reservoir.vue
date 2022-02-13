<template lang="pug">
  .container-fluid.bottom-space(v-if="loading === false")
    modal(v-if="showModal" @close="closeModal")
      FarmReservoirTaskForm(:data="reservoir" :asset="asset")
    .row
      .col
        .title-page
          BtnAddNew(
            :title="$gettext('Add Task')"
            customClass="float-right"
            v-on:click.native="openModal()"
          )

          h3.title-page {{ reservoir.name }}

    // basic info
    .row
      .col-xs-12.col-sm-12.col-md-6
        .basicinfo
          b-card(
            :title="$gettext('Basic Info')"
            class="card-ui"
          )
            .row
              .col-xs-12.col-sm-12.col-md-6
                small.text-muted
                  translate Source Type
                h4 {{ getReservoirType(reservoir.water_source.type).label }}
              .col-xs-12.col-sm-12.col-md-6
                small.text-muted
                  translate Capacity
                h4 {{ reservoir.water_source.capacity }}

              .col-xs-12.col-sm-12.col-md-6
                small.text-muted
                  translate Used In
                h4(v-for="area in reservoir.installed_to_area")
                  span.areatag {{ area.name }}

              .col-xs-12.col-sm-12.col-md-6
                small.text-muted
                  translate Created On
                h4 {{ reservoir.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}

      .col-xs-12.col-sm-12.col-md-6
        b-card(
          :title="$gettext('Notes')"
          class="card-ui"
        )
          .panel-body
            b-form(@submit.prevent="validateBeforeSubmit")
              .input-group
                b-form-input#content.form-control.input-sm(
                  type="text"
                  placeholder="Create a note"
                  v-validate="'required'"
                  :class="{'input': true, 'text-danger': errors.has('note.content') }"
                  v-model="note.content"
                  name="note.content"
                )
                b-input-group-append
                  b-button.btn.btn-sm.btn-success(type="submit")
                    i.fa.fa-paper-plane
                span.help-block.text-danger(
                  v-show="errors.has('note.content')"
                )
                  | {{ errors.first('note.content') }}

          b-list-group.list-notes
            b-list-group-item(v-for="reservoirNote in reservoir.notes" :key="reservoirNote.uid")
              .row
                .col-xs-8.col-sm-8.col-md-9.col-lg-10
                  span {{ reservoirNote.content }}
                    small.text-muted.clear.text-ellipsis
                      |
                      | {{ reservoirNote.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}

                .col-xs-4.col-sm-4.col-md-3.col-lg-2
                  button.btn.btn-xs.btn-default.text-center(
                    v-on:click="deleteNote(reservoirNote.uid)"
                  )
                    i.fa.fa-trash
    .row
      .col-12
        b-card(
          :title="$gettext('Tasks')"
          class="card-ui"
        )
          TasksList(:domain="'RESERVOIR'" :asset_id="reservoir.uid" :reload="reload")
</template>

<script>
import { mapActions } from 'vuex';
import { StubReservoir, StubNote } from '../../stores/stubs';
import { FindReservoirType } from '../../stores/helpers/farms/reservoir';
import Modal from '../../components/modal.vue';
import FarmReservoirTaskForm from './tasks/task-form.vue';
import TasksList from './tasks/task-list.vue';
import BtnAddNew from '../../components/common/btn-add-new.vue';
import BtnCancel from '../../components/common/btn-cancel.vue';
import BtnSave from '../../components/common/btn-save.vue';

export default {
  name: 'Reservoir',
  components: {
    FarmReservoirTaskForm,
    TasksList,
    Modal,
    BtnAddNew,
    BtnCancel,
    BtnSave,
  },
  data() {
    return {
      asset: 'Reservoir',
      loading: true,
      note: Object.assign({}, StubNote),
      reload: false,
      reservoir: Object.assign({}, StubReservoir),
      reservoirNotes: [],
      showModal: false,
    };
  },
  created() {
    this.getReservoirByUid(this.$route.params.id)
      .then(({ data }) => {
        this.loading = false;
        this.reservoir = data;
      })
      .catch(error => error);
  },
  methods: {
    ...mapActions([
      'getReservoirByUid',
      'createReservoirNotes',
      'deleteReservoirNote',
    ]),
    closeModal() {
      this.showModal = false;
      this.reload = !this.reload;
    },
    create() {
      this.note.obj_uid = this.$route.params.id;
      this.createReservoirNotes(this.note)
        .then((data) => {
          this.reservoir = data;
          this.note.content = '';
          this.$nextTick(() => this.$validator.reset());
        })
        .catch(({ data }) => {
          this.message = data;
          return this.message;
        });
    },
    deleteNote(uidNote) {
      this.note.obj_uid = this.$route.params.id;
      this.note.uid = uidNote;
      this.deleteReservoirNote(this.note)
        .then((data) => {
          this.reservoir = data;
          return this.reservoir;
        })
        .catch(({ data }) => {
          this.message = data;
          return this.message;
        });
    },
    getReservoirType(key) {
      return FindReservoirType(key);
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
i.fa.fa-paper-plane {
  text-align: center;
}

i.fas.fa-plus {
  text-align: left;
  width: 30px;
}

.title-page {
  margin: 20px 0 30px 0;
}

.bottom-space {
  padding-bottom: 60px;
}

.card-ui {
  margin-bottom: 20px;

  i {
    width: 30px;
  }

  .see-all {
    display: block;
    margin-bottom: 15px;
  }
}

.list-notes {
  margin-top: 20px;
}
</style>
