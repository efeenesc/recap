<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { SettingInputType } from "../../types/ExtendedSettings.interface.ts";
    import { SelectFolder } from "$lib/wailsjs/go/app/AppMethods.js"
    import Toggle from "../checkbox/Toggle.svelte";

    export let inputType: SettingInputType;
    export let inputValue: string | number;
    export let checked: boolean = false;
    export let id: string;
    let _class: string;
    export { _class as class };

    const dispatch = createEventDispatcher();

    function handleChange(event: any) {
      console.log(event);
        dispatch("changed", { changed: event || checked });
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
    <div class="flex items-center justify-center gap-4 my-2">
        <div class="bg-neutral-200 flex-shrink-0 bg-opacity-80 rounded-xl p-2">
            <div on:click={pickFolder} class="flex cursor-pointer items-center justify-center px-4 py-2 transition-all rounded-xl text-lg bg-white text-black hover:scale-[99%] active:scale-[95%]">
                Select folder
            </div>
        </div>

        <div class="bg-neutral-200 bg-opacity-80 rounded-xl p-2 flex-grow">
          <div
              class="rounded-xl text-xl bg-white text-black px-4 py-2"
          >
              {inputValue}
          </div>
      </div>
    </div>
{:else if inputType === "ExtendedTextInput"}
    <textarea
        {id}
        class="w-full min-h-[200px] {_class}"
        on:change={handleChange}
        value={inputValue}
    ></textarea>
{:else if inputType === "Boolean"}
    <Toggle {id} class={_class} bind:checked on:change={handleChange}></Toggle>
{:else if inputType === "NumberInput"}
    <input
        {id}
        class={_class}
        on:change={handleChange}
        type="number"
        value={inputValue}
    />
{:else if inputType === "TimePicker"}
    <input
        {id}
        class={_class}
        on:change={handleChange}
        type="time"
        value={inputValue}
    />
{:else if inputType === "APIPicker"}
    <div
        class="flex my-2 p-2 gap-2 w-fit bg-opacity-80 bg-neutral-200 rounded-xl text-lg text-neutral-900"
    >
        <div
            on:click={() => handleChange("ollama")}
            class="py-2 px-4 cursor-pointer transition-all rounded-xl {inputValue ===
            'ollama'
                ? 'bg-white text-neutral-950 font-bold'
                : ''}"
        >
            Ollama
        </div>
        <div
            on:click={() => handleChange("gemini")}
            class="py-2 px-4 cursor-pointer transition-all rounded-xl {inputValue ===
            'gemini'
                ? 'bg-white text-neutral-950 font-bold'
                : ''}"
        >
            Gemini
        </div>
    </div>
{/if}

<slot></slot>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;

    .neu {
        box-shadow:
            20px 20px 60px #232222,
            -20px -20px 60px #2f2e2e;
    }
</style>
