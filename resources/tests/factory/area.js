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

export function check_empty_form() {
  // Label
  cy.get('label#label-name').should('contain', 'Area Name')
  cy.get('label#label-size').should('contain', 'Size')
  cy.get('label#label-type').should('contain', 'Type')
  cy.get('label#label-location').should('contain', 'Locations')
  cy.get('label#label-reservoir').should('contain', 'Select Reservoir')

  cy.get('button[type=submit]').click()
  cy.get('span.help-block.text-danger').should('contain', 'The name field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The size field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The reservoir field is required.')
}

export function filled_seeding_area() {
  // Label
  cy.get('label#label-name').should('contain', 'Area Name')
  cy.get('label#label-size').should('contain', 'Size')
  cy.get('label#label-type').should('contain', 'Type')
  cy.get('label#label-location').should('contain', 'Locations')
  cy.get('label#label-reservoir').should('contain', 'Select Reservoir')

  // Typing the form
  cy.get('#name').type('Seeding Area '+moment().valueOf())
  cy.get('#size').type('4').should('have.value', '4')
  cy.get('#size_unit').select('Ha')
  cy.get('#type').select('SEEDING')
  cy.get('#location').select('INDOOR')
  cy.get('#reservoir').select('Reservoir1')
  cy.get('button[type=submit]').click()
  cy.get('form').should('not.be.visible')
}

export function filled_growing_area() {
  // Label
  cy.get('label#label-name').should('contain', 'Area Name')
  cy.get('label#label-size').should('contain', 'Size')
  cy.get('label#label-type').should('contain', 'Type')
  cy.get('label#label-location').should('contain', 'Locations')
  cy.get('label#label-reservoir').should('contain', 'Select Reservoir')

  // Typing the form
  cy.get('#name').type('Growing Area '+moment().valueOf())
  cy.get('#size').type('3').should('have.value', '3')
  cy.get('#size_unit').select('Ha')
  cy.get('#type').select('GROWING')
  cy.get('#location').select('OUTDOOR')
  cy.get('#reservoir').select('Reservoir1')
  cy.get('button[type=submit]').click()
  cy.get('form').should('not.be.visible')
}
