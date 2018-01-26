import { login_as_user } from '../factory/auth'

describe('Logout specs', () => {
  it ('should have logout link', () => {
    login_as_user()
    cy.get('#signout').should('contain', 'Sign Out')
    cy.get('#signout').click()
    cy.location().should( location => {
      expect(location.hash).to.eq('#/auth/login')
    })
  })

  it ('should not visit the dashboard page if the user already logout', () => {
    login_as_user()
    cy.get('#signout').should('contain', 'Sign Out')
    cy.get('#signout').click()

    cy.visit('/#/')
    cy.location().should( location => {
      expect(location.hash).to.eq('#/auth/login?redirect=%2F')
    })

    cy.visit('/#/intro/area')
    cy.location().should( location => {
      expect(location.hash).to.eq('#/auth/login?redirect=%2Fintro%2Farea')
    })
  })

  it ('should not have localstorage with vuex key in the browser', () => {
    login_as_user()
    cy.get('#signout').should('contain', 'Sign Out')
    cy.get('#signout').click()

    expect(localStorage.getItem('vuex')).to.eq(null)
  })
})
