import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue';
import VeeValidate, { Validator } from 'vee-validate';
import moment from 'moment-timezone';
import VueMoment from 'vue-moment';
import Toasted from 'vue-toasted';
import GetTextPlugin from 'vue-gettext';
import AppComponent from './components/app.vue';
import VueRouter from './router';
import VuexStore from './stores';
import Component from './component';
import {
  http,
  IsAlphanumSpaceHyphenUnderscore,
  IsFloat,
  IsLatitude,
  IsLongitude,
} from './services';
import translations from '../../languages/translations.json';
import '../sass/app.scss';

Validator.extend('alpha_num_space', IsAlphanumSpaceHyphenUnderscore);
Validator.extend('float', IsFloat);
Validator.extend('latitude', IsLatitude);
Validator.extend('longitude', IsLongitude);

Vue.use(Component);
Vue.use(BootstrapVue);
Vue.use(VeeValidate);
Vue.use(VueMoment, { moment });
Vue.use(Toasted, {
  theme: 'bubble',
  position: 'bottom-center',
  duration: 1500,
});

Vue.use(GetTextPlugin, {
  availableLanguages: {
    en_GB: 'British English',
    id_ID: 'Bahasa Indonesia',
    hu_HU: 'Magyar Nyelv',
    pt_BR: 'Brazilian Portuguese',
  },
  defaultLanguage: 'en_GB',
  translations,
  silent: false,
});

const App = () => new Vue({
  el: '#app',
  router: VueRouter,
  store: VuexStore,
  created() {
    http.init();
  },
  render: h => h(AppComponent),
});

App();
