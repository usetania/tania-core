<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../images/logobig.png")
    h3.text-lt.text-center.wrapper.m-t Awesome! Now let's create a new<br/> water source for your farm.
    .m-b-lg
      .wrapper
        .panel.panel-default
          .panel-heading
            h4.text-lt Create Reservoir
          .panel-body
            error-message(:message="message.error_message")
            form(@submit.prevent="validateBeforeSubmit")
              .form-group
                label(for="name") Reservoir Name
                input.form-control#name(type="text" v-validate="'required|alpha_num|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('reservoir.name') }" v-model="reservoir.name" name="reservoir.name")
                span.help-block.text-danger(v-show="errors.has('reservoir.name')") {{ errors.first('reservoir.name') }}
              .form-group
                label(for="type") Source
                select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('reservoir.type') }" v-model="reservoir.type" name="reservoir.type" @change="typeChanged($event.target.value)")
                  option(value="") Please select source
                  option(v-for="option in options" :value="option.key") {{ option.label }}
                span.help-block.text-danger(v-show="errors.has('reservoir.type')") {{ errors.first('reservoir.type') }}
              .form-group(v-if="reservoir.type == 'bucket'")
                input.form-control(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('reservoir.capacity') }" v-model="reservoir.capacity" placeholder="Capacity (litre)" name="reservoir.capacity")
                span.help-block.text-danger(v-show="errors.has('reservoir.capacity')") {{ errors.first('reservoir.capacity') }}
              .form-group
                button.btn.btn-addon.btn-primary.pull-right(type="submit")
                  | Continue
                  i.fa.fa-long-arrow-right
                router-link.btn.btn-addon.btn-default(:to="{name: 'IntroFarmCreate'}")
                  i.fa.fa-long-arrow-left
                  | Back
</template>


<script>
import { ReservoirTypes } from '@/stores/helpers/farms/reservoir'
import { StubReservoir, StubMessage } from '@/stores/stubs'
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'ReservoirIntro',
  data () {
    return {
      message: Object.assign({}, StubMessage),
      reservoir: Object.assign({}, StubReservoir),
      options: Array.from(ReservoirTypes)
    }
  },

  computed: {
    ...mapGetters({
      currentReservoir: 'introGetReservoir'
    })
  },

  mounted () {
    if (this.currentReservoir) {
      this.reservoir = Object.assign({}, this.currentReservoir)
    }
  },

  methods: {
    ...mapActions([
      'introSetReservoir'
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    typeChanged (type) {
      if (type === 'bucket') {
        this.reservoir.capacity = ''
      } else {
        this.reservoir.capacity = 0
      }
    },
    create () {
      this.introSetReservoir(this.reservoir)
      this.$router.push({ name: 'IntroAreaCreate' })
    }
  }
}
</script>
