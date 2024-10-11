import type { ExtendedSettings, BasicSetting, CategorizedSettings } from "../types/ExtendedSettings.interface.ts";

export const joinDisplaySettings = (config: BasicSetting, displayVals: ExtendedSettings): CategorizedSettings => {
  const val: CategorizedSettings = {};

  Object.keys(config).forEach(key => {
    const category = displayVals[key].Category;

    if (!val[category])
      val[category] = {}
    
    val[category][key] = { ...displayVals[key], Value: config[key] }
  })

  return val;
}