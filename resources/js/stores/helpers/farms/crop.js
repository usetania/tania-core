export const Containers = [
  { key: 'POT',  label: 'Pot' },
  { key: 'TRAY',  label: 'Tray' },
]

export function FindContainer(key) {
  return Containers.find(item => item.key === key)
}
