import Vue from 'vue/dist/vue.common.js'
import Vuex from 'vuex/dist/vuex.js'

Vue.use(Vuex)

import user from './user'
import farm from './farm'

export default new Vuex.Store({
  modules: {
    user, farm
  }
})
