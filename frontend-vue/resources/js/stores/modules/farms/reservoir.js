import NProgress from 'nprogress'

import * as types from '../../mutation-types'
import FarmApi from '../../api/farm'

const state = {
  reservoir: {},
  reservoirs: [],
  reservoirnotes: [],
}

const getters = {
  getCurrentReservoir: state => state.reservoir,
  getAllReservoirs: state => state.reservoirs
}

const actions = {
  submitReservoir ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      if (payload.uid != '') {
        FarmApi.ApiUpdateReservoir(payload.uid, payload, ({ data }) => {
          commit(types.UPDATE_RESERVOIR, data.data)
          resolve(payload)
        }, error => reject(error.response))
      } else {
        FarmApi
          .ApiCreateReservoir(farm.uid, payload, ({ data }) => {
            commit(types.CREATE_RESERVOIR, data.data)
            resolve(payload)
          }, error => reject(error.response))
      }
    })
  },
  fetchReservoirs ({ commit, state, getters }, payload) {
    const farm = getters.getCurrentFarm
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFetchReservoir(farm.uid, ({ data }) => {
        commit(types.FETCH_RESERVOIR, data.data)
        resolve(data)
      }, error => reject(error.response))
    })
  },
  getReservoirByUid ({ commit, state, getters }, reservoirId) {
    const farm = getters.getCurrentFarm

    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi.ApiFindReservoirByUid(farm.uid, reservoirId, ({ data }) => {
        resolve(data)
      }, error => reject(error.response))
    })
  },
  createReservoirNotes ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let reservoirId = payload.obj_uid
      FarmApi
        .ApiCreateReservoirNotes(reservoirId, payload, ({ data }) => {
          payload = data.data
          commit(types.CREATE_RESERVOIR_NOTES, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  deleteReservoirNote ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let reservoirId = payload.obj_uid
      let noteId = payload.uid
      FarmApi
        .ApiDeleteReservoirNotes(reservoirId, noteId, payload, ({ data }) => {
          payload = data.data
          commit(types.DELETE_RESERVOIR_NOTES, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },
}

const mutations = {
  [types.CREATE_RESERVOIR] (state, payload) {
    state.reservoirs.push(payload)
  },
  [types.UPDATE_RESERVOIR] (state, payload) {
    const reservoirs = state.reservoirs
    state.reservoirs = reservoirs.map(reservoir => (reservoir.uid === payload.uid) ? payload : reservoir)
  },
  [types.SET_RESERVOIR] (state, payload) {
    state.reservoir = payload
  },
  [types.FETCH_RESERVOIR] (state, payload) {
    state.reservoirs = payload
  },
  [types.CREATE_RESERVOIR_NOTES] (state, payload) {
    state.reservoirnotes.push(payload)
  },
  [types.DELETE_RESERVOIR_NOTES] (state, payload) {
    state.reservoirnotes.push(payload)
  }
}

export default {
  state, getters, actions, mutations
}
