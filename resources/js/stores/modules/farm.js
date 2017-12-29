import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import { http } from '@/services'
import stub from '@/stores/stubs/farm'

const state = {
  current: stub,
  all: [],
  types: []
}

const getters = {
  getCurrentFarm: state => state.current,
  getAllFarm: state => state.all,
  getAllFarmTypes: state => state.types,
  haveFarms: state => state.all.length > 0 ? true : false
}

const actions = {
  createFarm ({ commit, state }, payload) {
    NProgress.start()

    return new Promise((resolve, reject) => {
      http.post('farms', payload, ({ data }) => {
        resolve(data)
      }, error => reject(error))
    })
  },

  fetchFarmTypes ({ commit, state }) {
    NProgress.start()

    return new Promise((resolve, reject) => {
      http.get('farms/types', ({ data }) => {
        commit(types.FETCH_FARM_TYPES, data)
        resolve(data)
      }, error => reject(error))
    })
  }
}

const mutations = {
  [types.FETCH_FARM_TYPES] (state, payload) {
    state.types = payload
  }
}

export default {
  state, getters, actions, mutations
}
