import moment from 'moment';
import { login } from './auth'

export function open_reservoir_page() {
  login()
  cy.get('#production').click()
  cy.get('#reservoirs').click()
  cy.location().should( location => {
    expect(location.hash).to.eq('#/reservoirs')
  })
}

export function open_reservoir_form() {
  cy.get('#reservoirsform').click() 
  cy.get('form').should('be.visible') 
}

export function check_empty_form() {
  // Label
  cy.get('label#label-name').should('contain', 'Reservoir Name')
  cy.get('label#label-type').should('contain', 'Source')

  cy.get('button[type=submit]').click()
  cy.get('span.help-block.text-danger').should('contain', 'The name field is required.')
}

export function filled_tap_reservoir() {
  // Label
  cy.get('label#label-name').should('contain', 'Reservoir Name')
  cy.get('label#label-type').should('contain', 'Source')

  // Typing the form
  cy.get('#name').type('Tap Reservoir '+moment().valueOf())
  cy.get('#type').select('TAP')
  cy.get('#capacity').should('not.be.visible')
  cy.get('button[type=submit]').click()
  cy.get('form').should('not.be.visible')
}

export function filled_well_reservoir() {
  // Label
  cy.get('label#label-name').should('contain', 'Reservoir Name')
  cy.get('label#label-type').should('contain', 'Source')

  // Typing the form
  cy.get('#name').type('Bucket Reservoir '+moment().valueOf())
  cy.get('#type').select('BUCKET')
  cy.get('#capacity').type('10')
  cy.get('button[type=submit]').click()
  cy.get('form').should('not.be.visible')
}
