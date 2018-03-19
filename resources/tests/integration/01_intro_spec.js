import { login } from '../factory/auth'
import { filled_farm, filled_reservoir, filled_area } from '../factory/farm'

describe('Intro specs', () => {

  describe('Intro farm page', () => {
    it('should display farm intro pages from new user', () => {

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
      cy.get('span.help-block.text-danger').should('contain', 'The name field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The type field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The latitude field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The latitude field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The country field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The city field is required.')
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
      cy.get('span.help-block.text-danger').should('contain', 'The name field is required.')
    })

    it ('should to redirect to the reservoir page, if the reservoir empty', () => {
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
      cy.get('span.help-block.text-danger').should('contain', 'The name field is required.')
      cy.get('span.help-block.text-danger').should('contain', 'The size field is required.')
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
