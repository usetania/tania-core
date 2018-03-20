import * as MaterialTask from '../factory/task'

describe('Inventories', () => {
    it ('should show the materials page', () => {
      MaterialTask.open_material_page()
    })

    it ('should show the materials page and open the materials form modal', () => {
      MaterialTask.open_material_page()
      MaterialTask.open_material_form()
    })

    it ('should create a seed material', () => {
      MaterialTask.open_material_page()
      MaterialTask.open_material_form()
      MaterialTask.filled_material_seed()
    })

    it ('should create an agrochemical material', () => {
      MaterialTask.open_material_page()
      MaterialTask.open_material_form()
      MaterialTask.filled_material_agrochemical()
    })
})
