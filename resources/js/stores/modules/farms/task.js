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
  getTasksByCategoryAndPriorityAndStatus ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFindTasksByCategoryAndPriorityAndStatus(payload.category, payload.priority, payload.status, ({ data }) => {
          resolve(data)
        }, error => reject(error.response))
    })
  },
  submitTask ({ commit, state, getters }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      if (payload.uid != '') {
        FarmApi
          .ApiUpdateTask(payload.uid, payload, ({ data }) => {
            payload = data.data
            commit(types.UPDATE_TASK, payload)
            resolve(payload)
          }, error => reject(error.response))
      } else {
        FarmApi
          .ApiCreateTask(payload, ({ data }) => {
            payload = data.data
            commit(types.CREATE_TASK, payload)
            resolve(payload)
          }, error => reject(error.response))
      }
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
  [types.UPDATE_TASK] (state, payload) {
    const tasks = state.tasks
    state.tasks = tasks.map(task => (task.uid === payload.uid) ? payload : task)
  },
  [types.FETCH_TASKS] (state, payload) {
    state.tasks = payload
  },
}

export default {
  state, getters, actions, mutations
}
