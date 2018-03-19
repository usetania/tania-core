import { login } from '../factory/auth'

describe('Logout specs', () => {
  it ('should have logout link', () => {
    login()
    cy.get('#signout').should('contain', 'Sign Out')
    cy.get('#signout').click()
    cy.location().should( location => {
      expect(location.hash).to.eq('#/auth/login')
    })
  })

  it ('should not visit the dashboard page if the user already logout', () => {
    login()
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
    login()
    cy.get('#signout').should('contain', 'Sign Out')
    cy.get('#signout').click()

    expect(localStorage.getItem('vuex')).to.eq(null)
  })
})
