import localStorage from 'local-storage'
import { ls } from '../../../../js/services'

describe('services/ls', () => {
  beforeEach(() => localStorage.remove('foo'))

  describe('#get', () => {
    it('correctly gets an existing item from local storage', () => {
      localStorage('foo', 'bar')
      expect(ls.get('foo')).toEqual('bar')
    })

    it('correctly returns the default value for a non exising item', () => {
      let baz = ls.get('baz', 'qux')
      expect(baz).toEqual('qux')
    })
  })

  describe('#set', () => {
    it('correctly sets an item into local storage', () => {
      ls.set('foo', 'bar')
      expect(localStorage('foo')).toEqual('bar')
    })
  })

  describe('#remove', () => {
    it('correctly removes an item from local storage', () => {
      localStorage('foo', 'bar')
      ls.remove('foo')
      expect(localStorage('foo')).toEqual(null)
    })
  })
})
