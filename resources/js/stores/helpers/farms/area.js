export const AreaTypes = [
  { key: 'nursery',  label: 'Nursery / Seeding' },
  { key: 'growing', label: 'Growing Area' }
]

export const AreaLocations = [
  { key: 'outdoor',  label: 'Field (Outdoor)' },
  { key: 'indoor', label: 'Greenhouse (Indoor)' }
]

export const AreaSizeUnits = [
  { key: 'Ha',  label: 'hectare' },
  { key: 'm2', label: 'meter square' }
]

export function FindAreaType(key) {
  return AreaTypes.find(item => item.key === key)
}

export function FindAreaLocation(key) {
  return AreaLocations.find(item => item.key === key)
}

export function FindAreaSizeUnit(key) {
  return AreaSizeUnits.find(item => item.key === key)
}
