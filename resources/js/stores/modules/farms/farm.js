import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'
import stubFarm from '@/stores/stubs/farm'
import stubReservoir from '@/stores/stubs/reservoir'
import stubCrop from '@/stores/stubs/crop'
import stub from '@/stores/stubs/farm'

const state = {
  farm: stub,
  farms: [],
  types: [],
}

const getters = {
  getCurrentFarm: state => state.farm,
  getAllFarms: state => state.farms,
  getAllFarmTypes: state => state.types,
  haveFarms: state => state.farms.length > 0 ? true : false
}

const actions = {
  createFarm ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiCreateFarm(payload, ({ data }) => {
          commit(types.CREATE_FARM, data.data)
          commit(types.SET_FARM, data.data)
          resolve(data.data)
        }, error => reject(error.response))
    })
  },
  fetchFarm ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFetchFarm(({ data }) => {
          commit(types.FETCH_FARM, data.data)

          // select the current farm for the first array
          if (data.data.length > 0) {
            commit(types.SET_FARM, data.data[0])
          }

          resolve(data)
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
  setCurrentFarm({ commit, state }, farmId) {
    return new Promise((resolve, reject) => {
      let farm = state.farms.find(item => item.uid === farmId)
      if (farm) {
        commit(types.SET_FARM, farm)
        resolve(farm)
      } else {
        reject()
      }
    }, error => reject(error))
  },
}

const mutations = {
  [types.FETCH_FARM] (state, payload) {
    state.farms = payload
  },
  [types.FETCH_FARM_TYPES] (state, payload) {
    state.types = payload
  },
  [types.CREATE_FARM] (state, payload) {
    state.farms.push(payload)
  },
  [types.SET_FARM] (state, payload) {
    state.farm = payload
  }
}

export default {
  state, getters, actions, mutations
}
