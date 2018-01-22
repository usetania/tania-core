import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'

const state = {
  crops: [],
  cropnotes: [],
}

const getters = {
  getAllCrops: state => state.crops
}

const actions = {

  createCrop ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      let areaId = payload.initial_area
      FarmApi
        .ApiCreateCrop(areaId, payload, ({ data }) => {
          payload = data.data
          payload.farm_id = farm.uid
          commit(types.CREATE_CROP, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  fetchCrops ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchCrop(farm.uid, ({ data }) => {
        commit(types.FETCH_CROP, data.data)
        resolve(data)
      }, error => reject(error.response))
    })
  },
  getCropByUid ({ commit, state, getters }, cropId) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFindCropByUid(farm.uid, cropId, ({ data }) => {
        resolve(data)
      }, error => reject(error.response))
    })
  },
  createCropNotes ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let cropId = payload.obj_uid
      FarmApi
        .ApiCreateCropNotes(cropId, payload, ({ data }) => {
          payload = data.data
          commit(types.CREATE_CROP_NOTES, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  }
}

const mutations = {
  [types.CREATE_CROP] (state, payload) {
    state.crops.push(payload)
  },
  [types.SET_CROP] (state, payload) {
    state.crops = payload
  },
  [types.FETCH_CROP] (state, payload) {
    state.crops = payload
  },
  [types.CREATE_CROP_NOTES] (state, payload) {
    state.cropnotes.push(payload)
  }
}

export default {
  state, getters, actions, mutations
}
