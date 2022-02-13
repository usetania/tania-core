// setup JSDOM
require('jsdom-global')()
window.Date = Date
global.performance = window.performance;

// make sure polyfill is loaded before generators
require('@babel/polyfill')

require('chai').should()

// // make common utils available globally as well
global.expect = require('expect')
global.sinon = require('sinon')

const testUtils = require('vue-test-utils')
global.shallow = testUtils.shallow
global.mount = testUtils.mount
window.__UNIT_TESTING__ = true
