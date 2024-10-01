<script lang="ts">
  import { goto } from "$app/navigation";
  import { dialogStore, type DialogData } from "$lib/stores/DialogStore.ts";
  import { onDestroy, onMount } from "svelte";
  import type { Unsubscriber } from "svelte/store";
  import { removeFirstDialog } from "../../utils/dialog.ts";

  let dialogs: DialogData[] = [];
  let currentDialog: DialogData | undefined;
  let dialogUnsubscriber: Unsubscriber;

  function subscribeToDialog(): void {
    dialogUnsubscriber = dialogStore.subscribe((data) => {
      dialogs = data;

      if (dialogs && dialogs.length > 0) {
        showDialog(dialogs[0]);
      } else {
        currentDialog = undefined;
      }
    });
  }

  function showDialog(dialog: DialogData): void {
    console.log("Showing", dialog);
    currentDialog = dialog;
  }
  
  function invokeCallback(buttonType: 'primary' | 'secondary' | 'tertiary'): void {
    if (!currentDialog) return;

    try {
      switch (buttonType) {
        case 'primary':
          currentDialog.primaryButtonCallback?.();
          break;
        case 'secondary':
          currentDialog.secondaryButtonCallback?.();
          break;
        case 'tertiary':
          currentDialog.tertiaryButtonCallback?.();
          break;
      }
    } catch (error) {
      console.error(`Error invoking ${buttonType} callback:`, error);
    } finally {
      closeDialog();
    }
  }

  function closeDialog(): void {
    removeFirstDialog();
  }

  onDestroy(() => {
    if (dialogUnsubscriber) {
      dialogUnsubscriber();
    }
  });

  onMount(() => {
    subscribeToDialog();
  });
</script>

<div class="w-full h-full flex justify-center items-center fixed top-0 left-0 right-0 bottom-0 transition-all"
      class:pointer-events-auto={currentDialog}
      class:opacity-100={currentDialog}
      class:pointer-events-none={!currentDialog}
      class:opacity-0={!currentDialog}>
  {#if currentDialog}
    <div on:click={closeDialog} 
          class="absolute w-full h-full z-40 bg-black opacity-60"
          role="presentation">
    </div>
    <div class="flex flex-col z-50 rounded-xl bg-white text-black min-w-[40%] max-w-[70%] p-6">
      <h1 class="text-2xl font-bold mb-4">
        {currentDialog.title}
      </h1>
      <p class="text-lg">
        {currentDialog.description}
      </p>
      {#if currentDialog.primaryButtonCallback}
        <button on:click={() => invokeCallback('primary')} 
                class="flex w-fit p-3 mt-4 bg-blue-500 text-white rounded">
          {currentDialog.primaryButtonName}
        </button>
      {/if}
      {#if currentDialog.secondaryButtonCallback}
        <button on:click={() => invokeCallback('secondary')} 
                class="flex w-fit p-3 mt-2 bg-gray-300 text-black rounded">
          {currentDialog.secondaryButtonName}
        </button>
      {/if}
      {#if currentDialog.tertiaryButtonCallback}
        <button on:click={() => invokeCallback('tertiary')} 
                class="flex w-fit p-3 mt-2 bg-gray-100 text-black rounded">
          {currentDialog.tertiaryButtonName}
        </button>
      {/if}
    </div>
  {/if}
</div>

<style lang="postcss">
  @tailwind utilities;
  @tailwind components;
  @tailwind base;
</style>