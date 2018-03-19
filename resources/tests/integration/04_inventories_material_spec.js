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

    it ('should show the materials page and open the materials form modal', () => {
      login()
      cy.get('#inventories').click()
      cy.get('#materials').click()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/materials')
      })

      cy.get('#materialsform').click() 
      cy.get('form').should('be.visible') 
    })
})
