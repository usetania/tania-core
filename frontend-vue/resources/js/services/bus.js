import Vue from 'vue'

class EventBus {

  constructor () {
    this.bus = new Vue()
  }

  emit (name, ...args) {
    this.bus.$emit(name, ...args)
  }

  on () {
    if (arguments.length === 2) {
      this.bus.$on(arguments[0], arguments[1])
    } else {
      Object.keys(arguments[0]).forEach(key => this.bus.$on(key, arguments[0][key]))
    }

    return this
  }
}


const event = new EventBus()

export { event }
