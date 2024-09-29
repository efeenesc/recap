<script lang="ts">
    import { onMount } from "svelte";
    import ClipboardIcon from "../../icons/ClipboardIcon.svelte";
    import ScreenshotIcon from "../../icons/ScreenshotIcon.svelte";
    import SettingsIcon from "../../icons/SettingsIcon.svelte";
    import { CheckTimers, EmitStartStopScrTimer, EmitStartStopLLMTimer } from "$lib/wailsjs/go/app/AppMethods.js";
    import { page } from '$app/stores'; 
    import Toggle from "../checkbox/Toggle.svelte";
    import { goto } from '$app/navigation'
    
    let currentRoute: string;
    let scrTimer: boolean = true;
    let llmTimer: boolean = true;
    let sidePanelFull: boolean;
    let showOverlay: boolean;
    let windowWidth: number;
    const routes = [
        {
            path: "/",
            icon: ScreenshotIcon,
            title: "Home",
        },
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
        {
            path: "/settings/",
            icon: SettingsIcon,
            title: "Settings",
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
        scrTimer = !scrTimer;
    }

    async function startStopLLMTimer() {
        await EmitStartStopLLMTimer(!llmTimer);
        llmTimer = !llmTimer;
    }

    onMount(() => {
        window.addEventListener("resize", onResize);
        onResize();
    });

    function showOverlayClick() {
        showOverlay = !showOverlay;
    }

    function routeClicked(path: string) {
        goto(path);
    }

    getTimerStatus()
    $: currentRoute = $page.url.pathname;
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="side-panel-container">
    <div class="w-full h-screen flex flex-col items-center text-2xl font-light pt-5">
        {#if sidePanelFull}
            <div class="h-full w-[250px] flex flex-col px-4 gap-2">
                <div class="flex flex-col gap-2">
                    {#each routes as r}
                        <div on:click={() => routeClicked(r.path)} class="cursor-pointer px-2 py-2 {currentRoute === r.path ? "dark:bg-neutral-100 dark:text-neutral-950 rounded-lg" : "" }">
                            <div
                                class="w-fit h-min flex gap-2 items-center justify-center overflow-hidden relative"
                            >
                                <div class="sidepanel-icon w-8 h-8">
                                    <svelte:component this={r.icon} strokeColor={"#666666"}></svelte:component>
                                </div>
                                
                                <h2 class="text-lg tracking-wide">{r.title}</h2>
                            </div>
                        </div>
                    {/each}
                </div>
                
                
                <div class="block flex-grow min-h-1"></div>

                <div class="flex flex-col justify-center items-center px-2 py-4 dark:bg-neutral-100 dark:text-neutral-950 rounded-lg mb-6">
                    <div class="mb-4 subpixel-antialiased">Schedule</div>
                    <div class="flex flex-wrap justify-between gap-2">
                        <h2 class="text-lg tracking-wide">Screenshots </h2>
    
                        <Toggle checked={scrTimer} 
                        neu={true}
                        id="scrTimerToggle" class="w-fit"
                        on:checked={(e) =>startStopScrTimer()}></Toggle>
    
                        <h2 class="text-lg tracking-wide col-start-1 row-start-2">Processing </h2>
                        <Toggle checked={llmTimer} 
                        neu={true}
                        id="llmTimerToggle" class="w-fit"
                        on:checked={(e) => startStopLLMTimer()}></Toggle>
                    </div>
                </div>
                
            </div>
        {:else}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <!-- svelte-ignore a11y-no-static-element-interactions -->
            <div
                on:click={showOverlayClick}
                class="w-[100px] bg-neutral-100 h-24"
            ></div>
        {/if}
    </div>

    {#if showOverlay}
        <div
            class="fixed left-0 top-0 w-[300px] h-screen bg-neutral-100 outline outline-1 outline-white z-20"
        ></div>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <!-- svelte-ignore a11y-no-static-element-interactions -->
        <div
            on:click={showOverlayClick}
            class="fixed left-0 top-0 w-screen h-screen opacity-10 bg-black z-10"
        ></div>
    {/if}
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
