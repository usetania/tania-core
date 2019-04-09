<template lang="pug">
  .container-fluid
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
          translate Awesome! Now let's create a new water source for your farm.

    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        b-card(
          :title="$gettext('Create Reservoir')"
        )
          error-message(:message="message.error_message")
          b-form(@submit.prevent="validateBeforeSubmit")
            .form-group
              label#label-name(for="name")
                translate Reservoir Name
              input.form-control#name(
                type="text"
                v-validate="'required|alpha_num_space|min:5|max:100'"
                :class="{'input': true, 'text-danger': errors.has('name') }"
                v-model="reservoir.name"
                name="name"
              )
              span.help-block.text-danger(v-show="errors.has('name')") {{ errors.first('name') }}
            .form-group
              label#label-source(for="type")
                translate Source
              select.form-control#type(
                v-validate="'required'"
                :class="{'input': true, 'text-danger': errors.has('type') }"
                v-model="reservoir.type"
                name="type"
                @change="typeChanged($event.target.value)"
              )
                option(value="")
                  translate Please select source
                option(v-for="option in options" :value="option.key") {{ option.label }}
              span.help-block.text-danger(v-show="errors.has('type')") {{ errors.first('type') }}
            .form-group(v-if="reservoir.type == 'BUCKET'")
              input#capacity.form-control(
                type="text"
                v-validate="'required'"
                :class="{'input': true, 'text-danger': errors.has('capacity') }"
                v-model="reservoir.capacity"
                :placeholder="$gettext('Capacity (litre)')"
                name="capacity"
              )
              span.help-block.text-danger(v-show="errors.has('capacity')")
                | {{ errors.first('capacity') }}
            .form-group
              button.btn.btn-addon.btn-primary.float-right(type="submit")
                i.fas.fa-long-arrow-alt-right
                translate Continue
              router-link#back.btn.btn-addon.btn-default(:to="{name: 'IntroFarmCreate'}")
                i.fas.fa-long-arrow-alt-left
                translate Back
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import { ReservoirTypes } from '../../stores/helpers/farms/reservoir';
import { StubReservoir, StubMessage } from '../../stores/stubs';

export default {
  name: 'ReservoirIntro',
  data() {
    return {
      message: Object.assign({}, StubMessage),
      reservoir: Object.assign({}, StubReservoir),
      options: Array.from(ReservoirTypes),
    };
  },
  computed: {
    ...mapGetters({
      currentReservoir: 'introGetReservoir',
      currentFarm: 'introGetFarm',
    }),
  },

  mounted() {
    if (this.currentReservoir) {
      this.reservoir = Object.assign({}, this.currentReservoir);
    }

    if (this.currentFarm.name === '') {
      this.$router.push({ name: 'IntroFarmCreate' });
    }
  },

  methods: {
    ...mapActions([
      'introSetReservoir',
    ]),
    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.create();
        }
      });
    },
    typeChanged(type) {
      if (type === 'BUCKET') {
        this.reservoir.capacity = '';
      } else {
        this.reservoir.capacity = 0;
      }
    },
    create() {
      this.introSetReservoir(this.reservoir);
      this.$router.push({ name: 'IntroAreaCreate' });
    },
  },
};
</script>
