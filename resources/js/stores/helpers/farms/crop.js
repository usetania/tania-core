export const CropContainers = [
  { key: 'pot',  label: 'Pot' },
  { key: 'tray',  label: 'Tray' },
]

export function FindCropContainer(key) {
  return CropContainers.find(item => item.key === key)
}