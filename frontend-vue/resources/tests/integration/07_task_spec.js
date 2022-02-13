import * as Task from '../factory/task'

describe('Task', () => {
    it ('should show the task page', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
    })

    it ('should show the task form', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
      Task.open_task_form()
    })

    it ('should show the task form and close it', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
      Task.open_task_form()
      cy.get('button').contains('Cancel').first().click()
    })

    it ('should show the task form and show errors for empty form submission', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
      Task.open_task_form()
      Task.check_empty_form()
    })

    it ('should create a general task', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
      Task.open_task_form()
      Task.filled_task("GENERAL")
    })

    it ('should create a safety task', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
      Task.open_task_form()
      Task.filled_task("SAFETY")
    })

    it ('should create a sanitation task', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
      Task.open_task_form()
      Task.filled_task("SANITATION")
    })

    it ('should create a inventory task', () => {
      cy.clearLocalStorage()
      Task.open_task_page()
      Task.open_task_form()
      Task.filled_task("INVENTORY")
    })
})
