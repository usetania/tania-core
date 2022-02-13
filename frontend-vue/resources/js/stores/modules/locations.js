import NProgress from 'nprogress'

import * as types from '../mutation-types'
import { http } from '../../services'

const state = {
  countries: [],
  cities: []
}

const getters = {
  getCountries: state => state.countries,
  getCities: state => state.cities
}

const actions = {
  fetchCountries ({ commit, state }) {
    NProgress.start()

    return new Promise((resolve, reject) => {
      http.get('locations/countries', ({ data }) => {
        commit(types.FETCH_COUNTRIES, data)
        resolve(data)
      }, error => reject(error.response))
    })
  },

  fetchCitiesByCountryCode ({ commit, state }, payload) {
    NProgress.start()

    return new Promise((resolve, reject) => {
      http.get('locations/cities?country_id=' + payload, ({ data }) => {
        commit(types.FETCH_CITIES, data)
        resolve(data)
      }, error => reject(error.response))
    })
  }
}

const mutations = {
  [types.FETCH_COUNTRIES] (state, payload) {
    state.countries = payload
  },
  [types.FETCH_CITIES] (state, payload) {
    state.cities = payload
  }
}

export default {
  state, actions, getters, mutations
}
