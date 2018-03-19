import { login } from '../factory/auth'
import { filled_farm, filled_reservoir, filled_area } from '../factory/farm'

describe('Inventories', () => {
    it ('should show the materials page', () => {
      login()
      cy.get('#inventories').click()
      cy.get('#materials').click()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/materials')
      })
    })
})
