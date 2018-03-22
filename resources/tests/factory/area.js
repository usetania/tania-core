import moment from 'moment';
import { login } from './auth'

export function open_area_page() {
  login()
  cy.get('#production').click()
  cy.get('#areas').click()
  cy.location().should( location => {
    expect(location.hash).to.eq('#/areas')
  })
}

export function open_area_form() {
  cy.get('#areasform').click() 
  cy.get('form').should('be.visible') 
}

export function filled_seeding_area() {
  // Label
  cy.get('label#label-name').should('contain', 'Area Name')
  cy.get('label#label-size').should('contain', 'Size')
  cy.get('label#label-type').should('contain', 'Type')
  cy.get('label#label-location').should('contain', 'Locations')
  cy.get('label#label-reservoir').should('contain', 'Select Reservoir')

  // Typing the form
  cy.get('#name').type('Seeding Area '+moment().format())
  cy.get('#size').type('4').should('have.value', '4')
  cy.get('#size_unit').select('Ha')
  cy.get('#type').select('SEEDING')
  cy.get('#location').select('INDOOR')
  cy.get('#reservoir').select('Reservoir1')
  cy.get('button.btn[type=submit]').click()
}

export function filled_growing_area() {
  // Label
  cy.get('label#label-name').should('contain', 'Area Name')
  cy.get('label#label-size').should('contain', 'Size')
  cy.get('label#label-type').should('contain', 'Type')
  cy.get('label#label-location').should('contain', 'Locations')
  cy.get('label#label-reservoir').should('contain', 'Select Reservoir')

  // Typing the form
  cy.get('#name').type('Growing Area '+moment().format())
  cy.get('#size').type('3').should('have.value', '3')
  cy.get('#size_unit').select('Ha')
  cy.get('#type').select('GROWING')
  cy.get('#location').select('OUTDOOR')
  cy.get('#reservoir').select('Reservoir1')
  cy.get('button.btn[type=submit]').click()
}
