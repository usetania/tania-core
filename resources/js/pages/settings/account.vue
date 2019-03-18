<template lang="pug">
  .user-form.col
    .wrapper-md
      h1.m-n.font-thin.h3.text-primary Account Settings
    .wrapper-md
      .col-xs-6
        .panel
          .panel-body
            form(@submit.prevent="validateBeforeSubmit")
              .form-group
                label(for="username") Username
                input.form-control#username(type="text" v-validate="'required|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('username') }" v-model="user.username" name="username" readonly="true")
                span.help-block.text-danger(v-show="errors.has('username')") {{ errors.first('username') }}
              .form-group
                label(for="old_password") Old Password
                input.form-control#old_password(type="password" v-validate="'required|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('old password') }" v-model="user.old_password" name="old password" placeholder="")
                span.help-block.text-danger(v-show="errors.has('old password')") {{ errors.first('old password') }}
              .form-group
                label(for="password") Password
                input.form-control#new_password(type="password" v-validate="'required|confirmed|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('password') }" v-model="user.password" name="password" placeholder="")
                span.help-block.text-danger(v-show="errors.has('password')") {{ errors.first('password') }}
              .form-group
                label(for="password_confirmation") Confirm Password
                input.form-control#confirm_new_password(type="password" :class="{'input': true, 'text-danger': errors.has('password_confirmation') }" v-model="user.password_confirmation" name="password_confirmation" placeholder="")
                span.help-block.text-danger(v-show="errors.has('password_confirmation')") {{ errors.first('password_confirmation') }}
              .form-group
                button.btn.btn-addon.btn-success(type="submit") SAVE
</template>

<script>
import { StubUser } from '../../stores/stubs'
import { mapActions, mapGetters } from 'vuex'
export default {
  name: "Account",
  computed : {
    ...mapGetters({
      current_user: 'getCurrentUser',
    })
  },
  data () {
    return {
      user: Object.assign({}, StubUser),
    }
  },
  methods: {
    ...mapActions([
      'userChangePassword',
    ]),
    submit () {
      this.userChangePassword({
        userid : this.user.uid,
        old_password : this.user.old_password,
        new_password : this.user.password,
        confirm_new_password : this.user.password_confirmation,
        }).then(() => {
          this.$toasted.show('Password update successful')
          this.user.old_password = ''
          this.user.password = ''
          this.user.password_confirmation = ''
          this.$nextTick(() => this.$validator.reset())
        }).catch(() => this.$toasted.error('Error in user password update'))
    },
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.submit()
        }
      }).catch(() => {
        console.log(this.$validator.errors)
      })
    },
  },
  mounted () {
    this.user.username = this.current_user.username
    this.user.uid = this.current_user.uid
  }
}
</script>
