import { type Writable, writable } from "svelte/store";

export interface DialogData {
  title: string
  description: string,
  primaryButtonCallback?: () => void,
  primaryButtonName?: string,
  secondaryButtonCallback?: () => void,
  secondaryButtonName?: string,
  tertiaryButtonCallback?: () => void,
  tertiaryButtonName?: string,
}

export const dialogStore = writable<DialogData[]>();