import Vue from 'vue'
import VueRouter from './router'
import VuexStore from './stores'

import Component from './component'
import VeeValidate, { Validator } from 'vee-validate'

import moment from 'moment-timezone'
import VueMoment from 'vue-moment'
import vClickOutside from 'v-click-outside'

import AppComponent from './components/app.vue'
import {
  http,
  IsAlphanumSpaceHyphenUnderscore,
  IsFloat,
  IsLatitude,
  IsLongitude
} from './services'

Validator.extend('alpha_num_space', IsAlphanumSpaceHyphenUnderscore)
Validator.extend('float', IsFloat)
Validator.extend('latitude', IsLatitude)
Validator.extend('longitude', IsLongitude)

Vue.use(Component)
Vue.use(VeeValidate)
Vue.use(VueMoment, { moment })
Vue.use(vClickOutside)

new Vue({
  el: '#app',
  router: VueRouter,
  store: VuexStore,
  render: h => h (AppComponent),
  created () {
    http.init()
  }
})
