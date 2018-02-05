export const InventoryTypes = [
  { key: 'seed',  label: 'Seed' },
  { key: 'growing_medium', label: 'Growing Medium' },
  { key: 'agrochemical', label: 'Agrochemical' },
  { key: 'label_and_crop_support', label: 'Label and Crop Support' },
  { key: 'seeding_container', label: 'Seeding Container' },
  { key: 'post_harvest_supply', label: 'Post Harvest Supply' },
  { key: 'other', label: 'Other Material' },
]

export function FindInventoryType(key) {
  var inventoryType = InventoryTypes.find(item => item.key === key.toLowerCase())
  return inventoryType ? inventoryType.label : ''
}

export const QuantityUnits = [
  { key: 'SEEDS',  label: 'Seeds' },
  { key: 'PACKETS', label: 'Packets' },
  { key: 'GRAM', label: 'Gram' },
  { key: 'KILOGRAM', label: 'Kilogram' },
]

export function FindQuantityUnit(key) {
  var quantityUnit =  QuantityUnits.find(item => item.key === key.toLowerCase())
  return quantityUnit ? quantityUnit.label : ''
}