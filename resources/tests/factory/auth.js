export function login() {
  cy.clearLocalStorage()

  cy.visit('/#/')
  // Label
  cy.get('label#label-username').should('contain', 'Username')
  cy.get('label#label-password').should('contain', 'Password')

  // Typing the form
  cy.get('input#username').type('demo').should('have.value', 'demo')
  cy.get('input#password').type('password')
  cy.get('button.btn').click()
}

export function login_as_user() {
  cy.clearLocalStorage()

  cy.visit('/#/')
  // Label
  cy.get('label#label-username').should('contain', 'Username')
  cy.get('label#label-password').should('contain', 'Password')

  // Typing the form
  cy.get('input#username').type('user').should('have.value', 'user')
  cy.get('input#password').type('password')
  cy.get('button.btn').click()
}
