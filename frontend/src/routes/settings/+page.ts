import { GetConfig, GetDisplayValues } from "$lib/wailsjs/go/app/AppMethods.js"
import { config } from "$lib/wailsjs/go/models.ts"
import type { BasicSetting, ExtendedSettings } from "../../types/ExtendedSettings.interface.ts";
import { joinDisplaySettings } from "../../utils/setting.ts";

// Pulls screenshot from database
async function pullSettingsFromDb(): Promise<config.AppConfig | undefined> {
  try {
    const res = await GetConfig();
    return res;
  } catch (err) {
    console.error(`Could not get settings`, err);
    return undefined;
  }
}

async function pullDisplayValuesFromDb() {
  try {
    const res = await GetDisplayValues();
    return res as ExtendedSettings;
  } catch (err) {
    console.error(`Could not get display values`, err);
    return undefined;
  }
}

/** @type {import('./$types').PageLoad} */
export const load = async () => {
  return {
    streamed: {
      items: new Promise(async (resolve, reject) => {
        const basicSettings = await pullSettingsFromDb() as unknown as BasicSetting;
        const displayVals = await pullDisplayValuesFromDb();

        if (!basicSettings || !displayVals) reject(0);

        const result = joinDisplaySettings(basicSettings, displayVals!);

        //! Temporarily remove these keys. Automatic reporting is not implemented yet.
        delete result["Reports"]["ReportAutoEnabled"]
        delete result["Reports"]["ReportAutoAt"]

        console.log(result);

        // Handle cases where result is not found
        if (!result) {
          reject(0);
        } else {
          resolve(result);
        }
      })
    }
  };
};
