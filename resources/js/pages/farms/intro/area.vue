<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../../images/logobig.png")
    h3.text-lt.text-center.wrapper.m-t Perfect! Let's create a new area.
    .m-b-lg
      .wrapper
        .panel.panel-default
          .panel-heading
            h4.text-lt Add New Area
          .panel-body
            p.text-muted
              | Area is a space where you grow your plants. It could be a seeding tray, a garden bed, or a
              | pot or anything that describes the different physical locations in your facility.
            form(@submit.prevent="validateBeforeSubmit")
              .line.line-dashed.b-b.line-lg
              .row
                .col-xs-6
                  .form-group
                    label(for="name") Area Name
                    input.form-control#name(type="text" v-validate="'required|alpha_num|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('area.name') }" v-model="area.name" name="area.name")
                    span.help-block.text-danger(v-show="errors.has('area.name')") {{ errors.first('area.name') }}
                .col-xs-6
                  .from-group
                    label Size
                    .row
                      .col-xs-6
                        input.form-control#size(type="text" v-validate="'required|numeric'" :class="{'input': true, 'text-danger': errors.has('area.size') }" v-model="area.size" name="area.size")
                        span.help-block.text-danger(v-show="errors.has('area.size')") {{ errors.first('area.size') }}
                      .col-xs-6
                        select.form-control(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('area.unit') }" v-model="area.unit" name="area.unit")
                          option hectare
                          option meter square
                        span.help-block.text-danger(v-show="errors.has('area.unit')") {{ errors.first('area.unit') }}
              .row
                .col-xs-6
                  .form-group
                    label(for="type") Type
                    select.form-control#type(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('area.type') }" v-model="area.type" name="area.type")
                      option Nursery / Seeding
                      option Growing Area
                    span.help-block.text-danger(v-show="errors.has('area.type')") {{ errors.first('area.type') }}
                .col-xs-6
                  .form-group
                    label(for="locations") Locations
                    select.form-control#locations(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('area.location') }" v-model="area.location" name="area.location")
                      option Field (Outdoor)
                      option Greenhouse (Indoor)
                    span.help-block.text-danger(v-show="errors.has('area.location')") {{ errors.first('area.location') }}
              .row
                .col-xs-6
                  .form-group
                    label Select Reservoir
                    select.form-control(v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('area.reservoir') }" v-model="area.reservoir_id" name="area.reservoir")
                      option Please select reservoir
                      option(v-for="reservoir in reservoirs") {{ reservoir.name }}
                    span.help-block.text-danger(v-show="errors.has('area.reservoir')") {{ errors.first('area.reservoir') }}
                .col-xs-6
                  .form-group
                    label Select photo <small class="text-muted">(if any)</small>
                    input(type="file")
              .form-group
                button.btn.btn-addon.btn-success.pull-right(type="submit")
                  | Finish Setup
                  i.fa.fa-long-arrow-right
                router-link.btn.btn-addon.btn-default(:to="{name: 'IntroReservoirCreate'}")
                  i.fa.fa-long-arrow-left
                  | Back

</template>

<script>
import stub from '@/stores/stubs/area'
import stubMessage from '@/stores/stubs/message'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: 'AreaIntro',
  data () {
    return {
      message: Object.assign({}, stubMessage),
      area: Object.assign({}, stub)
    }
  },
  computed: {
    ...mapGetters({
      reservoirs: 'getAllReservoirs'
    })
  },
  methods: {
    ...mapActions([
      'userCompletedIntro',
      'createArea'
    ]),
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },
    create () {
      this.createArea(this.area)
        .then(data => {
          this.userCompletedIntro()
          this.$router.push({ name: 'Home' })
        })
        .catch(({ data }) => this.message = data)
    }
  }
}
</script>

