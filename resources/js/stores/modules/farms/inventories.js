import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'

const state = {
  inventories: []
}

const getters = {
  getAllFarmInventories: state => state.inventories,
}

const actions = {
  fetchFarmInventories ({ commit, state }) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFetchFarmInventories(({ data }) => {
          commit(types.FETCH_FARM_INVENTORIES, data.data)
          resolve(data)
        }, error => reject(error.response))
    })
  }
}

const mutations = {
  [types.FETCH_FARM_INVENTORIES] (state, payload) {
    state.inventories = payload
  },
}

export default {
  state, getters, actions, mutations
}
