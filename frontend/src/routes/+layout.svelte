<script lang="ts">
    import VirtualScrollbar from "./../components/virtual-scrollbar/VirtualScrollbar.svelte";
    import SidePanel from "./../components/side-panel/SidePanel.svelte";
    import { onNavigate } from "$app/navigation";
    import { onMount } from "svelte";
    import Dialog from "../components/dialog/Dialog.svelte";
    import { updateScroll } from "$lib/stores/ScrollStore.ts";
    import { createLazyIntersect } from "../components/lazy-intersect/LazyIntersect.ts";
    import { ReadInfo, UpdateSettings, UpdateInfo } from "$lib/wailsjs/go/app/AppMethods.js"
    import FirstTimeSetup from "../components/first-time-setup/FirstTimeSetup.svelte";
    import type { BasicSetting } from "../types/ExtendedSettings.interface.ts";

    let bodyFullHeight: number;
    let scrollHeight: number;
    let bodyContent: HTMLDivElement;
    let bodyInnerHeight: number;
    let showFirstTimeSetup: boolean = false;

    async function checkFirstTimeSetup() {
        const firstTimeSetupDone = (await ReadInfo("FirstTimeTutorialShown")).Value;
        if (firstTimeSetupDone != 1)
            showFirstTimeSetup = true;
    }

    async function firstTimeSetupFinished(ev: { detail: any }) {
        await UpdateSettings(ev.detail.settings);
        UpdateInfo({"FirstTimeTutorialShown": "1"}).then(() => {
            showFirstTimeSetup = false;
        });
    }

    onMount(() => {
        bodyContent.addEventListener("scroll", (ev) => {
            const target = (ev.target as HTMLDivElement);
            scrollHeight = target.scrollTop;
            bodyFullHeight = target.scrollHeight;
        });

        const { intersectionObserver, mutationObserver } = createLazyIntersect();
        checkFirstTimeSetup();

        return () => {
            mutationObserver.disconnect();
            intersectionObserver.disconnect();
        }
    });

    onNavigate((navigation) => {
        if (!document.startViewTransition) return;

        updateScroll(0);

        return new Promise((resolve) => {
            document.startViewTransition(async () => {
                
                resolve();
                await navigation.complete;
            });
        });
    });

    $: updateScroll(scrollHeight);
</script>

<main
    class="main-container h-screen w-full inter bg-opacity-30 select-none text-[#1f1f1f] dark:text-white"
>
    <!-- <div style="--wails-draggable:drag" class="z-50 fixed top-0 left-0 w-screen h-[6vh] flex justify-end">
        <div class="w-[40px] h-[30px] flex items-center justify-center">-</div>
    </div> -->
    <div class="h-full w-full flex flex-row overflow-hidden pt-2 relative">
        <div class="w-fit h-full grid grid-flow-row">
            <SidePanel></SidePanel>
        </div>
        <FirstTimeSetup isOpen={showFirstTimeSetup} on:finished={firstTimeSetupFinished} class="fixed top-0 left-0 w-screen h-screen z-50"></FirstTimeSetup>
        <Dialog class="z-50"></Dialog>
        <div
            bind:this={bodyContent}
            bind:clientHeight={bodyInnerHeight}
            class="scroller w-full h-full rounded-l-lg outline outline-1 bg-opacity-90 outline-neutral-200 dark:outline-neutral-900 px-4 lg:px-6 xl:px-10 overflow-hidden overflow-y-scroll dark:bg-opacity-50 bg-white dark:bg-black"
        >
            <!-- <div id="scroll-gradient" class="absolute h-[40%] top-2 left-[25%] right-0 z-20 opacity-0"></div> -->
            <div class="page-container h-max pt-10">
                <slot></slot>
            </div>
        </div>
        <VirtualScrollbar
            class="fixed h-[95vh] top-[1vh] right-2 w-2"
            bodyInner={bodyInnerHeight}
            bodyHeight={bodyFullHeight}
            bodyScroll={scrollHeight}
        ></VirtualScrollbar>
    </div>
</main>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;

    /* #scroll-gradient {
        background: linear-gradient(180deg, rgba(0,0,0,1) 0%, rgba(255,255,255,0) 100%);
    } */

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
