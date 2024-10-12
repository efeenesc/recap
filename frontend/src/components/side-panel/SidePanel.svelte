<script lang="ts">
    import InfoIcon from "./../../icons/InfoIcon.svelte";
    import { onDestroy, onMount } from "svelte";
    import ClipboardIcon from "../../icons/ClipboardIcon.svelte";
    import ScreenshotIcon from "../../icons/ScreenshotIcon.svelte";
    import SettingsIcon from "../../icons/SettingsIcon.svelte";
    import {
        CheckTimers,
        EmitStartStopScrTimer,
        EmitStartStopLLMTimer,
    } from "$lib/wailsjs/go/app/AppMethods.js";
    import { page } from "$app/stores";
    import Toggle from "../checkbox/Toggle.svelte";
    import { goto } from "$app/navigation";
    import HomeIcon from "../../icons/HomeIcon.svelte";
    import { EventsOff, EventsOn } from "$lib/wailsjs/runtime/runtime.js";

    interface Route {
        path: string;
        icon: any;
        title: string;
    }
    interface RouteGroup {
        title?: string;
        routes: Route[];
    }

    let currentRoute: string;
    let scrTimer: boolean = true;
    let llmTimer: boolean = true;
    let sidePanelFull: boolean;
    let windowWidth: number;
    const routes: RouteGroup[] = [
        {
            title: "Overview",
            routes: [
                {
                    path: "/",
                    icon: HomeIcon,
                    title: "Home",
                },
            ],
        },
        {
            title: "Spaces",
            routes: [
                {
                    path: "/screenshots/",
                    icon: ScreenshotIcon,
                    title: "Screenshots",
                },
                {
                    path: "/reports/",
                    icon: ClipboardIcon,
                    title: "Reports",
                },
            ],
        },
        {
            title: "System",
            routes: [
                {
                    path: "/settings/",
                    icon: SettingsIcon,
                    title: "Settings",
                },
                {
                    path: "/about/",
                    icon: InfoIcon,
                    title: "About",
                },
            ],
        },
    ];

    function onResize() {
        windowWidth = window.innerWidth;
        sidePanelFull = windowWidth > 900;
    }

    async function getTimerStatus() {
        const { scr, llm } = await CheckTimers();
        scrTimer = scr;
        llmTimer = llm;
    }

    async function startStopScrTimer() {
        await EmitStartStopScrTimer(!scrTimer);
    }

    async function startStopLLMTimer() {
        await EmitStartStopLLMTimer(!llmTimer);
    }

    onMount(() => {
        window.addEventListener("resize", onResize);
        EventsOn("rcv:llmstate", (newState: boolean) => {
            llmTimer = newState;
        });
        EventsOn("rcv:screenshotstate", (newState: boolean) => {
            scrTimer = newState;
        });
        onResize();

        return () => {
            EventsOff("rcv:llmstate", "rcv:screenshotstate");
        };
    });

    function routeClicked(path: string) {
        goto(path);
    }

    getTimerStatus();
    $: currentRoute = $page.url.pathname;
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="side-panel-container">
    <div
        class="w-full h-screen flex flex-col items-center text-2xl font-light pt-5"
    >
        {#if sidePanelFull}
            <div class="h-full w-[250px] flex flex-col px-4 gap-2">
                <div class="flex flex-col gap-2">
                    {#each routes as rg}
                        <div class="flex flex-col gap-2">
                            {#if rg.title}
                                <div class="px-2 text-sm dark:text-neutral-400">
                                    {rg.title}
                                </div>
                            {/if}

                            {#each rg.routes as r}
                                <div
                                    on:click={() => routeClicked(r.path)}
                                    class="cursor-pointer px-2 py-2 {currentRoute ===
                                    r.path
                                        ? 'fill-neutral-200 stroke-neutral-200 dark:fill-[#666666] dark:stroke-[#666666] bg-neutral-500 text-neutral-100 dark:bg-neutral-100 dark:text-neutral-950 rounded-lg'
                                        : 'fill-neutral-800 stroke-neutral-800 dark:fill-[#666666] dark:stroke-[#666666] text-neutral-800 dark:text-white'}"
                                >
                                    <div
                                        class="w-fit h-min flex gap-2 items-center justify-center overflow-hidden relative"
                                    >
                                        <div class="sidepanel-icon w-8 h-8">
                                            <svelte:component
                                                this={r.icon}
                                            ></svelte:component>
                                        </div>

                                        <h2 class="text-lg tracking-wide">
                                            {r.title}
                                        </h2>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {/each}
                </div>

                <div class="block flex-grow min-h-1"></div>

                <div
                    class="flex flex-col justify-center px-2 py-1 rounded-lg mb-6"
                >
                    <div class="mb-2 text-base text-neutral-400">Schedule</div>
                    <div class="flex flex-wrap justify-between gap-2">
                        <h2 class="text-lg tracking-wide text-black dark:text-white">Screenshots</h2>

                        <Toggle
                            checked={scrTimer}
                            id="scrTimerToggle"
                            class="w-fit"
                            on:checked={(e) => startStopScrTimer()}
                        ></Toggle>

                        <h2
                            class="text-lg tracking-wide col-start-1 row-start-2"
                        >
                            Processing
                        </h2>
                        <Toggle
                            checked={llmTimer}
                            id="llmTimerToggle"
                            class="w-fit"
                            on:checked={(e) => startStopLLMTimer()}
                        ></Toggle>
                    </div>
                </div>
            </div>
        {:else}
            <div class="flex flex-col gap-2">
                {#each routes as rg}
                    <!-- {#if rg.title}
                        <div class="mb-2 px-2 text-base dark:text-neutral-400">
                            {rg.title}
                        </div>
                    {/if} -->
                    {#each rg.routes as r}
                        <div
                            on:click={() => routeClicked(r.path)}
                            class="mx-1 cursor-pointer px-2 py-2 {currentRoute ===
                            r.path
                                ? 'dark:bg-neutral-100 dark:text-neutral-950 rounded-lg'
                                : ''}"
                        >
                            <div
                                title={r.title}
                                class="w-fit h-min flex gap-2 items-center justify-center overflow-hidden relative"
                            >
                                <div class="sidepanel-icon w-8 h-8 fill-[#1f1f1f] dark:fill-[#666666]">
                                    <svelte:component
                                        this={r.icon}
                                    ></svelte:component>
                                </div>
                            </div>
                        </div>
                    {/each}
                {/each}
            </div>
        {/if}
    </div>
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
