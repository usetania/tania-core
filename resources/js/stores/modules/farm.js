import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'
import stubFarm from '@/stores/stubs/farm'
import stubReservoir from '@/stores/stubs/reservoir'
import stub from '@/stores/stubs/farm'

const state = {
  farm: stub,
  farms: [],
  reservoir: {},
  reservoirs: [],
  area: {},
  areas: [],
  types: []
}

const getters = {
  getCurrentFarm: state => state.farm,
  getCurrentReservoir: state => state.reservoir,
  getCurrentArea: state => state.area,
  getAllFarms: state => state.farms,
  getAllReservoirs: state => state.reservoirs,
  getAllAreas: state => state.areas,
  getAllFarmTypes: state => state.types,
  haveFarms: state => state.farms.length > 0 ? true : false
}

const actions = {
  createFarm ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiCreateFarm(payload, ({ data }) => {
          payload.ui = data.data
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
      const farmId = state.farm.uid
      FarmApi
        .ApiCreateReservoir(farmId, payload, ({ data }) => {
          payload.uid = data
          payload.farm_id = farmId
          commit(types.CREATE_RESERVOIR, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  fetchReservoirs ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      const farmId = state.farm.uid
      FarmApi.ApiFetchReservoir(farmId, ({ data }) => {
        commit(types.FETCH_RESERVOIR, data.data)
        resolve(data)
      }, error => reject(error.response))
    })
  },
  createArea ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchReservoir(farmId, ({ data }) => {
        commit(types.CREATE_AREA, payload)
        resolve(payload)
      }, error => reject(error.response))
    })
  },
  fetchAreas ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      const farmId = state.farm.uid
      FarmApi.ApiFetchArea(farmId, ({ data }) => {
        commit(types.FETCH_AREA, data.data)
        resolve(data)
      }, error => reject(error.response))
    })
  },
  getAreaByUid ({ commit, state }, areaId) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      const farmId = state.farm.uid
      FarmApi.ApiFindAreaByUid(farmId, areaId, ({ data }) => {
        resolve(data)
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
    state.farm = payload
  },
  [types.CREATE_RESERVOIR] (state, payload) {
    state.reservoirs.push(payload)
  },
  [types.SET_RESERVOIR] (state, payload) {
    state.reservoir = payload
  },
  [types.FETCH_RESERVOIR] (state, payload) {
    state.reservoirs = payload
  },
  [types.CREATE_AREA] (state, payload) {
    state.areas.push(payload)
  },
  [types.SET_AREA] (state, payload) {
    state.area = payload
  },
  [types.FETCH_AREA] (state, payload) {
    state.areas = payload
  },
}

export default {
  state, getters, actions, mutations
}
