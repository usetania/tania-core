<template lang="pug">
.container-fluid.bottom-space
  .row
    .col
      h3.title-page
        translate Water Reservoir
  .row
    .col
      modal(v-if="showModal" @close="showModal = false")
        FarmReservoirForm(:data="data")

      BtnAddNew(:title="$gettext('Add Reservoir')" v-on:click.native="openModal()")

  .table-responsive.table-wrapper
    table.table
      thead
        tr
          th
            translate Name
          th
            translate Created On
          th
            translate Source Type
          th
            translate Capacity
          th
            translate Used In
          th
      tbody
        tr(v-for="reservoir in reservoirs")
          td: router-link(:to="{ name: 'FarmReservoir', params: { id: reservoir.uid } }")
            u {{ reservoir.name }}
          td {{ reservoir.created_date | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
          td {{ getType(reservoir.water_source.type).label }}
          td {{ reservoir.water_source.capacity }}
          td
            span(v-for="(area, index) in reservoir.installed_to_area")
              | {{ area.name }}
              span(v-if="index+1 < reservoir.installed_to_area.length") ,
          td
            a.pull-right(style="cursor: pointer;" @click="openModal(reservoir)")
              i.fa.fa-edit
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import { FindReservoirType } from '../../stores/helpers/farms/reservoir';
import Modal from '../../components/modal.vue';
import BtnAddNew from '../../components/common/btn-add-new.vue';
import FarmReservoirForm from './reservoirs-form.vue';

export default {
  name: 'FarmReservoirs',
  components: {
    FarmReservoirForm,
    Modal,
    BtnAddNew,
  },
  data() {
    return {
      showModal: false,
    };
  },
  computed: {
    ...mapGetters({
      reservoirs: 'getAllReservoirs',
    }),
  },
  mounted() {
    this.fetchReservoirs();
  },
  methods: {
    ...mapActions([
      'fetchReservoirs',
    ]),
    getType(key) {
      return FindReservoirType(key);
    },
    openModal(data) {
      this.showModal = true;
      if (data) {
        this.data = data;
      } else {
        this.data = {};
      }
    },
  },
};
</script>

<style lang="scss" scoped>
h3.title-page {
  margin: 20px 0 30px 0;
}

i {
  text-align: left;
  width: 30px;
}

.table-wrapper {
  margin-top: 20px;
}

.bottom-space {
  padding-bottom: 60px;
}
</style>
