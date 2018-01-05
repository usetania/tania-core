import NProgress from 'nprogress'
import * as types from '@/stores/mutation-types'

const state = {
  areas: []
}

const getters = {
  getAllAreas: state => state.areas
}

const actions = {
  createArea ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      commit(types.CREATE_RESERVOIR, payload)
      resolve(payload)
    })
  }
}

const mutations = {
  [types.CREATE_RESERVOIR] (state, payload) {
    state.areas.push(payload)
  }
}

export default {
  state, getters, actions, mutations
}
