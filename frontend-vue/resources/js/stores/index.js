import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from 'vuex-persistedstate'
import createLogger from 'vuex/dist/logger'

import farm from '../stores/modules/farm'
import intro from '../stores/modules/intro'
import locations from '../stores/modules/locations'
import user from '../stores/modules/user'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  modules: {
    farm, intro, locations, user
  },
  strict: debug,
  plugins: debug ? [createLogger(), createPersistedState()] : [createPersistedState()]
})
