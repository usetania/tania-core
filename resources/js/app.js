import Vue from 'vue'
import VueRouter from './router'
import VuexStore from './stores'

import Component from './component'
import VeeValidate from 'vee-validate'

import AppComponent from './components/app.vue'
import { http } from './services'

Vue.use(Component)
Vue.use(VeeValidate)

new Vue({
  el: '#app',
  router: VueRouter,
  store: VuexStore,
  render: h => h (AppComponent),
  created () {
    http.init()
  }
})
