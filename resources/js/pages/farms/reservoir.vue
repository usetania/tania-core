<template lang="pug">
  .reservoir-detail.col(v-if="loading === false")
    .wrapper-md
      a.btn.m-b-xs.btn-primary.btn-addon.pull-right(style="cursor: pointer;" id="show-modal" @click="showModal = true")
        i.fa.fa-plus
        | Add Task
      h1.m-n.font-thin.h3.text-black {{ reservoir.name }}
    .wrapper
      .row.basicinfo
        .col-md-4
          .panel
            .panel-heading
              span.h4.text-lt Basic Info
            .panel-body
              .row.m-b
                .col-md-6
                  small.text-muted Source Type
                  .h4.text-lt {{ getReservoirType(reservoir.water_source.type).label }}
                .col-md-6
                  small.text-muted Capacity
                  .h4.text-lt {{ reservoir.water_source.capacity }}
              .row.m-b
                .col-md-6
                  small.text-muted Used In
                  .h4.text-lt(v-for="area in reservoir.installed_to_areas")
                    span.areatag {{ area.name }}
                .col-md-6
                  small.text-muted Device Connected
                  .h4.text-lt 3
              .row.m-b
                .col-md-6
                  small.text-muted Created On
                  .h4.text-lt {{ reservoir.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
        .col-md-8
          .panel
            .panel-heading
              span.h4.text-lt Watering Schedule
            .panel-body
              span.h4.text-lt Graph here
      .row
        .col-sm-6
          .panel
            .panel-heading
              span.h4.text-lt Tasks
            table.table.m-b-none
              thead
                tr
                  th(style="width: 20%") Status
                  th Description
              tbody
                tr
                  td
                    span.label.label-danger URGENT
                  td
                    a(href="task-detail.html")
                      div Drain all the water and cleanse it with cleaning solutions
                      small.text-muted 29/12/2017
                tr
                  td
                    span.label.label-info NORMAL
                  td
                    a(href="task-detail.html")
                      div Check the solenoid valve
                      small.text-muted 27/12/2017
        .col-sm-6.col-xs-12
          .panel
            .panel-heading
              span.h4.text-lt Notes
            .panel-body
              form(@submit.prevent="validateBeforeSubmit")
                .input-group
                  input.form-control.input-sm#content(type="text" placeholder="Create a note" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('note.content') }" v-model="note.content" name="note.content")
                  span.input-group-btn
                    button.btn.btn-sm.btn-success(type="submit")
                      i.fa.fa-send
                  span.help-block.text-danger(v-show="errors.has('note.content')") {{ errors.first('crop.container_cell') }}
            ul.list-group.list-group-lg.no-bg.auto
              li.list-group-item.row(v-for="reservoirNote in reservoir.notes")
                .col-sm-9
                  span {{ reservoirNote.content }}
                  small.text-muted.clear.text-ellipsis {{ reservoirNote.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
                .col-sm-3
                  button.btn.btn-xs.btn-default.pull-right(v-on:click="deleteNote(reservoirNote.uid)")
                    i.fa.fa-trash
</template>

<script>
import { StubReservoir, StubNote } from '@/stores/stubs'
import { FindReservoirType } from '@/stores/helpers/farms/reservoir' 
import { mapActions } from 'vuex'
export default {
  name: 'Reservoir',
  data () {
    return {
      loading: true,
      reservoir: Object.assign({}, StubReservoir),
      note: Object.assign({}, StubNote),
      reservoirNotes: [],
      showModal: false,
    }
  },
  created () {
    this.getReservoirByUid(this.$route.params.id)
      .then(({ data }) =>  {
        this.loading = false
        this.reservoir = data
      })
      .catch(error => console.log(error))
  },
  methods: {
    ...mapActions([
      'getReservoirByUid',
      'createReservoirNotes',
      'deleteReservoirNote',
    ]),
    deleteNote(note_uid) {
      this.note.obj_uid = this.$route.params.id
      this.note.uid = note_uid
      this.deleteReservoirNote(this.note)
        .then(data => this.reservoir = data)
        .catch(({ data }) => this.message = data)
    },
    getReservoirType(key) {
      return FindReservoirType(key)
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
      this.createReservoirNotes(this.note)
        .then(data => this.reservoir = data)
        .catch(({ data }) => this.message = data)
    },
  }
}
</script>

