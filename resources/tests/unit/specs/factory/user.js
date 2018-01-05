import user from '@/stores/modules/user'

// TODO : implmenet inject loader when api is ready

export default {
  state: user.state,
  getters: user.getters,
  actions: user.actions,
  mutations: user.mutations
}
