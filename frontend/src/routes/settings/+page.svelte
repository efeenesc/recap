<script lang="ts">
    import gsap from "gsap";
    import { writable } from "svelte/store";
    import { db } from "$lib/wailsjs/go/models.ts";
    import type { CategorizedSettings } from "../../types/ExtendedSettings.interface.ts";
    import InputSwitch from "../../components/input-switch/InputSwitch.svelte";

    interface Data {
        streamed: {
            items: Promise<CategorizedSettings | undefined>;
        };
    }

    export let data: Data;
    let rcvSet = writable<CategorizedSettings | undefined>(undefined);
    let newSet = writable<CategorizedSettings | undefined>(undefined);

    async function getData(): Promise<CategorizedSettings | undefined> {
        try {
            const items = await data.streamed.items;
            rcvSet.set(items);
            newSet.set(items);
            console.log(items);
            setTimeout(() => animateLoadForAllDivs(), 100);
            return items;
        } catch (err) {
            console.error(err);
            rcvSet.set(undefined);
            newSet.set(undefined);
        }
    }

    function animateLoadForAllDivs() {
        gsap.to(".setting-div", {
            opacity: 1,
            scale: 1,
            duration: 1,
            ease: "expo.out",
        });
    }

    function changedEvent(event: CustomEvent, categoryKey: string, settingKey: any) {
        const changed = event.detail.changed;
        if (!changed) return;
        
        newSet.update(prev => {
            if (!prev) return prev;

            prev[categoryKey][settingKey].Value = changed;
            return prev;
        })
        // console.log("Received event:", changed, categoryKey, settingKey);
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-max min-h-screen inline bypass-pad">
    <div class="pb-2 gap-5 flex flex-col">
        <div class="pt-10 sticky left-0 top-0 z-30 flex">
            <h1 class="page-title text-2xl -tracking-wide opacity-85 w-full">
                Settings
            </h1>
            <div
                class="flex gap-2 justify-self-end -tracking-wide text-xl z-30 text-black font-semibold"
            ></div>
        </div>

        {#await getData()}
            Loading settings...
        {:then}
            {#if $newSet}
                {#each Object.keys($newSet) as cat}
                    <div class="flex flex-col my-4">
                        <div class="flex flex-col top-16 sticky">
                            <h1 class="category font-bold text-3xl mb-4">
                                {cat}
                            </h1>
                        </div>
                        <div class="flex flex-col">
                            {#each Object.keys($newSet[cat]) as set}
                                <div>
                                    <h3 class="text-xl">
                                        {$newSet[cat][set].DisplayName}
                                    </h3>
                                    <p>{$newSet[cat][set].Description}</p>
                                    <InputSwitch
                                        id={cat + '-' + set}
                                        class="text-md my-4 pl-4 py-2 text-neutral-950"
                                        inputType={$newSet[cat][set].InputType}
                                        inputValue={$newSet[cat][set].Value}
                                        on:changed={(e) => changedEvent(e, cat, set)}
                                    ></InputSwitch>
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
</style>
