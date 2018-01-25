export const ReservoirTypes = [
  { key: "TAP", label: "Tap / Well" },
  { key: "BUCKET", label: "Water Tank / Cistern" }
]

export function FindReservoirType(key) {
  return ReservoirTypes.find(item => item.key === key)
}
