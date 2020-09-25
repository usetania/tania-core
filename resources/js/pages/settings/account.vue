<template lang="pug">
  .container-fluid.bottom-space
    .row
      .col
        h3.title-page
          translate Account Settings

    .row
      .col-xs-12.col-sm-12.col-md-6
        b-card
          b-form(@submit.prevent="validateBeforeSubmit")
            .form-group
              label(for="username")
                translate Username
              input.form-control#username(type="text" v-validate="'required|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('username') }" v-model="user.username" name="username" readonly="true")
              span.help-block.text-danger(v-show="errors.has('username')") {{ errors.first('username') }}
            .form-group
              label(for="old_password")
                translate Old Password
              input.form-control#old_password(type="password" v-validate="'required|min:5|max:100'" :class="{'input': true, 'text-danger': errors.has('old password') }" v-model="user.old_password" name="old password" placeholder="")
              span.help-block.text-danger(v-show="errors.has('old password')") {{ errors.first('old password') }}
            .form-group
              label(for="password")
                translate Password
              input.form-control#new_password(type="password" v-validate="'required|min:5|max:100'" ref="password" :class="{'input': true, 'text-danger': errors.has('password') }" v-model="user.password" name="password" placeholder="")
              span.help-block.text-danger(v-show="errors.has('password')") {{ errors.first('password') }}
            .form-group
              label(for="password_confirmation")
                translate Confirm Password
              input.form-control#confirm_new_password(type="password" v-validate="'required|confirmed:password'" data-vv-as="password confirmation" :class="{'input': true, 'text-danger': errors.has('password_confirmation') }" v-model="user.password_confirmation" name="password_confirmation" placeholder="")
              span.help-block.text-danger(v-show="errors.has('password_confirmation')") {{ errors.first('password_confirmation') }}
            .form-group
              BtnSave
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { StubUser } from '../../stores/stubs';
import BtnSave from '../../components/common/btn-save.vue';

export default {
  name: 'Account',
  components: {
    BtnSave,
  },
  data() {
    return {
      user: Object.assign({}, StubUser),
    };
  },
  computed: {
    ...mapGetters({
      current_user: 'getCurrentUser',
    }),
  },
  mounted() {
    this.user.username = this.current_user.username;
    this.user.uid = this.current_user.uid;
    const dict = {
      custom: { password_confirmation: { confirmed: 'The password confirmation does not match.' } }
    }
    this.$validator.localize('en', dict);
  },
  methods: {
    ...mapActions([
      'userChangePassword',
    ]),
    submit() {
      this.userChangePassword({
        userid: this.user.uid,
        old_password: this.user.old_password,
        new_password: this.user.password,
        confirm_new_password: this.user.password_confirmation,
      }).then(() => {
        this.$toasted.show('Password update successful');
        this.user.old_password = '';
        this.user.password = '';
        this.user.password_confirmation = '';
        this.$nextTick(() => this.$validator.reset());
      }).catch(() => this.$toasted.error('Error in user password update'));
    },
    validateBeforeSubmit() {
      this.$validator.validateAll().then((result) => {
        if (result) {
          this.submit();
        }
      }).catch(() => {
        console.log(this.$validator.errors);
      });
    },
  },
};
</script>

<style lang="scss" scoped>
h3.title-page {
  margin: 20px 0 30px 0;
}

.bottom-space {
  padding-bottom: 60px;
}
</style>
