import farm from '@/stores/modules/farm'

const actionsInjector = require('inject-loader!@/stores/modules/farm')
const injector = actionsInjector({
  '../api/farm': {
    ApiCreateFarm: function() {
      return new Promise((resolve, reject) => {
        resolve()
      })
    },
    ApiGetFarmTypes: function() {
      return new Promise((resolve, reject) => {
        const data = [
          { Code: 'organic', Name: 'Organic / Soil-Based' },
          { Code: 'hydroponic', Name: 'Hydroponic' }
        ]
        resolve(data)
      })
    },

  }
})

export default {
  state: farm.state,
  getters: farm.getters,
  actions: injector.actions,
  mutations: farm.mutations
}
