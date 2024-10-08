import { type Writable, writable } from "svelte/store";

// Thanks Svelte, I didn't want to pass props to my <slot> elements anyway. Who needs that when you have stores?
// Let's have a scroll store. Stores for everything
export const scrollStore = writable<number>(0);

// I should be creating a function in /utils instead of here. This file shouldn't exist in the first place, I'm not creating another file
export const updateScroll = (scrollTop: number) => {
  scrollStore.set(scrollTop);
}