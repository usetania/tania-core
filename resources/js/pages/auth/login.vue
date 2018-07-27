<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../images/logobig.png")

    .m-b-lg
      .wrapper
        .panel.panel-default
          .panel-body
            form(@submit.prevent="validateBeforeSubmit")
              .form-group(:class="{ 'control': true }")
                label#label-username Username
                input.form-control#username(type="text" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('username') }" placeholder="Please input username" v-model="username" name="username")
                span.help-block.text-danger(v-show="errors.has('username')") {{ errors.first('username') }}
              .form-group(:class="{ 'control': true }")
                label#label-password Password
                input.form-control#password(type="password" v-validate="'required'" :class="{'input': true, 'text-danger': errors.has('password') }" placeholder="Please input username" v-model="password" name="password")
                span.help-block.text-danger(v-show="errors.has('password')") {{ errors.first('password') }}
              .form-group.text-center.m-t
                  button.btn.btn-addon.btn-primary(type="submit")
                    i.fas.fa-long-arrow-alt-right
                    | Login
</template>

<script>
import Nprogres from 'nprogress'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: 'Login',

  data () {
    return {
      username: '',
      password: ''
    }
  },

  computed : {
    ...mapGetters({
      user: 'getCurrentUser',
      IsNewUser: 'IsNewUser'
    })
  },

  mounted () {
    // redirect if the user already auntenticated
    if (this.user.uid !== '') {
      this.$router.push({ name: 'Home' })
    }
  },

  methods: {
    ...mapActions([
      'userLogin',
      'fetchCountries',
      'fetchFarmTypes',
      'fetchFarm',
      'fetchFarmInventories',
    ]),
    validateBeforeSubmit() {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.login()
        }
      })
    },
    login () {
      this.userLogin({
        username: this.username,
        password: this.password,
        client_id: process.env.CLIENT_ID,
        response_type: 'token',
        redirect_uri: location.protocol+"//"+location.host,
        state: 'random-string',
      }).then(this.redirector)
      .catch(() => this.$toasted.error('Incorrect Username and/or password'))
    },
    redirector (response) {
      Promise.all([
        this.fetchCountries(),
        this.fetchFarm(),
        this.fetchFarmTypes(),
        this.fetchFarmInventories()
      ]).then(response => {
        if (this.IsNewUser === true) {
          this.$router.push({ name: 'IntroFarmCreate' })
        } else {
          this.$router.push({ name: 'Home' })
        }
        Nprogres.done()
      }).catch(error => console.log(error))
    }
  }
}
</script>
