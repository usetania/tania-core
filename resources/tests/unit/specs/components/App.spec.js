import Vue from 'vue/dist/vue.common.js'
import AppComponent from '@/components/app.vue'
import VuexFactory from '../factory'

describe('components/app', () => {
  it ('should load initial rendering', done => {

    // mock component
    const routerView = {
      name: 'router-view',
      render: h => h('div'),
    };

    // register mock component
    Vue.component('router-view', routerView);

    const vm = new Vue({
      template: '<AppComponent></AppComponent>',
      store: VuexFactory,
      components: {
        AppComponent
      }
    }).$mount()

    Vue.nextTick()
      .then(() => {
        expect(vm.$el.querySelector('#footer')).toBe(null)
        done()
      })
      .catch(done)

  })
})
