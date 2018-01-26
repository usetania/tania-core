export function filled_farm() {
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

export function filled_reservoir() {
  // Label
  cy.get('label#label-name').should('contain', 'Reservoir Name')
  cy.get('label#label-source').should('contain', 'Source')

  // Typing the form
  cy.get('#name').type('Reservoir1').should('have.value', 'Reservoir1')
  cy.get('#type').select('TAP')
  cy.get('#capacity').should('not.exist')
  cy.get('#type').select('BUCKET')
  cy.get('#capacity').should('exist').type('12').should('have.value', '12')
  cy.get('button.btn[type=submit]').click()
}

export function filled_area() {
  // Label
  cy.get('label#label-name').should('contain', 'Area Name')
  cy.get('label#label-size').should('contain', 'Size')
  cy.get('label#label-type').should('contain', 'Type')
  cy.get('label#label-location').should('contain', 'Locations')
  cy.get('label#label-reservoir').should('contain', 'Select Reservoir')

  // Typing the form
  cy.get('#name').type('Area1').should('have.value', 'Area1')
  cy.get('#size').type('2').should('have.value', '2')
  cy.get('#size_unit').select('Ha')
  cy.get('#type').select('SEEDING')
  cy.get('#location').select('OUTDOOR')
  cy.get('#reservoir').select('Reservoir1')
  cy.get('button.btn[type=submit]').click()
}
