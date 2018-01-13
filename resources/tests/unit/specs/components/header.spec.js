import Vuex from 'vuex/dist/vuex.js'
import VeeValidate from 'vee-validate'
import Header from '@/components/header.vue'
import { shallow, createLocalVue } from 'vue-test-utils';
import farmStore from '@/stores/modules/farm'

const localVue = createLocalVue()
localVue.use(Vuex)
localVue.use(VeeValidate)

const state = {
  farm: {
    uid: "d2e5bf67-79c1-4684-b14c-766a1b92155b",
    name: "FarmName1",
    description: "Tanibox farm"
  },
  farms: [
    {
      uid: "d2e5bf67-79c1-4684-b14c-766a1b92155b",
      name: "FarmName1",
      description: "Tanibox farm"
    },
    {
      uid: "509e68ce-8c77-49b6-8c53-84d67069441f",
      name: "FarmName2",
      description: "Tanibox farm 2"
    }
  ]
}

const VuexFactory = new Vuex.Store({
  state: state,
  actions: farmStore.actions,
  getters: farmStore.getters,
  mutations: farmStore.mutations
})

describe('components/header', () => {

  it('should display current farm from state', () => {
    const wrapper = shallow(Header, {
      mocks: { $store: VuexFactory },
      localVue
    })
    expect(wrapper.find('li.dropdown.farmswitch').classes()).toContain('closed')
    expect(wrapper.find('a.farm-current span').text().trim()).toEqual('FarmName1')
  })

  it('should display dropdown farms', () => {
    const wrapper = shallow(Header, {
      mocks: { $store: VuexFactory },
      localVue
    })
    wrapper.find('a.farm-current').trigger('click')
    expect(wrapper.find('li.dropdown.farmswitch').classes()).toContain('open')
  })

  it('should not display dropdown farms', () => {
    const wrapper = shallow(Header, {
      mocks: { $store: VuexFactory },
      localVue
    })
    wrapper.find('a.farm-current').trigger('click')
    expect(wrapper.find('li.dropdown.farmswitch').classes()).toContain('open')
    wrapper.find('a.farm-current').trigger('click')
    expect(wrapper.find('li.dropdown.farmswitch').classes()).toContain('closed')
  })

  it('should change farm context when the user select a fam in dropdown list', () => {
    const wrapper = shallow(Header, {
      mocks: { $store: VuexFactory },
      localVue
    })
    // open dropdown
    wrapper.find('a.farm-current').trigger('click')
    expect(wrapper.find('li.dropdown.farmswitch').classes()).toContain('open')

    // select a farm
    wrapper.find('ul.dropdown-menu li a#FarmName2').trigger('click')
    expect(wrapper.find('a.farm-current span').text().trim()).toEqual('FarmName2')
  })
})
