import NProgress from 'nprogress'

import * as types from '../../mutation-types'
import { calculateNumberOfPages, pageLength } from '../../constants'
import FarmApi from '../../api/farm'

const state = {
  inventories: [],
  materials: [],
  fertilizers: [],
  pesticides: [],
  pages: 0,
}

const getters = {
  getAllFarmInventories: state => state.inventories,
  getAllMaterials: state => state.materials,
  getMaterialsNumberOfPages: state => state.pages,
  getAllPesticides: state => state.pesticides,
  getAllFertilizers: state => state.fertilizers,
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
  fetchMaterials ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFetchMaterial(payload.pageId, ({ data }) => {
          commit(types.FETCH_MATERIALS, data.data)
          commit(types.SET_PAGES, data.total)
          resolve(data)
        }, error => reject(error.response))
    })
  },
  fetchAgrochemicalMaterials ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFetchAgrochemicalMaterial(payload.type, ({ data }) => {
          if (payload.type == 'FERTILIZER') {
            commit(types.FETCH_FERTILIZERS, data.data)
          } else if (payload.type == 'PESTICIDE') {
            commit(types.FETCH_PESTICIDES, data.data)
          }
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
    state.materials.unshift(payload)
    if (state.materials.length > pageLength) {
      state.materials.pop()
    }
    state.pages = calculateNumberOfPages(state.materials.length + 1)
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
  [types.SET_PAGES] (state, pages) {
    state.pages = calculateNumberOfPages(pages)
  },
  [types.FETCH_FERTILIZERS] (state, payload) {
    state.fertilizers = payload
  },
  [types.FETCH_PESTICIDES] (state, payload) {
    state.pesticides = payload
  },
}

export default {
  state, getters, actions, mutations
}
