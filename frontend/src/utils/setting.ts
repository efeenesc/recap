import { db } from "$lib/wailsjs/go/models.ts"
import type { ExtendedSettings, ExtendedSettingDisplay, BasicSetting, CategorizedSettings } from "../types/ExtendedSettings.interface.ts";

export const joinDisplaySettings = (config: BasicSetting, displayVals: ExtendedSettings): CategorizedSettings => {
  const val: CategorizedSettings = {};

  Object.keys(config).forEach(key => {
    const category = displayVals[key].Category;
    
    const value: ExtendedSettingDisplay = {
      [key]: {
        ...displayVals[key],
        Value: config[key]
      }
    }

    console.log(value);
    if (!val[category])
      val[category] = {}
    
    val[category][key] = { ...displayVals[key], Value: config[key] }
  })

  return val;
}