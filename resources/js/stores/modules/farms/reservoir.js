import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'
import stubFarm from '@/stores/stubs/farm'
import stubReservoir from '@/stores/stubs/reservoir'
import stubCrop from '@/stores/stubs/crop'
import stub from '@/stores/stubs/farm'

const state = {
  reservoir: {},
  reservoirs: []
}

const getters = {
  getCurrentReservoir: state => state.reservoir,
  getAllReservoirs: state => state.reservoirs
}

const actions = {

  createReservoir ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiCreateReservoir(farm.uid, payload, ({ data }) => {
          payload = data.data
          payload.farm_id = farm.uid
          commit(types.CREATE_RESERVOIR, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  fetchReservoirs ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchReservoir(farmId, ({ data }) => {
        commit(types.FETCH_RESERVOIR, data.data)
        resolve(data)
      }, error => reject(error.response))
    })
  }
}

const mutations = {
  [types.CREATE_RESERVOIR] (state, payload) {
    state.reservoirs.push(payload)
  },
  [types.SET_RESERVOIR] (state, payload) {
    state.reservoir = payload
  },
  [types.FETCH_RESERVOIR] (state, payload) {
    state.reservoirs = payload
  }
}

export default {
  state, getters, actions, mutations
}
