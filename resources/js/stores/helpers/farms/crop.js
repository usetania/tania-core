export const CropContainers = [
  { key: 'POT',  label: 'Pot' },
  { key: 'TRAY',  label: 'Tray' },
]

export function FindCropContainer(key) {
  return CropContainers.find(item => item.key === key)
}