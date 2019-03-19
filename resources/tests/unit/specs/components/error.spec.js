import Vue from 'vue/dist/vue.common.js'
import ErrorMessage from '../../../../js/components/error.vue'
import { shallowMount } from '../../../../../node_modules/vue-test-utils';

describe('components/error', () => {
  it('should render error message', () => {
    const wrapper = shallow(ErrorMessage, {
      propsData: {
        message: 'This field is required'
      }
    })

    expect(wrapper.props().message).toBe('This field is required')
    expect(wrapper.classes()).toContain('alert')
    expect(wrapper.text()).toBe('This field is required')
  })
  it('should not render error message', () => {
    const wrapper = shallow(ErrorMessage, {
      propsData: {
        message: ''
      }
    })

    expect(wrapper.props().message).toBe('')
    expect(wrapper.text()).toBe('')
  })
})
