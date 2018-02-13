import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import FarmApi from '@/stores/api/farm'

const state = {
  tasks: [],
}

const getters = {
  getAllTasks: state => state.tasks,
}

const actions = {
  fetchTasks ({ commit, state }) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFetchTask(({ data }) => {
          commit(types.FETCH_TASKS, data.data)
          resolve(data)
        }, error => reject(error.response))
    })
  },
  getTasksByDomainAndAssetId ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFindTasksByDomainAndAssetId(payload.domain, payload.assetId, ({ data }) => {
          resolve(data)
        }, error => reject(error.response))
    })
  },
  createTask ({ commit, state, getters }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiCreateTask(payload, ({ data }) => {
          payload = data.data
          commit(types.CREATE_TASK, payload)
          resolve(payload)
        }, error => reject(error.response))
    })
  },
  setTaskDue ({ commit, state, getters }, taskId) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiSetTaskDue(taskId, ({ data }) => {
          resolve(data)
        }, error => reject(error.response))
    })
  },
  setTaskCompleted ({ commit, state, getters }, taskId) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiSetTaskCompleted(taskId, ({ data }) => {
          resolve(data)
        }, error => reject(error.response))
    })
  },
}

const mutations = {
  [types.CREATE_TASK] (state, payload) {
    state.tasks.push(payload)
  },
  [types.FETCH_TASKS] (state, payload) {
    state.tasks = payload
  },
}

export default {
  state, getters, actions, mutations
}
