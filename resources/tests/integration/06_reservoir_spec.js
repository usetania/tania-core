import * as Reservoir from '../factory/reservoir'

describe('Reservoir', () => {
    it ('should show the reservoir page', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
    })

    it ('should show the reservoir form', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      Reservoir.open_reservoir_form()
    })

    it ('should show the reservoir form and close it', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      Reservoir.open_reservoir_form()
      cy.get('button').contains('CANCEL').first().click()
    })

    it ('should show the reservoir form and show errors for empty form submission', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      Reservoir.open_reservoir_form()
      Reservoir.check_empty_form()
    })

    it ('should create a tap reservoir', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      Reservoir.open_reservoir_form()
      Reservoir.filled_tap_reservoir()
    })

    it ('should create a well reservoir', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      Reservoir.open_reservoir_form()
      Reservoir.filled_well_reservoir()
    })
})
