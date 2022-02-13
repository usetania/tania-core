import Vue from 'vue';
import VueRouter from 'vue-router';
import store from './stores';

Vue.use(VueRouter);

const router = new VueRouter({
  routes: [
    {
      path: '/',
      name: 'Home',
      meta: { requiresAuth: true },
      component: () => import('./pages/home.vue'),
    },

    // Authenticated
    {
      path: '/auth/login',
      name: 'AuthLogin',
      component: () => import('./pages/auth/login.vue'),
    },

    // intro
    {
      path: '/intro/farm',
      name: 'IntroFarmCreate',
      meta: { requiresAuth: true },
      component: () => import('./pages/intro/farm.vue'),
    },
    {
      path: '/intro/reservoir',
      name: 'IntroReservoirCreate',
      meta: { requiresAuth: true },
      component: () => import('./pages/intro/reservoir.vue'),
    },
    {
      path: '/intro/area',
      name: 'IntroAreaCreate',
      meta: { requiresAuth: true },
      component: () => import('./pages/intro/area.vue'),
    },

    // Farm
    {
      path: '/farms/create',
      name: 'FarmCreate',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/create.vue'),
    },
    {
      path: '/reservoirs',
      name: 'FarmReservoirs',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/reservoirs.vue'),
    },
    {
      path: '/reservoirs/:id',
      name: 'FarmReservoir',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/reservoir.vue'),
    },
    {
      path: '/areas',
      name: 'FarmAreas',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/areas.vue'),
    },
    {
      path: '/areas/:id',
      name: 'FarmArea',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/area.vue'),
    },
    {
      path: '/crops',
      name: 'FarmCrops',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/crops.vue'),
    },
    {
      path: '/crop/:id',
      name: 'FarmCrop',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/crop.vue'),
    },
    {
      path: '/crop/notes/:id',
      name: 'FarmCropNotes',
      meta: { requiresAuth: true },
      component: () => import('./pages/farms/crop-notes.vue'),
    },
    {
      path: '/tasks',
      name: 'Task',
      meta: { requiresAuth: true },
      component: () => import('./pages/tasks/task.vue'),
    },

    // Inventories
    {
      path: '/materials',
      name: 'InventoriesMaterials',
      meta: { requiresAuth: true },
      component: () => import('./pages/inventories/materials.vue'),
    },

    // Settings
    {
      path: '/settings/account',
      name: 'Account',
      meta: { requiresAuth: true },
      component: () => import('./pages/settings/account.vue'),
    },
  ],
});

function middleware(to, from, next) {
  // if user have property intro == true a.k.a new user we will redirect into intro pages
  if (store.getters.IsNewUser === true) {
    const positioning = store.getters.introGetUserPosition;
    if (to.name === positioning) {
      next();
    } else {
      const introMaps = [
        { from: 'IntroReservoirCreate', to: 'IntroFarmCreate' }, // back button
        { from: 'IntroAreaCreate', to: 'IntroReservoirCreate' }, // back button
        { from: 'IntroFarmCreate', to: 'IntroReservoirCreate' }, // next button
        { from: 'IntroReservoirCreate', to: 'IntroAreaCreate' }, // next button
      ].filter(item => from.name === item.from && to.name === item.to);

      // check if the route is available in intro maps
      if (introMaps.length === 1) {
        next();
      } else {
        next({ name: store.getters.introGetUserPosition });
      }
    }
  } else {
    next();
  }
}

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (store.getters.IsUserAuthenticated === false) {
      next({ name: 'AuthLogin', query: { redirect: to.fullPath } });
    } else {
      middleware(to, from, next);
    }
  } else {
    next();
  }
});

export default router;
