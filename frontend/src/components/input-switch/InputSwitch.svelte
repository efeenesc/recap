<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { SettingInputType } from "../../types/ExtendedSettings.interface.ts";
    import { SelectFolder } from "$lib/wailsjs/go/app/AppMethods.js";
    import Toggle from "../checkbox/Toggle.svelte";

    export let inputType: SettingInputType;
    export let inputValue: string | number;
    export let inputOptions: string[];
    export let checked: boolean = false;
    export let id: string;
    let _class: string;
    export { _class as class };

    $: checked = inputValue == 1 ? true : false; // Not doing type equality to make 1 == '1' return true

    const dispatch = createEventDispatcher();

    function handleChange(event: any) {
        let returnValue;

        // Check if the event is a string
        if (typeof event === "string") {
            returnValue = event;
        } else if (event.target && event.target.value !== undefined) {
            // If there's a target with a value, use that
            returnValue = event.target.value;
        } else if (event.detail && event.detail.checked !== undefined) {
            // If there's a detail with checked, use that
            returnValue = event.detail.checked;
        }

        // Fallback to checked if returnValue is undefined
        dispatch("changed", { changed: returnValue || checked });
    }

    async function pickFolder() {
        try {
            const dirHandle = await SelectFolder();
            handleChange(dirHandle);
        } catch (err) {
            console.error(err);
        }
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
{#if inputType === "FolderPicker"}
    <div
        class="flex items-center justify-center gap-2 my-4 h-12 overflow-hidden"
    >
        <button
            on:click={pickFolder}
            class="flex-shrink-0 h-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg shadow transition duration-300 ease-in-out transform focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
        >
            Select folder
        </button>
        <div
            class="folder-name bg-gray-200 h-12 overflow-hidden rounded-lg p-3 shadow-inner overflow-x-scroll"
        >
            <span class="text-gray-900 text-nowrap">{inputValue}</span>
        </div>
    </div>
{:else if inputType === "ExtendedTextInput"}
    <textarea
        {id}
        class="w-full min-h-[100px] bg-gray-200 focus:bg-white p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-y {_class}"
        on:change={handleChange}
        value={inputValue}
    ></textarea>
{:else if inputType === "Boolean"}
    <div
        class="flex items-center bg-gray-200 w-fit h-fit rounded-lg px-4 py-2 text-neutral-900 my-2"
    >
        <div class="w-[30px]">
            {checked ? "On" : "Off"}
        </div>

        <Toggle
            {id}
            bind:checked
            class="my-2 flex ml-4 pl-4"
            on:checked={(e) => handleChange(e)}
        ></Toggle>
    </div>
{:else if inputType === "NumberInput"}
    <input
        {id}
        class="w-fit max-w-24 p-2 bg-gray-200 focus:bg-white border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent {_class}"
        on:change={handleChange}
        type="number"
        value={inputValue}
    />
{:else if inputType === "TimePicker"}
    <input
        {id}
        class="w-fit p-2 bg-gray-200 focus:bg-white border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent {_class}"
        on:change={handleChange}
        type="time"
        value={inputValue}
    />
{:else if inputType === "APIPicker"}
    <div class="flex w-fit my-4 p-1 gap-2 bg-gray-200 rounded-lg shadow-inner">
        {#each inputOptions as option}
            <button
                on:click={() => handleChange(option)}
                class="py-2 px-4 rounded-md transition-all capitalize {inputValue ===
                option
                    ? 'bg-white text-blue-600 font-bold shadow'
                    : 'text-gray-600 hover:bg-gray-300'}"
            >
                {option}
            </button>
        {/each}
    </div>
{:else if inputType === "APIModelPicker"}
    <input
        {id}
        type="text"
        class="w-fit p-3 border bg-gray-200 focus:bg-white border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-y {_class}"
        on:change={handleChange}
        value={inputValue}
    />
{:else if inputType === "URLInput"}
    <input
        {id}
        type="url"
        class="w-fit p-3 border bg-gray-200 focus:bg-white border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-y {_class}"
        on:change={handleChange}
        value={inputValue}
    />
{:else}
    <input
        {id}
        type="text"
        class="w-full max-w-4xl p-3 focus:bg-white bg-gray-200 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-y {_class}"
        on:change={handleChange}
        value={inputValue}
    />
{/if}

<slot></slot>

<style global lang="postcss">
    @tailwind base;
    @tailwind components;
    @tailwind utilities;

    .folder-name::-webkit-scrollbar {
        background-color: rgba(0,0,0,0);
        width: 16px;
        height: 10px;
        overflow: hidden;
    }

    /* background of the scrollbar except button or resizer */
    .folder-name::-webkit-scrollbar-track {
        background-color: rgba(0,0,0,0);
        overflow: hidden;
    }

    /* scrollbar itself */
    .folder-name::-webkit-scrollbar-thumb {
        background-color:#babac0;
        border-radius:9999px;
    }

    /* set button(top and bottom of the scrollbar) */
    .folder-name::-webkit-scrollbar-button {display:none}
</style>
