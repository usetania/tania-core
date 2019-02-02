import 'bootstrap'
import Vue from 'vue'
import VueRouter from './router'
import VuexStore from './stores'

import { library } from 'fortawesomeDir/fontawesome-svg-core'
import { 
  faLongArrowAltRight, 
  faHome, 
  faLeaf, 
  faAngleDown, 
  faAngleRight,
  faAngleDoubleRight,
  faClipboard,
  faArchive,
  faUser,
  faChevronCircleLeft,
  faChevronCircleRight,
  faDiceFour
} from 'fortawesomeDir/free-solid-svg-icons'
import { FontAwesomeIcon } from 'fortawesomeDir/vue-fontawesome'

import Component from './component'
import VeeValidate, { Validator } from 'vee-validate'

import moment from 'moment-timezone'
import VueMoment from 'vue-moment'
import vClickOutside from 'v-click-outside'

import Toasted from 'vue-toasted'

import AppComponent from './components/app.vue'
import {
  http,
  IsAlphanumSpaceHyphenUnderscore,
  IsFloat,
  IsLatitude,
  IsLongitude
} from './services'

library.add(
  faLongArrowAltRight, 
  faHome, 
  faLeaf, 
  faAngleDown, 
  faAngleRight, 
  faClipboard, 
  faArchive, 
  faUser, 
  faChevronCircleRight, 
  faChevronCircleLeft, 
  faDiceFour, 
  faAngleDoubleRight
)

Vue.component('font-awesome-icon', FontAwesomeIcon)

Validator.extend('alpha_num_space', IsAlphanumSpaceHyphenUnderscore)
Validator.extend('float', IsFloat)
Validator.extend('latitude', IsLatitude)
Validator.extend('longitude', IsLongitude)

Vue.use(Component)
Vue.use(VeeValidate)
Vue.use(VueMoment, { moment })
Vue.use(vClickOutside)
Vue.use(Toasted, { 
   theme: "bubble", 
   position: "bottom-center", 
   duration : 1500
})

new Vue({
  el: '#app',
  router: VueRouter,
  store: VuexStore,
  render: h => h (AppComponent),
  created () {
    http.init()
  }
})
