<template lang="pug">
.login-wrapper
  .container-fluid
    .row
      .col-xs-12.col-sm-12.col-md-6.offset-md-3.col-lg-4.offset-lg-4
        b-card(
          class="card-block"
        )
          .text-center
            img(
              src="../../../images/logobig.png"
              alt="Tania Logo"
              width="200"
            )
          b-form(@submit.prevent="validateBeforeSubmit")
            .form-group(:class="{ 'control': true }")
              label#label-username
                translate Username
              input.form-control#username(
                type="text"
                v-validate="'required'"
                :class="{'input': true, 'text-danger': errors.has('username') }"
                :placeholder="$gettext('Input your username here')"
                v-model="username"
                name="username"
              )
              span.help-block.text-danger(
                v-show="errors.has('username')"
              ) {{ errors.first('username') }}
            .form-group(:class="{ 'control': true }")
              label#label-password
                translate Password
              input.form-control#password(
                type="password"
                v-validate="'required'"
                :class="{'input': true, 'text-danger': errors.has('password') }"
                :placeholder="$gettext('Your password here')"
                v-model="password" name="password"
              )
              span.help-block.text-danger(
                v-show="errors.has('password')"
              ) {{ errors.first('password') }}
            .form-group.text-center.m-t
              button.btn.btn-addon.btn-primary(type="submit")
                i.fas.fa-unlock
                translate Login
</template>

<script>
import Nprogres from 'nprogress';
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'Login',

  data() {
    return {
      username: '',
      password: '',
    };
  },

  computed: {
    ...mapGetters({
      user: 'getCurrentUser',
      IsNewUser: 'IsNewUser',
    }),
  },

  mounted() {
    // redirect if the user already auntenticated
    if (this.user.uid !== '') {
      this.$router.push({ name: 'Home' });
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
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.login();
        }
      });
    },
    login() {
      this.userLogin({
        username: this.username,
        password: this.password,
        client_id: process.env.CLIENT_ID,
        response_type: 'token',
        redirect_uri: `${window.location.protocol}//${window.location.host}`,
        state: 'random-string',
      }).then(this.redirector)
        .catch(() => this.$toasted.error('Incorrect Username and/or password'));
    },
    redirector() {
      Promise.all([
        this.fetchCountries(),
        this.fetchFarm(),
        this.fetchFarmTypes(),
        this.fetchFarmInventories(),
      ]).then(() => {
        if (this.IsNewUser === true) {
          this.$router.push({ name: 'IntroFarmCreate' });
        } else {
          this.$router.push({ name: 'Home' });
        }
        Nprogres.done();
      }).catch(error => error);
    },
  },
};
</script>

<style lang="scss" scoped>
.login-wrapper {
  height: 100%;
  display: flex;
  align-items: center;
}
</style>
