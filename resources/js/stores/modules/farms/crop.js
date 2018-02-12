import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'

const state = {
  crops: [],
  cropnotes: [],
  cropactivities: [],
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
  },
  deleteCropNote ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let cropId = payload.obj_uid
      let noteId = payload.uid
      FarmApi
        .ApiDeleteCropNotes(cropId, noteId, payload, ({ data }) => {
          payload = data.data
          commit(types.DELETE_CROP_NOTES, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  moveCrop ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let cropId = payload.obj_uid
      FarmApi
        .ApiMoveCrop(cropId, payload, ({ data }) => {
          payload = data.data
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  harvestCrop ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let cropId = payload.obj_uid
      FarmApi
        .ApiHarvestCrop(cropId, payload, ({ data }) => {
          payload = data.data
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  dumpCrop ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let cropId = payload.obj_uid
      FarmApi
        .ApiDumpCrop(cropId, payload, ({ data }) => {
          payload = data.data
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  photoCrop ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let cropId = payload.obj_uid
      const formData = new FormData()
      formData.set('photo', payload.photo)
      formData.set('description', payload.description)
      FarmApi
        .ApiPhotoCrop(cropId, formData, ({ data }) => {
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  fetchActivities ({ commit, state, getters }, cropId) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchActivity(cropId, ({ data }) => {
        resolve(data)
      }, error => reject(error.response))
    })
  },
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
  },
  [types.DELETE_CROP_NOTES] (state, payload) {
    state.cropnotes.push(payload)
  },
  [types.FETCH_CROP_ACTIVITIES] (state, payload) {
    state.cropactivities = payload
  },
  [types.CREATE_CROP_ACTIVITIES] (state, payload) {
    state.cropactivities.push(payload)
  },
}

export default {
  state, getters, actions, mutations
}
