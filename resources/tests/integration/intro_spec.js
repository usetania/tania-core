function login() {
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

function filled_farm() {
  // Label
  cy.get('label#label-name').should('contain', 'Farm Name')
  cy.get('label#label-description').should('contain', 'Farm Description')
  cy.get('label#label-type').should('contain', 'Farm Type')
  cy.get('label#label-location').should('contain', 'Location')
  cy.get('label#label-country').should('contain', 'Country')
  cy.get('label#label-city').should('contain', 'City')

  // Typing the form
  cy.get('input#name').type('FarmName2').should('have.value', 'FarmName2')
  cy.get('textarea#description').type('Farm long description criteria').should('have.value','Farm long description criteria')
  cy.get('select#type').select('organic')
  cy.get('input#latitude').type('-6.2499848').should('have.value', '-6.2499848')
  cy.get('input#longitude').type('106.68292059999999').should('have.value', '106.68292059999999')
  cy.get('select#country').select('ID')
  cy.get('select#city').select('JK')
  cy.get('button.btn[type=submit]').click()
}

function filled_reservoir() {
  // Label
  cy.get('label#label-name').should('contain', 'Reservoir Name')
  cy.get('label#label-source').should('contain', 'Source')

  // Typing the form
  cy.get('#name').type('Reservoir1').should('have.value', 'Reservoir1')
  cy.get('#type').select('tap')
  cy.get('#capacity').should('not.exist')
  cy.get('#type').select('bucket')
  cy.get('#capacity').should('exist').type('12').should('have.value', '12')
  cy.get('button.btn[type=submit]').click()
}

function filled_area() {
  // Label
  cy.get('label#label-name').should('contain', 'Area Name')
  cy.get('label#label-size').should('contain', 'Size')
  cy.get('label#label-type').should('contain', 'Type')
  cy.get('label#label-location').should('contain', 'Locations')
  cy.get('label#label-reservoir').should('contain', 'Select Reservoir')

  // Typing the form
  cy.get('#name').type('Area1').should('have.value', 'Area1')
  cy.get('#size').type('2').should('have.value', '2')
  cy.get('#size_unit').select('hectare')
  cy.get('#type').select('seeding')
  cy.get('#location').select('outdoor')
  cy.get('#reservoir').select('Reservoir1')
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

    it ('should redirect into farm intro, if the farm is empty', () => {
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

    it ('should go to the next page when the form filled correctly', () => {
      login()
      filled_farm()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/reservoir')
      })
    })
  })

  describe('Intro reservoir', () => {
    it ('should correctly to previous page', () => {
      login()
      filled_farm()
      // should have back navigation
      cy.get('a#back').click()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/farm')
      })
      cy.get('button.btn[type=submit]').click()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/reservoir')
      })
    })

    it ('should display error page when the form is not filled correctly', () => {
      login()
      filled_farm()
      cy.get('button.btn[type=submit]').click()
      cy.get('span.help-block.text-danger').should('contain', 'The reservoir.name field is required.')
    })

    it ('shoul to redirect to the reservoir page, if the reservoir empty', () => {
      login()
      filled_farm()
      // always redirect into farm intro when the farm is empty
      cy.visit('/#/intro/area')
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/reservoir')
      })
    })

    it ('should go to next page when form is filled correctly', () => {
      login()
      filled_farm()
      filled_reservoir()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/area')
      })
    })
  })

  describe('Intro Area', () => {
    it ('should correctly to previous page', () => {
      login()
      filled_farm()
      filled_reservoir()

      // should have back navigation
      cy.get('a#back').click()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/reservoir')
      })
      cy.get('button.btn[type=submit]').click()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/intro/area')
      })
    })
    it ('should display error page when the form is not filled correctly', () => {
      login()
      filled_farm()
      filled_reservoir()
      cy.get('button.btn[type=submit]').click()
      cy.get('span.help-block.text-danger').should('contain', 'The area.name field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The area.size field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The area.size_unit field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The area.type field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The area.location field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The area.reservoir field is required.')
    })
    it ('should go to next page when form is filled correctly', () => {
      login()
      filled_farm()
      filled_reservoir()
      filled_area()
      cy.location().should( location => {
        expect(location.hash).to.eq('#/')
      })
    })
  })
})
