import moment from 'moment';
import { login } from './auth'

export function open_task_page() {
  login()
  cy.get('#tasks').click()
  cy.location().should( location => {
    expect(location.hash).to.eq('#/tasks')
  })
}

export function open_task_form() {
  cy.get('#tasksform').click() 
  cy.get('form').should('be.visible') 
}

export function check_empty_form() {
  cy.get('button[type=submit]').click()

  cy.get('span.help-block.text-danger').should('contain', 'The due date field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The priority field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The category field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The title field is required.')
}

export function filled_task(type) {
  // Label
  cy.get('label#label-due-date').should('contain', 'Due Date')
  cy.get('label#label-priority').should('contain', 'Is this task urgent?')
  cy.get('label#label-category').should('contain', 'Task Category')
  cy.get('label#label-title').should('contain', 'Title')
  cy.get('label#label-description').should('contain', 'Description')

  // Typing the form
  cy.get('input#due_date').click()
  cy.get('span.next').first().click()
  cy.get('span').contains(15).first().click()

  cy.get('label').contains('Yes').first().click()

  cy.get('select#category').select(type)
  cy.get('input#title').type(type + ' Task '+moment().valueOf())
  cy.get('textarea#description').type(type + ' Task Description'+moment().valueOf())
  cy.get('button[type=submit]').click()
  cy.get('form').should('not.be.visible')
}
