export const IsAlphanumSpaceHyphenUnderscore = {
  getMessage(field, params, data) {
    return (data && data.message) || 'The ' + field + ' should be alphanumeric, space, hypen, or underscore'
  },
  validate (value) {
    let regexp = /^[a-zA-Z0-9]+[?:[\w -]*[a-zA-Z0-9]$/
    return value.search(regexp) === 0
  }
}

export const IsFloat = {
  getMessage(field, params, data) {
    return (data && data.message) || 'The ' + field + ' should be float only'
  },
  validate (value) {
    let regexp = /^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$/
    return regexp.test(value)
  }
}

export const IsLatitude = {
  getMessage(field, params, data) {
    return (data && data.message) || 'The ' + field + ' is not latitude value'
  },
  validate (value) {
    let regexp = /^(?=.)-?((8[0-5]?)|([0-7]?[0-9]))?(?:\.[0-9]{1,20})?$/
    return regexp.test(value)
  }
}

export const IsLongitude = {
  getMessage(field, params, data) {
    return (data && data.message) || 'The ' + field + ' is not longitude value'
  },
  validate (value) {
    let regexp = /^(?=.)-?((0?[8-9][0-9])|180|([0-1]?[0-7]?[0-9]))?(?:\.[0-9]{1,20})?$/
    return regexp.test(value)
  }
}
