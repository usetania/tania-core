import NProgress from 'nprogress'
import Api from '../api/farm'
import * as types from '../mutation-types'
import stubFarm from '../stubs/farm'
import stubReservoir from '../stubs/reservoir'
import stubArea from '../stubs/area'

const defaults = {
  farm: Object.assign({}, stubFarm),
  reservoir: Object.assign({}, stubReservoir),
  area: Object.assign({}, stubArea)
}

const state = Object.assign({}, defaults)

const getters = {
  introGetFarm : state => state.farm,
  introGetReservoir: state => state.reservoir,
  introGetArea: state => state.area,
  introGetUserPosition: (state, getters) => {
    if (getters.introGetFarm.name === '') {
      return 'IntroFarmCreate'
    } else if (getters.introGetReservoir.name === '') {
      return 'IntroReservoirCreate'
    } else if (getters.introGetArea.name === '') {
      return 'IntroAreaCreate'
    } else {
      return 'IntroAreaCreate'
    }
  }
}

const actions = {
  introSetFarm ({ commit, state }, payload) {
    commit(types.INTRO_SET_FARM, payload)
  },
  introSetReservoir ({ commit, state }, payload) {
    commit(types.INTRO_SET_RESERVOIR, payload)
  },
  introSetArea ({ commit, state }, payload) {
    commit(types.INTRO_SET_AREA, payload)
  },
  introCreateFarm ({ commit, state }) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      if (state.farm.uid !== '') {
        resolve(state.farm)
      } else {
        Api.ApiCreateFarm(state.farm, ({ data }) => {
          let farm = data.data
          commit(types.INTRO_SET_FARM, farm)
          resolve(farm)
        }, err => reject(err.response))
      }
    })
  },
  introCreateReservoir ({ commit, state }) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      if (state.reservoir.uid !== '') {
        resolve(state.reservoir)
      } else {
        Api.ApiCreateReservoir(state.farm.uid, state.reservoir, ({ data }) => {
          let reservoir = data.data
          let area = Object.assign({}, state.area, {reservoir_id: reservoir.uid, farm_id: state.farm.uid})
          commit(types.INTRO_SET_RESERVOIR, reservoir)
          commit(types.INTRO_SET_AREA, area)
          resolve(reservoir)
        }, err => reject(err.response))
      }
    })
  },
  introCreateArea ({ commit, state }) {
    NProgress.start()
    return new Promise((resolve, reject) => {
      if (state.area.uid !== '') {
        resolve(state.area)
      } else {
        var formData = new FormData()
        formData.append('name', state.area.name)
        formData.append('size', state.area.size)
        formData.append('size_unit', state.area.size_unit)
        formData.append('type', state.area.type)
        formData.append('location', state.area.location)
        formData.append('reservoir_id', state.area.reservoir_id)
        formData.append('photo', state.area.photo)

        Api.ApiCreateArea(state.farm.uid, formData, ({ data }) => {
          let area = {
            ...data.data,
            photo: '/api/farms/' + state.farm.uid + '/areas/' + data.data.uid + '/photos'
          }
          // COMMIT
          commit(types.CREATE_FARM, state.farm)
          commit(types.SET_FARM, state.farm)
          commit(types.CREATE_RESERVOIR, state.reservoir)
          commit(types.SET_RESERVOIR, state.reservoir)
          commit(types.CREATE_AREA, area)
          commit(types.SET_AREA, area)

          commit(types.USER_COMPLETED_INTRO)

          // reset intro
          commit(types.INTRO_SET_FARM, Object.assign({}, defaults.farm))
          commit(types.INTRO_SET_RESERVOIR, Object.assign({}, defaults.reservoir))
          commit(types.INTRO_SET_AREA, Object.assign({}, defaults.area))

          // resolve
          resolve(area)
        }, err => reject(err.response))
      }
    })
  },
}

const mutations = {
  [types.INTRO_SET_FARM] (state, payload) {
    state.farm = payload
  },
  [types.INTRO_SET_RESERVOIR] (state, payload) {
    state.reservoir = payload
  },
  [types.INTRO_SET_AREA] (state, payload) {
    state.area = payload
  }
}

export default {
  state, getters, actions, mutations
}
