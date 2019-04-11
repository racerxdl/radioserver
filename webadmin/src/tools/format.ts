const subMultiple = ["", "m", "&micro;", "n", "p", "f", "a", "z", "y"];
const multiple = ["", "k", "M", "G", "T", "P", "E", "Z", "Y"];

function toNotationUnit(v: number): [number, string] {
  let unit;
  let counter = 0;
  let value = v;
  if (value < 1) {
    while (value < 1) {
      counter++;
      value = value * 1e3;
      if (counter === 8) break;
    }
    unit = subMultiple[counter];
  } else {
    while (value > 1000) {
      counter++;
      value = value / 1e3;
      if (counter === 8) break;
    }
    unit = multiple[counter];
  }
  value = Math.round(value * 1e2) / 1e2;
  return [value, unit];
}

function toHzNotation(v: number): string {
  const z = toNotationUnit(v);
  return z[0].toLocaleString() + ' ' + z[1] + 'Hz';
}

function toBytesPerSecNotation(v: number): string {
  const z = toNotationUnit(v);
  return z[0].toLocaleString() + ' ' + z[1] + 'b/s';
}

export default {
  toNotationUnit,
  toHzNotation,
  toBytesPerSecNotation,
}
