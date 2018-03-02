<template lang="pug">
  .account-form.col
    .wrapper-md
      h1.m-n.font-thin.h3.text-primary Account Settings
    .wrapper-md
      .col-xs-6
        .panel
          .panel-body
            form(@submit.prevent="validateBeforeSubmit")
              .form-group
                label(for="username") Username
                input.form-control#username(type="text" v-validate="'required|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('username') }" v-model="account.username" name="username")
                span.help-block.text-danger(v-show="errors.has('username')") {{ errors.first('username') }}
              .form-group
                label(for="password") Password
                input.form-control#password(type="password" v-validate="'required|confirmed|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('password') }" v-model="account.password" name="password")
                span.help-block.text-danger(v-show="errors.has('password')") {{ errors.first('password') }}
              .form-group
                label(for="password_confirmation") Confirm Password
                input.form-control#password_confirmation(type="password" :class="{'input': true, 'text-danger': errors.has('password_confirmation') }" v-model="account.password_confirmation" name="password_confirmation")
              .form-group
                button.btn.btn-addon.btn-success(type="submit") SAVE
</template>

<script>
export default {
  name: "Account",
  data () {
    return {
      account: {},
    }
  },
  methods: {
    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.submit()
        }
      })
    },
  }
}
</script>
