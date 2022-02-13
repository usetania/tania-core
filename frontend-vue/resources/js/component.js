import ErrorMessage from './components/error.vue';

export default {
  ErrorMessage,
  install(Vue) {
    Vue.component('error-message', ErrorMessage);
  },
};
