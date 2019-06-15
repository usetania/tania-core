<template lang="pug">
  .container-fluid.container-intro
    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        .text-center
          img(
            src="../../../images/logobig.png"
            alt="Tania Logo"
            width="200"
          )

    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        h3.text-center
          translate Hello! Can you tell me a little about your farm?

    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        b-form(@submit.prevent="validateBeforeSubmit")
          b-card(
            :title="$gettext('Create Farm')"
          )
            error-message(:message="message.error_message")
            .form-group
              label#label-name(for="name")
                translate Farm Name
              input.form-control#name(
                type="text"
                v-validate="'required|alpha_num_space|min:5|max:100'"
                :class="{'input': true, 'text-danger': errors.has('name') }"
                placeholder=""
                v-model="farm.name"
                name="name"
              )
              span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
            .form-group
              label#label-description(for="description")
                translate Farm Description
              textarea.form-control#description(
                placeholder=""
                rows="3"
                v-model="farm.description"
              )
            .form-group
              BtnContinue(:title="$gettext('Continue')" customClass="float-right")

</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { StubFarm, StubMessage } from '../../stores/stubs';
import BtnContinue from '../../components/common/btn-continue.vue';

export default {
  components: {
    BtnContinue,
  },

  data() {
    return {
      message: Object.assign({}, StubMessage),
      farm: Object.assign({}, StubFarm),
      mapbox: {},
    };
  },

  computed: {
    ...mapGetters({
      currentFarm: 'introGetFarm',
    }),
  },

  mounted() {
    if (this.currentFarm) {
      this.farm = Object.assign({}, this.currentFarm);
    }
  },

  methods: {
    ...mapActions([
      'introSetFarm',
    ]),

    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.create();
        }
      });
    },

    create() {
      this.introSetFarm(this.farm);
      this.$router.push({ name: 'IntroReservoirCreate' });
    },
  },
};
</script>


<style lang="scss" scoped>
.container-intro {
  background-color: #f6f8f8;
  padding-top: 20px;
  padding-bottom: 20px;
}
</style>
