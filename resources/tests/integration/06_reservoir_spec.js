import * as Reservoir from '../factory/reservoir'
import * as Task from '../factory/task'

describe('Reservoirs', () => {
  describe('Reservoirs Listing Page', () => {
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

  describe('Reservoir Page', () => {
    it ('should show the reservoir page', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      cy.get('u').contains('Reservoir1').first().click()
    })

    it ('should show the reservoir page and create a note', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      cy.get('u').contains('Reservoir1').first().click()
      cy.get('input#content').type('This is the note to reservoir 1')
      cy.get('button[type=submit]').click()
    })

    it ('should delete the newly created note', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      cy.get('u').contains('Reservoir1').first().click()
      cy.get('input#content').type('This is the second note to reservoir 1')
      cy.get('button[type=submit]').click()
      cy.get('button').children('i').should('have.class', 'fa-trash').first().click()
    })

    it ('should create a task for the reservoir', () => {
      cy.clearLocalStorage()
      Reservoir.open_reservoir_page()
      cy.get('u').contains('Reservoir1').first().click()
      cy.get('#addTaskForm').click()
      Task.filled_task('RESERVOIR')
    })
  })
})
