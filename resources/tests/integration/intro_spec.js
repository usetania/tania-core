function login() {
  cy.clearLocalStorage()

  cy.visit('/#/')
  cy.get('label#label-username').should('contain', 'Username')
  cy.get('label#label-password').should('contain', 'Password')
  cy.get('input#username').type('demo').should('have.value', 'demo')
  cy.get('input#password').type('password')
  cy.get('button.btn').click()
}

function filled_farm() {
  cy.get('label#label-name').should('contain', 'Farm Name')
  cy.get('label#label-description').should('contain', 'Farm Description')
  cy.get('label#label-type').should('contain', 'Farm Type')
  cy.get('label#label-location').should('contain', 'Location')
  cy.get('label#label-country').should('contain', 'Country')
  cy.get('label#label-city').should('contain', 'City')


  cy.get('input#name').type('FarmName').should('have.value', 'FarmName')

  cy.get('button.btn[type=submit]').click()
}

describe('Intro specs', () => {

  describe('Intro farm page', () => {
    it('should display farm intro pages from non existing user', () => {

      login()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/farm')
      })
    })

    it ('should redirect into farm intro', () => {
      login()

      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/farm')
      })

      // always redirect into farm intro when the farm is empty
      cy.visit('#/intro/reservoir')
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/farm')
      })

      // always redirect into farm intro when the farm is empty
      cy.visit('/#/intro/area')
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/farm')
      })

    })

    it ('should display error page when the form is not filled correctly', () => {
      login()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/farm')
      })

      cy.get('button.btn[type=submit]').click()

      cy.get('span.h4.font-bold').should('contain', 'Create Farm')
      cy.get('span.help-block.text-danger').should('contain', 'The farm.name field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The farm.farm_type field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The farm.latitude field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The farm.latitude field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The farm.country_code field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The farm.city_code field is required.')
    })

    it ('should navigate into next page when the form filled correctly', () => {
      login()
      filled_farm()
    })
  })
})
