import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'

const state = {
  inventories: [],
  materials: [],
}

const getters = {
  getAllFarmInventories: state => state.inventories,
  getAllMaterials: state => state.materials,
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
  },
  fetchMaterials ({ commit, state }) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFetchMaterial(({ data }) => {
          commit(types.FETCH_MATERIALS, data.data)
          resolve(data)
        }, error => reject(error.response))
    })
  },
  submitMaterial ({ commit, state, getters }, payload) {
    NProgress.start()
    if (payload.uid != '') {
      return new Promise((resolve, reject) => {
        FarmApi
          .ApiUpdateMaterial(payload.uid, payload, ({ data }) => {
            payload = data.data
            commit(types.UPDATE_MATERIAL, payload)
            resolve(payload)
          }, error => reject(error.response))
      })
    } else {
      return new Promise((resolve, reject) => {
        FarmApi
          .ApiCreateMaterial(payload, ({ data }) => {
            payload = data.data
            commit(types.CREATE_MATERIAL, payload)
            resolve(payload)
          }, error => reject(error.response))
      })
    }
  },
}

const mutations = {
  [types.CREATE_MATERIAL] (state, payload) {
    state.materials.push(payload)
  },
  [types.UPDATE_MATERIAL] (state, payload) {
    const materials = state.materials
    state.materials = materials.map(material => (material.uid === payload.uid) ? payload : material)
  },
  [types.FETCH_MATERIALS] (state, payload) {
    state.materials = payload
  },
  [types.FETCH_FARM_INVENTORIES] (state, payload) {
    state.inventories = payload
  },
}

export default {
  state, getters, actions, mutations
}
