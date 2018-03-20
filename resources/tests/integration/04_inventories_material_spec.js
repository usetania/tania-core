import * as Material from '../factory/material'

describe('Inventories', () => {
    it ('should show the materials page', () => {
      Material.open_material_page()
    })

    it ('should show the materials page and open the materials form modal', () => {
      Material.open_material_page()
      Material.open_material_form()
    })

    it ('should create a seed material', () => {
      Material.open_material_page()
      Material.open_material_form()
      Material.filled_material_seed()
    })

    // it ('should create a growing medium material', () => {
    //   Material.open_material_page()
    //   Material.open_material_form()
    //   Material.filled_material_growing_medium()
    // })

    it ('should create an agrochemical material', () => {
      Material.open_material_page()
      Material.open_material_form()
      Material.filled_material_agrochemical()
    })
})
