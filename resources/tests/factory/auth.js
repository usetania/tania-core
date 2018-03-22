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

export function logout() {
  cy.get('#signout').should('contain', 'Sign Out')
  cy.get('#signout').click()
  cy.location().should( location => {
    expect(location.hash).to.eq('#/auth/login')
  })
}