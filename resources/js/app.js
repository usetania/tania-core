import Vue from 'vue'
import VueRouter from './router'
import VuexStore from './stores'
import Component from './component'
import VeeValidate, { Validator } from 'vee-validate'
import moment from 'moment-timezone'
import VueMoment from 'vue-moment'
import Toasted from 'vue-toasted'
import AppComponent from './components/app.vue'
import GetTextPlugin from 'vue-gettext'
import translations from '../../languages/translations.json'

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
Vue.use(Toasted, {
   theme: "bubble",
   position: "bottom-center",
   duration : 1500
})

Vue.use(GetTextPlugin, {
  availableLanguages: {
    en_GB: 'British English',
    id_ID: 'Bahasa Indonesia',
    hu_HU: 'Magyar Nyelv',
    pt_BR: 'Brazilian Portuguese'
  },
  defaultLanguage: 'en_GB',
  translations: translations,
  silent: false
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
