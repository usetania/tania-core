import VeeValidate from 'vee-validate'
import Vuex from 'vuex'
import Login from '../../../../js/pages/auth/login.vue'
import { shallowMount, createLocalVue } from '../../../../../node_modules/vue-test-utils'
import userStore from '../../../../js/stores/modules/user'
import GetTextPlugin from 'vue-gettext'
import translations from '../../../../../languages/translations.json'

const localVue = createLocalVue()
localVue.use(Vuex)
localVue.use(VeeValidate)
localVue.use(GetTextPlugin, {
  availableLanguages: {
    en_GB: 'British English',
    id_ID: 'Bahasa Indonesia',
    hu_HU: 'Magyar Nyelv'
  },
  defaultLanguage: 'en_GB',
  translations: translations,
  silent: false
})

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

  it('should render input fields', () => {
    const wrapper = shallow(Login, {
      mocks: {
        $store: VuexStore,
        $router: Router
      },
      localVue
    })
    expect(wrapper.contains('input')).toBe(true)
    expect(wrapper.contains('input')).toBe(true)
  })
})
