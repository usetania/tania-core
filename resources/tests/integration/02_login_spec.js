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

  it('should redirect to homepage for existing user who already have farm', () => {
    cy.clearLocalStorage()

    cy.visit('/#/')
    cy.get('input#username').type('user').should('have.value', 'user')
    cy.get('input#password').type('password')
    cy.get('button.btn').click()

    cy.location().should( location => {
      expect(location.hash).to.eq('#/')
    })

    cy.get('a.farm-current span').should('contain', 'FarmName2')

  })

  // skip this test, we will enable when the auth endpoint
  it.skip ('should redirect to intro page for non existing user', () => {
    cy.clearLocalStorage()

    cy.visit('/#/')
    cy.get('input#username').type('demo').should('have.value', 'demo')
    cy.get('input#password').type('password')
    cy.get('button.btn').click()

    cy.location().should( location => {
      expect(location.hash).to.eq('#/intro/farm')
    })
  })
})
