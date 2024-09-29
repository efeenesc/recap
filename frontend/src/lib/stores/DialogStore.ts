import { type Writable, writable } from "svelte/store";

export interface DialogData {
  title: string
  description: string,
  buttonLink?: string,
  buttonLinkDescription?: string
}

export const dialogStore = writable<DialogData[]>();