export const ReservoirTypes = [
  { key: "tap", label: "Tap / Well" },
  { key: "bucket", label: "Water Tank / Cistern" }
]

export function FindReservoirType(key) {
  return ReservoirTypes.find(item => item.key === key)
}
