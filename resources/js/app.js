import Vue from 'vue'
import VueRouter from './router'
import VuexStore from './stores'

import AppComponent from './components/app.vue'
import { http } from './services'

new Vue({
  el: '#app',
  router: VueRouter,
  store: VuexStore,
  render: h => h (AppComponent),
  created () {
    http.init()
  }
})
