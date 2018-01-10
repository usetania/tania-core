import Vue from 'vue'
import VueRouter from './router'
import VuexStore from './stores'

import Component from './component'
import VeeValidate from 'vee-validate'

import moment from 'moment-timezone'
import VueMoment from 'vue-moment'

import AppComponent from './components/app.vue'
import { http } from './services'

Vue.use(Component)
Vue.use(VeeValidate)
Vue.use(VueMoment, { moment })

new Vue({
  el: '#app',
  router: VueRouter,
  store: VuexStore,
  render: h => h (AppComponent),
  created () {
    http.init()
  }
})
