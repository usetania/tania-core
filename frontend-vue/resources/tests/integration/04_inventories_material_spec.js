import * as Material from '../factory/material'

describe('Inventories', () => {
    it ('should show the materials page', () => {
      cy.clearLocalStorage()
      Material.open_material_page()
    })

    it ('should show the materials page and open the materials form modal', () => {
      cy.clearLocalStorage()
      Material.open_material_page()
      Material.open_material_form()
    })

    it ('should create a seed material', () => {
      cy.clearLocalStorage()
      Material.open_material_page()
      Material.open_material_form()
      Material.filled_material_seed()
    })

    // it ('should create a growing medium material', () => {
    //   cy.clearLocalStorage()
    //   Material.open_material_page()
    //   Material.open_material_form()
    //   Material.filled_material_growing_medium()
    // })

    it ('should create an agrochemical material', () => {
      cy.clearLocalStorage()
      Material.open_material_page()
      Material.open_material_form()
      Material.filled_material_agrochemical()
    })

    // it ('should create a label and crop support material', () => {
    //   cy.clearLocalStorage()
    //   Material.open_material_page()
    //   Material.open_material_form()
    //   Material.filled_material_label_crop_support()
    // })

    // it ('should create other material', () => {
    //   cy.clearLocalStorage()
    //   Material.open_material_page()
    //   Material.open_material_form()
    //   Material.filled_material_other()
    // })
})
