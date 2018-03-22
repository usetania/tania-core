import * as Area from '../factory/area'

describe('Areas', () => {
    it ('should show the areas page', () => {
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
