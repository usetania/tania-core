import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './registerServiceWorker';
import '../scss/app.scss';

import Default from './layouts/Default.vue';
import NoSidebar from './layouts/NoSidebar.vue';

Vue.component('default-layout', Default);
Vue.component('no-sidebar-layout', NoSidebar);

Vue.config.productionTip = false;

Vue.use(BootstrapVue);

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
