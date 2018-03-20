import { login } from '../factory/auth'
import { filled_material_seed } from '../factory/task'

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

    it ('should create a seed material', () => {
      login()
      cy.get('#inventories').click()
      cy.get('#materials').click()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/materials')
      })

      cy.get('#materialsform').click() 
      cy.get('form').should('be.visible') 

      filled_material_seed()
    })
})
