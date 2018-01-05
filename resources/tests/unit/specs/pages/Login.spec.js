import Vue from 'vue/dist/vue.common.js'
import VeeValidate from 'vee-validate'
import Login from '@/pages/auth/login.vue'

Vue.use(VeeValidate)

describe('pages/auth/login', () => {
  it('should render username label', () => {
    const Constructor = Vue.extend(Login);
    const vm = new Constructor().$mount();
    expect(vm.$el.querySelector('#label-username').textContent)
    .toEqual('Username');
  })

  it('should render password label', () => {
    const Constructor = Vue.extend(Login);
    const vm = new Constructor().$mount();
    expect(vm.$el.querySelector('#label-username').textContent)
    .toEqual('Username');
  })
})
