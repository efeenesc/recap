import { dialogStore, type DialogData } from "$lib/stores/DialogStore.ts";

export const addNewDialog = (newDialog: DialogData) => {
  dialogStore.update(prev => {
    if (!prev) return [newDialog];

    return [newDialog, ...prev]
  })
}

export const removeFirstDialog = () => {
  dialogStore.update(prev => {
    if (!prev || prev.length === 0) return prev; // Return if no dialogs exist

    return prev.slice(1); // Remove the first dialog
  });
};