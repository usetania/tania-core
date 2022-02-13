<template lang="pug">
  .container-fluid.bottom-space
    .row
      .col
        h3.title-page
          translate Areas
    .row
      .col
        modal(v-if="showModal" @close="showModal = false")
          farmAreaForm(:data="data")

        BtnAddNew(:title="$gettext('Add Area')" v-on:click.native="openModal()")

    .cards-wrapper
      .row
        .col-xs-12.col-sm-12.col-md-6.col-lg-4.col-xl-4(v-for="area in areas")
          b-card(
            class="card-ui"
          )
            .panel-heading.description
              b-card-title
                router-link(:to="{ name: 'FarmArea', params: { id: area.uid } }") {{ area.name }}
                a.float-right(
                  v-if="area.plant_quantity === 0"
                  style="cursor: pointer;"
                  @click="openModal(area)"
                )
                  i.fa.fa-edit

              b-card-text.small.text-muted {{ getType(area.type).label }}

            .row
              .col-4
                b-card-text.small.text-muted Size ({{ getSizeUnit(area.size.unit.symbol).label }})
                span.text-md {{ area.size.value }}
              .col-4
                b-card-text.small.text-muted.block
                  translate Batches
                span.text-md {{ area.total_crop_batch }}
              .col-4
                b-card-text.small.text-muted.block
                  translate Quantity
                span.text-md {{ area.plant_quantity }}
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { FindAreaType, FindAreaSizeUnit } from '../../stores/helpers/farms/area';
import Modal from '../../components/modal.vue';
import BtnAddNew from '../../components/common/btn-add-new.vue';
import FarmAreaForm from './areas-form.vue';

export default {
  name: 'FarmAreas',
  components: {
    FarmAreaForm,
    Modal,
    BtnAddNew,
  },
  data() {
    return {
      showModal: false,
      data: {},
    };
  },
  computed: {
    ...mapGetters({
      areas: 'getAllAreas',
    }),
  },
  mounted() {
    this.fetchAreas();
  },
  methods: {
    ...mapActions([
      'fetchAreas',
    ]),
    getType(key) {
      return FindAreaType(key);
    },
    getSizeUnit(key) {
      return FindAreaSizeUnit(key);
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

i.fa.fa-plus {
  text-align: left;
  width: 30px;
}

.cards-wrapper {
  margin-top: 20px;

  .card-ui {
    margin-bottom: 20px;

    i {
      width: 30px;
    }

    .panel-heading {
      margin-bottom: 20px;
    }
  }
}

.bottom-space {
  padding-bottom: 60px;
}
</style>
