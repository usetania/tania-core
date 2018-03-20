import moment from 'moment';
import { login } from './auth'

export function open_material_page() {
  login()
  cy.get('#inventories').click()
  cy.get('#materials').click()
  cy.location().should( location => {
    expect(location.hash).to.eq('#/materials')
  })
}

export function open_material_form() {
  cy.get('#materialsform').click() 
  cy.get('form').should('be.visible') 
}

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

export function filled_material_growing_medium() {
  cy.get('select#material_type').select('growing_medium')
  // Label
  cy.get('label#label-material-type').should('contain', 'Choose type of material')
  cy.get('label#label-name').should('contain', 'Name')
  cy.get('label#label-produced-by').should('contain', 'Produced by')
  cy.get('label#label-quantity').should('contain', 'Quantity')
  cy.get('label#label-quantity-unit').should('contain', 'Unit')
  cy.get('label#label-price-per-unit').should('contain', 'Price per Unit')
  cy.get('label#label-notes').should('contain', 'Additional Notes')

  cy.get('button.btn[type=submit]').click()
  cy.get('span.help-block.text-danger').should('contain', 'The name field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The produced by field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The quantity field is required.')
  cy.get('span.help-block.text-danger').should('contain', 'The price field is required.')

  // Typing the form
  cy.get('input#name').type('Growing Medium'+moment().format())
  cy.get('input#produced_by').type('Growers')
  cy.get('input#quantity').type('1000')
  cy.get('select#quantity_unit').select('BAGS')
  cy.get('input#price_per_unit').type('1')
  cy.get('textarea#notes').type('Growers Notes')

  cy.get('button.btn[type=submit]').click()
}

export function filled_material_agrochemical() {
  cy.get('select#material_type').select('agrochemical')
  // Label
  cy.get('label#label-material-type').should('contain', 'Choose type of material')
  cy.get('label#label-name').should('contain', 'Name')
  cy.get('label#label-chemical-type').should('contain', 'Chemical Type')
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
  cy.get('input#name').type('Nitrogen Fertilizer'+moment().format())
  cy.get('select#chemical_type').select('FERTILIZER')
  cy.get('input#produced_by').type('Nitrogen Producer')
  cy.get('input#quantity').type('1000')
  cy.get('select#quantity_unit').select('PACKETS')
  cy.get('input#price_per_unit').type('1')
  cy.get('input#expiration_date').click()
  cy.get('span.next').first().click()
  cy.get('span').contains(15).first().click()
  cy.get('textarea#notes').type('FERTILIZER Notes')

  cy.get('button.btn[type=submit]').click()
}
