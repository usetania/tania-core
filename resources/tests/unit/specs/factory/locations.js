import locations from '@/stores/module/locations'

const actionsInjector = require('inject-loader!@/stores/module/locations')
const injector = actionsInjector({

})


export default {
  state: locations.state,
  getters: locations.getters,
  actions: locations.actions,
  mutations: locations.mutations
}
