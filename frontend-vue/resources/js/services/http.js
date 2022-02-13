import axios from 'axios'
import NProgress from 'nprogress'
import qs from 'qs';

import { ls } from '../services'

const state = {
  token: '',
  expires_in: 0
}

export const http = {
  request (method, url, data, successCb = null, errorCb = null, headers = {}) {
    axios.request({
      url,
      data: data instanceof FormData ? data : qs.stringify(data),
      method,
      headers: Object.assign({}, {
        'Content-Type': 'application/x-www-form-urlencoded'
      }, headers)
    }).then(successCb).catch(errorCb)
  },

  get (url, successCb = null, errorCb = null) {
    return this.request('get', url, {}, successCb, errorCb)
  },

  post (url, data, successCb = null, errorCb = null, headers = {}) {
    return this.request('post', url, data, successCb, errorCb, headers)
  },

  login (url, data, successCb = null, errorCb = null, headers = {}) {
    return axios.post(url, qs.stringify(data), { headers: Object.assign({}, {
        'Content-Type': 'application/x-www-form-urlencoded'
      }, headers) }).then(function(response) {
      var url = new URL(response.request.responseURL)
      var token = url.searchParams.get("access_token")
      var expires_in = url.searchParams.get("expires_in")
      ls.set('token', token)
      ls.set('expires_in', expires_in)
      return response
    }).catch(function () {
      throw new Error()
    })
  },

  put (url, data, successCb = null, errorCb = null) {
    return this.request('put', url, data, successCb, errorCb)
  },

  delete (url, data = {}, successCb = null, errorCb = null) {
    return this.request('delete', url, data, successCb, errorCb)
  },

  /**
   * Init the service.
   */
  init () {

    axios.defaults.baseURL = '/api'

    // Intercept the request to make sure the token is injected into the header.
    axios.interceptors.request.use(config => {
      // we intercept axios request and add authorizatio header before perform send a request to the server
      config.headers.Authorization = 'Bearer '+ ls.get('token')
      return config
    })

    // Intercept the response and…
    axios.interceptors.response.use(response => {
      NProgress.done()

      // …get the token from the header or response data if exists, and save it.
      // const token = response.headers['Authorization'] || response.data['token']
      // token && ls.set('jwt-token', token)

      return response
    }, error => {
      NProgress.done()
      // Also, if we receive a Bad Request / Unauthorized error
      if (error.response.status === 400 || error.response.status === 401) {
        // and we're not trying to login
        console.log(error)
      }

      return Promise.reject(error)
    })
  }
}

export default {
  state
}
