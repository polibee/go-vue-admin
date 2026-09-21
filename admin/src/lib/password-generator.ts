const characterSets = [
  'ABCDEFGHJKLMNPQRSTUVWXYZ',
  'abcdefghijkmnopqrstuvwxyz',
  '23456789',
  '!@#$%^&*_-+=',
]

function randomIndex(max: number) {
  const values = new Uint32Array(1)
  globalThis.crypto.getRandomValues(values)
  return values[0] % max
}

export function generatePassword(length = 16) {
  const size = Math.max(12, length)
  const chars = characterSets.map((set) => set[randomIndex(set.length)])
  const allCharacters = characterSets.join('')
  while (chars.length < size) chars.push(allCharacters[randomIndex(allCharacters.length)])
  for (let index = chars.length - 1; index > 0; index -= 1) {
    const swapIndex = randomIndex(index + 1)
    ;[chars[index], chars[swapIndex]] = [chars[swapIndex], chars[index]]
  }
  return chars.join('')
}
