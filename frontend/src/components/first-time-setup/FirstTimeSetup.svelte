<script lang="ts">
    import VirtualScrollbar from "./../virtual-scrollbar/VirtualScrollbar.svelte";
    import { writable } from "svelte/store";
    import InputSwitch from "./../input-switch/InputSwitch.svelte";
    import { fade, scale } from "svelte/transition";
    import { expoOut } from "svelte/easing";
    import type { CategorizedSettings } from "./../../types/ExtendedSettings.interface.ts";
    import { createEventDispatcher } from "svelte";
    import {
        GetConfig,
        GetDisplayValues,
    } from "$lib/wailsjs/go/app/AppMethods.js";
    import { config } from "$lib/wailsjs/go/models.ts";
    import type {
        BasicSetting,
        ExtendedSettings,
    } from "../../types/ExtendedSettings.interface.ts";
    import { joinDisplaySettings } from "../../utils/setting.ts";
    import { onMount } from "svelte";
    import { get } from "svelte/store";
    import { addNewDialog } from "../../utils/dialog.ts";

    export let isOpen = false;
    export let _class: string | undefined = "";
    export { _class as class };
    let settings = writable<CategorizedSettings | undefined>();
    let bodyClientHeight: number;
    let bodyInnerHeight: number;
    let bodyScrollTop: number = 0;
    let bodyContent: Element;
    let thisPage: Element;

    async function pullSettingsFromDb(): Promise<config.AppConfig | undefined> {
        try {
            const res = await GetConfig();
            return res;
        } catch (err) {
            console.error(`Could not get settings`, err);
            return undefined;
        }
    }

    async function pullDisplayValuesFromDb() {
        try {
            const res = await GetDisplayValues();
            return res as ExtendedSettings;
        } catch (err) {
            console.error(`Could not get display values`, err);
            return undefined;
        }
    }

    async function getSettings() {
        const basicSettings =
            (await pullSettingsFromDb()) as unknown as BasicSetting;
        const displayVals = await pullDisplayValuesFromDb();

        const result = joinDisplaySettings(basicSettings, displayVals!);

        //! Temporarily remove these keys. Automatic reporting is not implemented yet.
        delete result["Reports"]["ReportAutoEnabled"];
        delete result["Reports"]["ReportAutoAt"];

        settings.set(result);
    }

    /**
     * Read CategorizedSettings to a BasicSetting key-value dictionary. Used when sending the final state of changes
     */
    function readIntoBasicSetting(stng: CategorizedSettings | undefined) {
        if (!stng) return;
        const basicSettings: BasicSetting = {};

        Object.keys(stng).forEach((catKey) => {
            Object.keys(stng[catKey]).forEach((setKey) => {
                basicSettings[setKey] = stng[catKey][setKey].Value;
            });
        });

        return basicSettings;
    }

    /**
     * Converts number values to string, then returns the new string-only dictionary. Used before calling UpdateSettings,
     * which takes a dictionary with string-only values
     */
    function convertChangedSettingsToStr(settings: BasicSetting) {
        const stringOnly: { [key: string]: string } = {};
        Object.keys(settings).forEach((key) => {
            stringOnly[key] = settings[key].toString();
        });
        return stringOnly;
    }

    /**
     * Called by InputSwitch whenever a value changes. New value might be in `event.detail.changed`, or passed as-is as `string`
     */
    function changedEvent(
        event: string | CustomEvent,
        categoryKey: string,
        settingKey: any
    ) {
        let changed = typeof event === "string" ? event : event.detail.changed; // Check if new value is passed in event.detail.changed

        if (changed === undefined) return;

        if (typeof changed === "boolean") changed = changed ? "1" : "0"; // Convert to string

        // Update settings with the new value
        settings.update((prev) => {
            if (!prev) return prev;
            prev[categoryKey][settingKey].Value = changed;
            return prev;
        });
    }

    function onScroll(e: Event) {
        const target = e.target as HTMLDivElement;
        bodyScrollTop = target.scrollTop;
        bodyInnerHeight = target.scrollHeight;
    }

    const dispatch = createEventDispatcher();

    onMount(async () => {
        await getSettings();
    });

    function finishSetup() {
        const finishedSettings = readIntoBasicSetting(get(settings));
        if (!finishedSettings) {
            addNewDialog({
                title: "Error",
                description:
                    "Could not finish setting up Recap. Please report this as an issue on GitHub.",
            });

            return;
        }
        const stringKeyValDict = convertChangedSettingsToStr(finishedSettings);
        dispatch("finished", { settings: stringKeyValDict });
    }
