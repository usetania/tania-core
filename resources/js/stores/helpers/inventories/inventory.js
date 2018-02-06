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

export const AgrochemicalQuantityUnits = [
  { key: 'PACKETS',  label: 'Packets' },
  { key: 'BOTTLES', label: 'Bottles' },
  { key: 'BAGS', label: 'Bags' },
]

export function FindAgrochemicalQuantityUnit(key) {
  var quantityUnit =  AgrochemicalQuantityUnits.find(item => item.key === key.toLowerCase())
  return quantityUnit ? quantityUnit.label : ''
}

export const ChemicalTypes = [
  { key: 'DISINFECTANT',  label: 'Disinfectant & Sanitizer' },
  { key: 'FERTILIZER', label: 'Fertilizer' },
  { key: 'HORMONE', label: 'Hormone & Growth Agent' },
  { key: 'MANURE', label: 'Manure' },
  { key: 'PESTICIDE', label: 'Pesticide' },
]

export function FindChemicalType(key) {
  var chemicalType =  ChemicalTypes.find(item => item.key === key.toLowerCase())
  return chemicalType ? chemicalType.label : ''
}

export const GrowingMediumQuantityUnits = [
  { key: 'BAGS',  label: 'Bags' },
  { key: 'CUBIC_METRE', label: 'Cubic Metre' },
]

export function FindGrowingMediumQuantityUnit(key) {
  var quantityUnit =  GrowingMediumQuantityUnits.find(item => item.key === key.toLowerCase())
  return quantityUnit ? quantityUnit.label : ''
}
