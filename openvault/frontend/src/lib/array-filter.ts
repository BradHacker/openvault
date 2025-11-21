export function filterNullish<T>(this: Array<T>): Array<NonNullable<T>> {
  return this.filter((item) => item !== null && item !== undefined);
}
Array.prototype.filterNullish = filterNullish;

declare global {
  interface Array<T> {
    filterNullish(): Array<NonNullable<T>>;
  }
}
