import NProgress from 'nprogress'

import * as types from '../../mutation-types'
import { AddClicked } from '../../helpers/farms/crop'
import { calculateNumberOfPages, pageLength } from '../../constants'
import FarmApi from '../../api/farm'
import moment from 'moment-timezone'

const state = {
  tasks: [],
  pages: 0,
  total: 0,
}

const getters = {
  getTasks: state => state.tasks,
  getNumberOfTasks: state => state.total,
  getTasksNumberOfPages: state => state.pages,
}

const actions = {
  fetchTasks ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFetchTask(payload.pageId, ({ data }) => {
          commit(types.FETCH_TASKS, data.data)
          commit(types.SET_PAGES, data.total_rows)
          resolve(data)
        }, error => reject(error.response))
    })
  },
  getTasksByDomainAndAssetId ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      FarmApi
        .ApiFindTasksByDomainAndAssetId(payload.pageId, payload.domain, payload.assetId, ({ data }) => {
          commit(types.FETCH_TASKS, data.data)
          commit(types.SET_PAGES, data.total_rows)
          resolve(data)
        }, error => reject(error.response))
    })
  },
  getTasksByCategoryAndPriorityAndStatus ({ commit, state }, payload) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      let query = '&'
      if (payload.status == 'COMPLETED') {
        query += 'status=COMPLETED'
      } else if (payload.status == 'THISWEEK') {
        let due_start = moment().startOf('week').format('YYYY-MM-DD')
        let due_end = moment().endOf('week').format('YYYY-MM-DD')
        query += 'due_start=' + due_start +'&due_end=' + due_end
      }  else if (payload.status == 'THISMONTH') {
        let due_start = moment().startOf('month').format('YYYY-MM-DD')
        let due_end = moment().endOf('month').format('YYYY-MM-DD')
        query += 'due_start=' + due_start +'&due_end=' + due_end
      } else if (payload.status == 'OVERDUE') {
        query += 'is_due=true'
      } else if (payload.status == 'TODAY') {
        let due = moment().format('YYYY-MM-DD')
        query += 'due_date=' + due
      } else {
        query += 'status=CREATED'
      }
      FarmApi
        .ApiFindTasksByCategoryAndPriorityAndStatus(payload.pageId, payload.category, payload.priority, query, ({ data }) => {
          commit(types.FETCH_TASKS, data.data)
          commit(types.SET_PAGES, data.total_rows)
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
    state.tasks.unshift(payload)
    if (state.tasks.length > pageLength) {
      state.tasks.pop()
    }
    state.pages = calculateNumberOfPages(state.tasks.length + 1)
  },
  [types.UPDATE_TASK] (state, payload) {
    const tasks = state.tasks
    state.tasks = tasks.map(task => (task.uid === payload.uid) ? payload : task)
    state.tasks = AddClicked(state.tasks)
  },
  [types.FETCH_TASKS] (state, payload) {
    state.tasks = AddClicked(payload)
  },
  [types.SET_PAGES] (state, payload) {
    state.total = payload
    state.pages = calculateNumberOfPages(payload)
  },
}

export default {
  state, getters, actions, mutations
}
