export function login() {
  cy.clearLocalStorage()

  cy.visit('/#/')
  // Label
  cy.get('label#label-username').should('contain', 'Username')
  cy.get('label#label-password').should('contain', 'Password')

  // Typing the form
  cy.get('input#username').type('tania').should('have.value', 'tania')
  cy.get('input#password').type('tania')
  cy.get('button.btn').click()
}
