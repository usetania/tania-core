export const pageLength = 10

export function calculateNumberOfPages (length) {
  return (length % pageLength == 0) ? Math.floor(length/pageLength) : Math.floor(length/pageLength) + 1
}