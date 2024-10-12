<script lang="ts">
	import { fade, scale } from 'svelte/transition';
    import { expoOut } from 'svelte/easing';
    import { dialogStore, type DialogData } from "$lib/stores/DialogStore.ts";
    import { onDestroy, onMount } from "svelte";
    import type { Unsubscriber } from "svelte/store";
    import { removeFirstDialog } from "../../utils/dialog.ts";

    let _class: string;
    export { _class as class };
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
        currentDialog = dialog;
    }

    function invokeCallback(
        buttonType: "primary" | "secondary" | "tertiary"
    ): void {
        if (!currentDialog) return;

        try {
            switch (buttonType) {
                case "primary":
                    currentDialog.primaryButtonCallback?.();
                    break;
                case "secondary":
                    currentDialog.secondaryButtonCallback?.();
                    break;
                case "tertiary":
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

<div
    class="w-full h-full flex justify-center items-center fixed top-0 left-0 right-0 bottom-0 transition-all {_class}"
    class:pointer-events-auto={currentDialog}
    class:opacity-100={currentDialog}
    class:pointer-events-none={!currentDialog}
    class:opacity-0={!currentDialog}
>
    {#if currentDialog}
        <div
            transition:fade={{ duration: 200 }}
            on:click={closeDialog}
            class="absolute w-full h-full z-40 bg-black opacity-60"
            role="presentation"
        ></div>
        <div
            transition:scale={{start: 0.9, opacity: 0, easing: expoOut, duration: 500 }}
            class="flex flex-col z-50 rounded-xl bg-white text-black min-w-[40%] max-w-[70%] p-6"
        >
            <h1 class="text-2xl lg:text-3xl font-extrabold mb-4">
                {currentDialog.title}
            </h1>
            <p class="text-lg lg:text-xl font-medium text-neutral-800">
                {currentDialog.description}
            </p>
            <div class="flex gap-5 mt-12 text-md lg:text-xl">
                {#if currentDialog.primaryButtonCallback}
                    <button
                        on:click={() => invokeCallback("primary")}
                        class="flex shadow-md w-full font-bold items-center justify-center py-4 transition-colors hover:bg-blue-400 bg-blue-500 text-white rounded-xl"
                    >
                        {currentDialog.primaryButtonName}
                    </button>
                {/if}
                {#if currentDialog.secondaryButtonCallback}
                    <button
                        on:click={() => invokeCallback("secondary")}
                        class="flex shadow-md w-full font-bold items-center justify-center py-4 transition-colors hover:bg-gray-200 bg-gray-300 text-neutral-700 rounded-xl"
                    >
                        {currentDialog.secondaryButtonName}
                    </button>
                {/if}
                {#if currentDialog.tertiaryButtonCallback}
                    <button
                        on:click={() => invokeCallback("tertiary")}
                        class="flex shadow-md items-center justify-center w-full font-bold py-4 transition-colors hover:bg-gray-50 bg-gray-100 text-neutral-700 rounded-xl"
                    >
                        {currentDialog.tertiaryButtonName}
                    </button>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
