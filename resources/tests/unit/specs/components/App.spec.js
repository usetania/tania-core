import { shallowMount, createLocalVue } from '../../../../../node_modules/@vue/test-utils'
import Vuex from 'vuex'
import AppComponent from '@/components/app.vue'

const routerView = {
  name: 'router-view',
  render: h => h('div'),
};
const localVue = createLocalVue()
localVue.use(Vuex)
localVue.component('router-view', routerView)

describe('components/app', () => {
  it ('should have "app-content" class if user already authenticated', () => {
    const getters = {
      IsUserAllowSeeNavigator: state => true
    }
    const actions = {
      fetchCountries () {
        return new Promise((resolve, reject) => resolve())
      },
      fetchFarmTypes () {
        return new Promise((resolve, reject) => resolve())
      }
    }

    const VuexFactory = new Vuex.Store({getters, actions})
    const wrapper = shallowMount(AppComponent, {
      mocks: { $store: VuexFactory },
      localVue
    })

    expect(wrapper.find('#content').classes()).toContain('app-content')
  })

  it ('should not have "app-content" class if user not authenticated', () => {
    const getters = {
      IsUserAllowSeeNavigator: state => false
    }
    const actions = {
      fetchCountries () {
        return new Promise((resolve, reject) => resolve())
      },
      fetchFarmTypes () {
        return new Promise((resolve, reject) => resolve())
      }
    }

    const VuexFactory = new Vuex.Store({getters, actions})
    const wrapper = shallowMount(AppComponent, {
      mocks: { $store: VuexFactory },
      localVue
    })

    expect(wrapper.find('#content.app-content').exists()).toBe(false)
  })
})
