<script lang="ts">
    import { goto } from "$app/navigation";
  import { dialogStore, type DialogData } from "$lib/stores/DialogStore.ts"
  import { onDestroy, onMount } from "svelte";
    import type { Unsubscriber } from "svelte/store";
    import { removeFirstDialog } from "../../utils/dialog.ts";

  let dialogs: DialogData[] = [];
  let currentDialog: DialogData | undefined;
  let dialogUnsubscriber: Unsubscriber;

  function subscribeToDialog() {
    dialogUnsubscriber = dialogStore.subscribe((data) => {
      dialogs = data;

      if (dialogs && dialogs.length > 0) {
        showDialog(dialogs[0]);
      } else {
        currentDialog = undefined;
      }
    })
  }

  function showDialog(dialog: DialogData) {
    console.log("Showing", dialog)
    currentDialog = dialog;
  }

  function closeDialog() {
    removeFirstDialog();
  }

  function buttonClicked() {
    if (!currentDialog?.buttonLink) return;

    const targetRoute = currentDialog!.buttonLink;
    
    closeDialog();
    goto(targetRoute);
  }

  onDestroy(() => {
    dialogUnsubscriber();
  })

  onMount(() => {
    subscribeToDialog();
  })
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-full flex justify-center items-center fixed top-0 left-0 right-0 bottom-0 transition-all {currentDialog ? 'pointer-events-auto opacity-100' : 'pointer-events-none opacity-0'}">
  {#if currentDialog}
    <div on:click={closeDialog} class="absolute w-full h-full z-40 bg-black opacity-60"></div>
    <div class="flex flex-col z-50 rounded-xl bg-white text-black min-w-[40%] max-w-[70%] p-6">
      <h1 class="text-2xl font-bold mb-4">
        {currentDialog.title}
      </h1>
      <p class="text-lg">
        {currentDialog.description}
      </p>
      {#if currentDialog.buttonLink}
        <div on:click={buttonClicked} class="flex w-fit p-3">
          {currentDialog.buttonLink}
        </div>
      {/if}
    </div>
  {/if}
</div>

<style global lang="postcss">
  @tailwind utilities;
  @tailwind components;
  @tailwind base;
</style>