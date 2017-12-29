<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../images/logobig.png")

    .m-b-lg
      .wrapper
        .panel.panel-default
          .panel-body
            form(@submit.prevent="login")
              .form-group
                label Username
                input.form-control(type="text" placeholder="Please input username" v-model="username")
              .form-group
                label Password
                input.form-control(type="password" placeholder="Please input username" v-model="password")
              .form-group.text-center.m-t
                  button.btn.btn-addon.btn-primary(type="submit")
                    i.fa.fa-long-arrow-right
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
      user: 'getCurrentUser'
    })
  },

  methods: {
    ...mapActions({
      loginUser: 'login'
    }),
    login () {
      this.loginUser({
        username: this.username,
        password: this.password
      }).then(response => {
        Nprogres.done()

        // if the current user is new user
        if (this.user.intro === true) {
          this.$router.push({ name: 'IntroFarmCreate' })
        } else {
          this.$router.push({ name: 'Home' })
        }

      }).catch(error => {

      })
    }
  }
}
</script>
