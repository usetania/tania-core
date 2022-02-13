import NProgress from 'nprogress'

import * as types from '../../mutation-types'
import { calculateNumberOfPages, pageLength } from '../../constants'
import FarmApi from '../../api/farm'

const state = {
  crops: [],
  cropnotes: [],
  cropactivities: [],
  cropPages: 0,
}

const getters = {
  getAllCrops: state => state.crops,
  getCropsNumberOfPages: state => state.cropPages,
}

const actions = {

  submitCrop ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm
    NProgress.start()
    return new Promise((resolve, reject) => {
      if (payload.uid != '') {
        FarmApi.ApiUpdateCrop(payload.uid, payload, ({ data }) => {
          commit(types.UPDATE_CROP, data.data)
          resolve(payload)
        }, error => reject(error.response))
      } else {
        let areaId = payload.initial_area
        FarmApi
          .ApiCreateCrop(areaId, payload, ({ data }) => {
            commit(types.CREATE_CROP, data.data)
            resolve(payload)
        }, error => reject(error.response))
      }
    })
  },
  fetchCrops ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchCrop(farm.uid, payload.pageId, payload.status, ({ data }) => {
        commit(types.FETCH_CROP, data.data)
        commit(types.SET_PAGES, data.total_rows)
        resolve(data)
      }, error => reject(error.response))
    })
  },
  fetchArchivedCrops ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchArchivedCrop(farm.uid, payload.pageId, ({ data }) => {
        commit(types.FETCH_CROP, data.data)
        commit(types.SET_PAGES, data.total_rows)
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
      var formData = new FormData()
      formData.append('photo', payload.photo)
      formData.append('description', payload.description)
      FarmApi
        .ApiPhotoCrop(cropId, formData, ({ data }) => {
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  waterCrop ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let cropId = payload.obj_uid
      FarmApi
        .ApiWaterCrop(cropId, payload, ({ data }) => {
          payload = data.data
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
    state.crops.unshift(payload)
    if (state.crops.length > pageLength) {
      state.crops.pop()
    }
    state.cropPages = calculateNumberOfPages(state.crops.length + 1)
  },
  [types.UPDATE_CROP] (state, payload) {
    const crops = state.crops
    state.crops = crops.map(crop => (crop.uid === payload.uid) ? payload : crop)
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
  [types.SET_PAGES] (state, payload) {
    state.cropPages = calculateNumberOfPages(payload)
  },
}

export default {
  state, getters, actions, mutations
}
