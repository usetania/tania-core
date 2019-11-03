import Vue from 'vue';
import Router from 'vue-router';
import Dashboard from './views/dashboard/Index.vue';
import Installer from './views/installer/Index.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: Dashboard,
    },
    {
      path: '/installer',
      name: 'installer',
      meta: { layout: 'no-sidebar' },
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: Installer,
    },
  ],
});
