<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../images/logobig.png")
    h3.text-lt.text-center.wrapper.m-t
      translate Awesome! Now let's create a new water source for your farm.
    .m-b-lg
      .wrapper
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
