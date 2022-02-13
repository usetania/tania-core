import * as Area from '../factory/area'
import * as Task from '../factory/task'

describe('Areas', () => {
  describe('Areas Listing Page', () => {
    it ('should show the areas listing page', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
    })

    it ('should show the areas form', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      Area.open_area_form()
    })

    it ('should show the areas form and close it', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      Area.open_area_form()
      cy.get('button').contains('CANCEL').first().click()
    })

    it ('should show the areas form and show errors for empty form submission', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      Area.open_area_form()
      Area.check_empty_form()
    })

    it ('should create a seeding area', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      Area.open_area_form()
      Area.filled_seeding_area()
    })

    it ('should create a growing area', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      Area.open_area_form()
      Area.filled_growing_area()
    })
  })

  describe('Area Page', () => {
    it ('should show the area page', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      cy.get('a').contains('Area1').first().click()
    })

    it ('should show the area page and create a note', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      cy.get('a').contains('Area1').first().click()
      cy.get('input#content').type('This is the note to area 1')
      cy.get('button[type=submit]').click()
    })

    it ('should delete the newly created note', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      cy.get('a').contains('Area1').first().click()
      cy.get('input#content').type('This is the second note to area 1')
      cy.get('button[type=submit]').click()
      cy.get('button').children('i').should('have.class', 'fa-trash').first().click()
    })

    it ('should create a task for the area', () => {
      cy.clearLocalStorage()
      Area.open_area_page()
      cy.get('a').contains('Area1').first().click()
      cy.get('#addTaskForm').click()
      Task.filled_task('AREA')
    })
  })
})
