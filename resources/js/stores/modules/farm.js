import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'
import stub from '@/stores/stubs/farm'

const state = {
  current: stub,
  farms: [],
  reservoirs: [],
  types: []
}

const getters = {
  getCurrentFarm: state => state.current,
  getAllFarms: state => state.farms,
  getAllReservoirs: state => state.reservoirs,
  getAllFarmTypes: state => state.types,
  haveFarms: state => state.farms.length > 0 ? true : false
}

const actions = {
  createFarm ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiCreateFarm(payload, ({ data }) => {
          payload.id = data.data
          commit(types.CREATE_FARM, payload)
          commit(types.SET_FARM, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },

  fetchFarmTypes ({ commit, state }) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiGetFarmTypes(({ data }) => {
          commit(types.FETCH_FARM_TYPES, data)
          resolve(data)
        }, error => reject(error.response))
    })
  },

  createReservoir ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      const farmId = state.current.id
      FarmApi
        .ApiCreateReservoir(farmId, payload, ({ data }) => {
          payload.id = data
          payload.farm_id = farmId
          commit(types.CREATE_RESERVOIR, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  }
}

const mutations = {
  [types.FETCH_FARM_TYPES] (state, payload) {
    state.types = payload
  },
  [types.CREATE_FARM] (state, payload) {
    state.farms.push(payload)
  },
  [types.SET_FARM] (state, payload) {
    state.current = payload
  },
  [types.CREATE_RESERVOIR] (state, payload) {
    state.reservoirs.push(payload)
  }
}

export default {
  state, getters, actions, mutations
}
