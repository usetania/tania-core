import { configure, addParameters } from '@storybook/vue';
import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue';
import GetTextPlugin from 'vue-gettext';
import translations from '../languages/translations.json';

import '../resources/sass/app.scss';

// Storybook viewport addons for responsive development
addParameters({ viewport: { defaultViewport: 'responsive' } });

Vue.use(BootstrapVue);
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

function loadStories() {
  const req = require.context('../resources/js/components', true, /\.stories\.js$/);
  req.keys().forEach(filename => req(filename));
}

configure(loadStories, module);
