export function transformObjectToOption<T extends object>(obj: T) {
  console.log(`transformObjectToOption',`, obj);
  return Object.entries(obj).map(([value, label]) => ({
    value,
    label
  })) as Common.OptionWithKey<keyof T>[];
}
