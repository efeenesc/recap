import type { db } from "$lib/wailsjs/go/models.ts";

export interface ExtendedScreenshot extends db.CaptureScreenshotImage {
  Visible?: boolean;
  Date: string;
  Time: string;
  Selected: boolean;
}

export interface DatedScreenshot {
  [key: string]: ExtendedScreenshot[]
}