</script>

{#if isOpen}
    <div
        bind:this={thisPage}
        class="{_class
            ? _class
            : ''} flex flex-col items-center justify-center transition-all"
    >
        <div
            on:scroll={onScroll}
            bind:clientHeight={bodyClientHeight}
            bind:this={bodyContent}
            transition:scale={{start: 0.9, opacity: 0, easing: expoOut, duration: 500 }}
            class="bg-white text-black w-[80%] h-[80%] rounded-xl p-6 overflow-y-scroll relative flex"
        >
            <div id="content" class="h-max w-full">
                <div class="pb-4 border-b-[1px] border-neutral-400 [&>p]:mb-2">
                    <div>
                        <h1 class="text-4xl lg:text-3xl font-extrabold mb-4">
                            Setup
                        </h1>
                    </div>

                    <p>Thank you for using Recap!</p>
                    <p>
                        A few settings require your attention before Recap
                        starts. Please review the default settings and configure
                        your AI models.
                    </p>
                    <p>
                        If you choose to use an API that requires an API key
                        (e.g. Gemini), please enter its API key.
                    </p>
                </div>

                {#await getSettings() then}
                    {#if $settings}
                        {#each Object.keys($settings) as cat}
                            <div class="flex flex-col">
                                <!-- Category title -->
                                <div
                                    class="flex flex-col font-bold text-xl first:mt-3 mt-8 mb-4"
                                >
                                    {cat}
                                </div>
                                <div class="flex flex-col">
                                    {#each Object.keys($settings[cat]) as set}
                                        <div
                                            class="border-b-[1px] last:border-b-0 border-neutral-300 mb-2"
                                        >
                                            <h3 class="text-xl">
                                                {$settings[cat][set]
                                                    .DisplayName}
                                            </h3>
                                            <p>
                                                {$settings[cat][set]
                                                    .Description}
                                            </p>

                                            <!-- Input section -->
                                            <div
                                                class="flex gap-2 items-center"
                                            >
                                                <!-- Input switch -->
                                                <InputSwitch
                                                    id={cat + "-" + set}
                                                    class="text-md my-4 pl-4 py-2 text-neutral-950"
                                                    inputType={$settings[cat][
                                                        set
                                                    ].InputType}
                                                    inputValue={$settings[cat][
                                                        set
                                                    ].Value}
                                                    inputOptions={$settings[cat][set].Options}
                                                    on:changed={(e) =>
                                                        changedEvent(
                                                            e,
                                                            cat,
                                                            set
                                                        )}
                                                ></InputSwitch>
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {/each}
                    {/if}
                {/await}
                <div class="flex w-full justify-end">
                    <button
                        class="h-full py-2 px-4 rounded-2xl text-2xl text-white bg-blue-400 hover:bg-blue-300"
                        on:click={finishSetup}
                    >
                        Complete setup
                    </button>
                </div>
            </div>
            <div class="w-[8px] relative -right-4">
                <VirtualScrollbar
                    bodyInner={bodyClientHeight}
                    bodyHeight={bodyInnerHeight}
                    bodyScroll={bodyScrollTop}
                    class="fixed w-2 -translate-x-1/2"
                    style="height: 75vh"
                ></VirtualScrollbar>
            </div>
        </div>
        <div
        transition:fade
            class="absolute top-0 left-0 w-full h-full bg-black opacity-50 -z-10"
        ></div>
    </div>
{/if}

<style global lang="postcss">
    @tailwind base;
    @tailwind components;
    @tailwind utilities;

    .page-container:has(> div.bypass-pad) {
        padding-top: 0px;
    }

    ::-webkit-scrollbar {
        display: none;
    }

    /* background of the scrollbar except button or resizer */
    ::-webkit-scrollbar-track {
        display: none;
    }

    /* scrollbar itself */
    ::-webkit-scrollbar-thumb {
        display: none;
    }

    /* set button(top and bottom of the scrollbar) */
    ::-webkit-scrollbar-button {
        display: none;
    }
</style>
