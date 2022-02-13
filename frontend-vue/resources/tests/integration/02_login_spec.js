describe('Login specs', () => {
  it('should display login page if the user is not authenticated', () => {
    cy.clearLocalStorage()

    cy.visit('/#/')
    cy.get('label#label-username').should('contain', 'Username')
    cy.location().should( location => {
      expect(location.hash).to.eq('#/auth/login?redirect=%2F')
    })
  })

  it ('should display error page if the user does not fill field correctly ', () => {
    cy.clearLocalStorage()

    cy.visit('/#/')
    cy.get('button.btn').click()
    cy.get('span.help-block.text-danger').should('contain', 'The username field is required.')
    cy.get('span.help-block.text-danger').should('contain', 'The password field is required.')
  })
})
