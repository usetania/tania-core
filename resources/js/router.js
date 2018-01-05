import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/stores'

Vue.use(VueRouter)

const router = new VueRouter({
  routes: [
    { path: '/', name: 'Home', meta: { requiresAuth: true }, component: () => import('./pages/home.vue') },

    // Authenticated
    { path: '/auth/login', name: 'AuthLogin', component: () => import('./pages/auth/login.vue') },

    // intro
    { path: '/intro/farm', name: 'IntroFarmCreate', meta: { requiresAuth: true }, component: () => import('./pages/farms/intro/farm.vue') },
    { path: '/intro/reservoir', name: 'IntroReservoirCreate', meta: { requiresAuth: true }, component: () => import('./pages/farms/intro/reservoir.vue') },
    { path: '/intro/area', name: 'IntroAreaCreate', meta: { requiresAuth: true }, component: () => import('./pages/farms/intro/area.vue') },

    // Farm
    { path: '/farm-add', name: 'FarmCreate', meta: { requiresAuth: true }, component: () => import('./pages/farms/create.vue') }
  ]
})


router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (store.getters.IsUserAuthenticated === false) {
      next({
        name: 'AuthLogin',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }

})

export default router
