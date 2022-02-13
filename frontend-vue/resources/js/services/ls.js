import Store from 'local-storage'

export const ls = {
  get (key, defaultVal = null) {
    return Store(key) || defaultVal
  },

  set (key, val) {
    return Store(key, val)
  },

  remove (key) {
    return Store.remove(key)
  }
}
