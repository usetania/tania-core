import moment from 'moment';

export function filled_material_seed() {
  // Label
  cy.get('label#label-material-type').should('contain', 'Choose type of material')
  cy.get('label#label-name').should('contain', 'Variety Name')
  cy.get('label#label-plant-type').should('contain', 'Plant Type')
  cy.get('label#label-produced-by').should('contain', 'Produced by')
  cy.get('label#label-quantity').should('contain', 'Quantity')
  cy.get('label#label-quantity-unit').should('contain', 'Unit')
  cy.get('label#label-price-per-unit').should('contain', 'Price per Unit')
  cy.get('label#label-expiration-date').should('contain', 'Expiration date')
  cy.get('label#label-notes').should('contain', 'Additional Notes')

  cy.get('button.btn[type=submit]').click()
  cy.get('span.help-block.text-danger').should('contain', 'The name field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The produced by field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The quantity field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The price field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The expiration date field is required.')

  // Typing the form
  cy.get('input#name').type('Onion'+moment().format())
  cy.get('select#plant_type').select('VEGETABLE')
  cy.get('input#produced_by').type('Onion Producer')
  cy.get('input#quantity').type('1000')
  cy.get('select#quantity_unit').select('SEEDS')
  cy.get('input#price_per_unit').type('1')
  cy.get('input#expiration_date').click()
  cy.get('span.next').first().click()
  cy.get('span').contains(15).first().click()
  cy.get('textarea#notes').type('Onion Notes')

  cy.get('button.btn[type=submit]').click()
}
