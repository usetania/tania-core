import VeeValidate from 'vee-validate'
import Vuex from 'vuex'
import Login from '@/pages/auth/login.vue'
import { shallowMount, createLocalVue } from '../../../../../node_modules/@vue/test-utils'
import userStore from '@/stores/modules/user'

const localVue = createLocalVue()
localVue.use(Vuex)
localVue.use(VeeValidate)

describe('pages/auth/login', () => {

  const VuexStore = new Vuex.Store({
    state: userStore.state,
    actions: userStore.actions,
    getters: userStore.getters,
    mutations: userStore.mutations
  })

  const Router = {
    push() {

    }
  }

  it('should render username and password label', () => {
    const wrapper = shallowMount(Login, {
      mocks: {
        $store: VuexStore,
        $router: Router
      },
      localVue
    })
    expect(wrapper.find('#label-username').text().trim()).toEqual('Username')
    expect(wrapper.find('#label-password').text().trim()).toEqual('Password')
  })
})
