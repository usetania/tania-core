import NProgress from 'nprogress'

import * as types from '@/stores/mutation-types'
import API from '@/stores/api/farm'
import stub from '@/stores/stubs/user'

const state = {
  current: stub
}

const getters = {
  getCurrentUser: state => state.current,
  IsUserAuthenticated: state => state.current.uid != '',
  IsNewUser: (state, getters) => getters.haveFarms === false,
  IsUserAllowSeeNavigator: (state, getters) => {
    return getters.IsUserAuthenticated && getters.IsNewUser === false
  }
}

const actions = {
  userLogin ({ commit, state }, payload) {
    NProgress.start()
    return new Promise(( resolve, reject ) => {
      API
        .ApiLogin(payload, ({ data }) => {
          commit(types.USER_LOGIN, {
            uid: data.uid,
            username: payload.username,
            email: 'hello@tanibox.com',
            intro: payload.username === 'user' ? false: true
          })
          resolve(data)
        }, error => reject(error.response))
    })
  },
  userChangePassword ({ commit, state }, payload) {
    NProgress.start()
    return new Promise(( resolve, reject ) => {
      API
        .ApiChangePassword(payload, ({ data }) => {
          resolve(data)
        }, error => reject(error.response))
    })
  },
  userCompletedIntro({ commit, state }) {
    commit(types.USER_COMPLETED_INTRO)
  },
  userSignOut({commit, state}, payload) {
    return new Promise((resolve, reject) => {
      commit(types.USER_LOGOUT)
      resolve()
    })
  }
}

const mutations = {
  [types.USER_LOGIN] (state, { uid, username, email, intro }) {
    state.current = { uid, username, email, intro }
  },
  [types.USER_COMPLETED_INTRO] (state) {
    state.current.intro = false
  },
  [types.USER_LOGOUT] (state, payload) {
    state.current.uid = ''
  }
}

export default {
  state, getters, actions, mutations
}
