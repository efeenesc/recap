<script lang="ts">
    import gsap from "gsap";
    import { get, writable } from "svelte/store";
    import { db } from "$lib/wailsjs/go/models.ts";
    import type {
        BasicSetting,
        CategorizedSettings,
    } from "../../types/ExtendedSettings.interface.ts";
    import InputSwitch from "../../components/input-switch/InputSwitch.svelte";
    import { deepClone } from "../../utils/deepclone.ts";
    import { UpdateSettings } from "$lib/wailsjs/go/app/AppMethods.js";
    import { addNewDialog } from "../../utils/dialog.ts";
    import RevertIcon from "../../icons/RevertIcon.svelte";
    import { beforeNavigate, goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { scrollStore } from "$lib/stores/ScrollStore.ts";

    interface Data {
        streamed: {
            items: Promise<CategorizedSettings | undefined>;
        };
    }

    export let data: Data;
    let newSet = writable<CategorizedSettings | undefined>(undefined);
    let rcvSet: BasicSetting = {};
    let changedSettings: BasicSetting = {};
    let allowPageChange: boolean = false;
    let wereSettingsChanged: boolean = false;
    let scrollTop: number = 0;
    let titleBackgroundOpacity: boolean = false;

    $: {
        if (scrollTop) {
            titleBackgroundOpacity = scrollTop > 100 ? true : false;
        }
    }

    async function getData(): Promise<CategorizedSettings | undefined> {
        let items = undefined;
        try {
            items = await data.streamed.items;
        } catch (err) {
            console.error(err);
        } finally {
            newSet.set(deepClone(items));
            readIntoBasicSetting(items);
            setTimeout(() => animateLoadForAllDivs(), 100);
            subscribeToChanges();
            return items;
        }
    }

    /**
     * Read CategorizedSettings received from +load to a BasicSetting key-value dictionary. Simplifies checking for setting changes
     */
    function readIntoBasicSetting(stng: CategorizedSettings | undefined) {
        rcvSet = {};
        if (!stng) return;

        Object.keys(stng).forEach((catKey) => {
            Object.keys(stng[catKey]).forEach((setKey) => {
                rcvSet[setKey] = stng[catKey][setKey].Value;
            });
        });
    }

    /**
     * Load 0 to 100 opacity transition for all divs
     */
    function animateLoadForAllDivs() {
        gsap.to(".setting-div", {
            opacity: 1,
            scale: 1,
            duration: 1,
            ease: "expo.out",
        });
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
        newSet.update((prev) => {
            if (!prev) return prev;

            // Not checking for type equality - this is intentional to get true out of 10 == '10'
            if (rcvSet[settingKey] == changed) {
                if (changedSettings[settingKey]) {
                    delete changedSettings[settingKey];
                }
            } else {
                changedSettings[settingKey] = changed;
            }

            prev[categoryKey][settingKey].Value = changed;

            return prev;
        });
    }

    /**
     * Ran every time a setting was changed to check if the changedSettings dictionary has any values. Sets 'settingsChanged' with the result
     */
    function checkSettingChanges() {
        wereSettingsChanged = Object.keys(changedSettings).length > 0;
    }

    /**
     * Converts number values to string, then returns the new string-only dictionary. Used before calling UpdateSettings,
     * which takes a dictionary with string-only values
     */
    function convertChangedSettingsToStr() {
        const stringOnly: { [key: string]: string } = {};
        Object.keys(changedSettings).forEach((key) => {
            stringOnly[key] = changedSettings[key].toString();
        });
        return stringOnly;
    }

    function subscribeToChanges() {
        newSet.subscribe(() => {
            checkSettingChanges();
        });
    }

    /**
     * Save changes by writing new settings to database and refreshing BasicSettings with the new values
     */
    async function saveChanges() {
        try {
            await UpdateSettings(convertChangedSettingsToStr());
            readIntoBasicSetting(get(newSet));
            changedSettings = {};
            wereSettingsChanged = false;
        } catch (error: any) {
            addNewDialog({
                title: "Error",
                description: `Could not update settings. The following error was received: ${error}`,
            });
        }
    }

    function revertChanges(cat: string, set: string): any {
        const prevVal = rcvSet[set];

        changedEvent(prevVal.toString(), cat, set);
    }

    beforeNavigate(async ({ to, cancel }) => {
        if (allowPageChange) {
            // Skip the dialog if allowPageChange is true. This is set in the callback functions below
            // to stop beforeNavigate from interrupting the transition
            allowPageChange = false;
            return;
        }

        if (wereSettingsChanged) {
            cancel(); // Cancel the navigation immediately

            addNewDialog({
                title: "Continue?",
                description:
                    "You have unsaved changes. Would you like to continue without saving?",
                primaryButtonName: "Save",
                primaryButtonCallback: async () => {
                    await saveChanges();
                    allowPageChange = true;
                    goto(to!.url.pathname);
                },
                secondaryButtonName: "Don't save",
                secondaryButtonCallback: () => {
                    allowPageChange = true;
                    goto(to!.url.pathname);
                },
                tertiaryButtonName: "Cancel",
                tertiaryButtonCallback: () => {
                    // Do nothing, navigation is already cancelled
                },
            });
        }
    });

    onMount(() => {
        const unsubscribe = scrollStore.subscribe(
            (scrollVal) => (scrollTop = scrollVal)
        );

        return () => {
            unsubscribe();
        }; // Unsubscribe from this mistake of a store
    });
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-max min-h-screen inline bypass-pad">
    <div class="pb-2 flex flex-col">
        <div
            class="top-gradient-bg {titleBackgroundOpacity
                ? 'after:opacity-100'
                : 'after:opacity-0'} pt-10 sticky left-0 top-0 z-30 flex justify-between"
        >
            <h1
                class="page-title w-1/2 text-2xl -tracking-wide opacity-85 z-40"
            >
                Settings
            </h1>
            <div
                class="flex gap-2 justify-self-end -tracking-wide text-xl z-30 text-black font-semibold"
            >
                <div
                    on:click={saveChanges}
                    class="-tracking-wide transition-all {wereSettingsChanged
                        ? 'opacity-100'
                        : 'opacity-0'} text-nowrap text-xl px-4 p-2 bg-opacity-80 cursor-pointer active:scale-[99%] hover:bg-opacity-90 bg-blue-400 text-black font-semibold rounded-lg"
                >
                    Save changes
                </div>
            </div>
        </div>

        {#await getData()}
            Loading settings...
        {:then}
            {#if $newSet}
                {#each Object.keys($newSet) as cat}
                    <div class="flex flex-col">
                        <!-- Category title -->
                        <div class="flex flex-col top-16 sticky z-40">
                            <h1 class="category font-bold text-3xl mb-4">
                                {cat}
                            </h1>
                        </div>
                        <div class="flex flex-col">
                            {#each Object.keys($newSet[cat]) as set}
                                <div
                                    class="border-b-[1px] border-neutral-800 mb-2"
                                >
                                    <h3 class="text-xl">
                                        {$newSet[cat][set].DisplayName}
                                    </h3>
                                    <p>{$newSet[cat][set].Description}</p>

                                    <!-- Input section -->
                                    <div class="flex gap-2 items-center">
                                        <!-- Input switch -->
                                        <InputSwitch
                                            id={cat + "-" + set}
                                            class="text-md my-4 pl-4 py-2 text-neutral-950"
                                            inputType={$newSet[cat][set]
                                                .InputType}
                                            inputValue={$newSet[cat][set].Value}
                                            inputOptions={["ollama", "gemini"]}
                                            on:changed={(e) =>
                                                changedEvent(e, cat, set)}
                                        ></InputSwitch>

                                        <!-- Revert button -->
                                        <div
                                            on:click={() =>
                                                revertChanges(cat, set)}
                                            class="cursor-pointer inline flex-shrink transition-all p-1 w-8 h-8 aspect-square rounded-full bg-gray-200 {changedSettings[
                                                set
                                            ] !== undefined
                                                ? 'opacity-100 scale-100'
                                                : 'opacity-0 scale-[95%]'}"
                                        >
                                            <RevertIcon strokeColor="#2966eb"
                                            ></RevertIcon>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/each}
            {/if}
        {/await}
    </div>
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;

    .top-gradient-bg::after {
        display: block;
        content: "";
        position: absolute;
        top: 0;
        left: -100px;
        right: -100px;
        height: 200px;
        z-index: 1;
        background: linear-gradient(
            180deg,
            rgb(0, 0, 0) 0%,
            rgba(0, 0, 0, 0) 100%
        );
    }
</style>
