<script lang="ts">
    import XMark from "../../../icons/XMark.svelte";
    import type { ExtendedScreenshot } from "../../../types/ExtendedScreenshot.interface.ts";
    import { beforeNavigate, goto } from "$app/navigation";
    import { onMount } from "svelte";
    import BackArrow from "../../../icons/BackArrow.svelte";

    interface Data {
        streamed: {
            items: Promise<ExtendedScreenshot | undefined>;
        };
    }

    /** @type {import('./$types').PageData} */
    export let data: Data;

    async function goBack() {
        await goto('/screenshots/');
    }

    function keyUp(key: KeyboardEvent) {
        if (key.code === "Escape")
            goBack();
    }
</script>

<svelte:window on:keyup={keyUp}></svelte:window>
<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-max min-h-screen inline bypass-pad">
    <div class="pb-2 gap-5 flex flex-col">
        <div class="flex gap-2 overflow-hidden z-30 pt-10 sticky top-0 h-min items-center">
            <div on:click={goBack} class="aspect-square flex-shrink h-8 cursor-pointer">
                <BackArrow strokeColor="#fff"></BackArrow>
            </div>
            <h1
                class="page-title text-2xl -tracking-wide opacity-85 left-0 top-0 w-full z-20"
            >
                Screenshots
            </h1>
        </div>
        {#await data.streamed.items then screenshot}
            {#if screenshot !== undefined}
                <div class="w-full h-max inline bypass-pad">
                    <div class="pb-2 gap-5 flex flex-col pt-2">
                        <div
                            id="s{screenshot.CaptureID}"
                            class="transition-box-container relative rounded-lg w-fit justify-center items-center bg-neutral-800 outline overflow-hidden outline-1 outline-neutral-900 p-1 shadow-2xl"
                        >
                            <div class="absolute left-0 top-0 right-0 bottom-0 bg-neutral-800 opacity-0 hidden">

                            </div>
                            <img
                                alt="screenshot"
                                class="transition-box-content max-h-[80vh] flex rounded-md object-contain select-none pointer-events-none"
                                loading="lazy"
                                src={screenshot.Screenshot}
                            />
                        </div>
                    </div>
                    <div class="flex flex-col justify-end items-end gap-5 mt-4">
                        {#if screenshot.Description !== null}
                            <p class="max-w-[50vw] bg-neutral-800 outline outline-1 outline-neutral-900 rounded-lg p-4 text-justify text-wrap">{screenshot.Description}</p>
                            <div class="flex justify-center items-center gap-1">
                                <p class="whitespace-pre">Generated with <span class="whitespace-pre animated-text">Gemini 1.5 Flash</span></p>
                                <img class="h-4 aspect-square" alt="Gemini logo" src="/Gemini_Logo.png">
                            </div>
                        {:else}
                            <p>A description hasn't been generated yet. Generate one now!</p>
                        {/if}
                    </div>
                </div>
            {:else}
                <p>Error: undefined</p>
            {/if}
        {:catch error}
            <p>Error: {error.message}</p>
        {/await}
    </div>
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;

    .animated-text {
        background: radial-gradient(circle at 100%, #1071ee, #ed30a2 50%, #807ae3 75%);
        font-weight: 600;
        background-size: 200% auto;
        color: #000;
        background-clip: text;
        -webkit-text-fill-color: transparent;
        animation: animatedTextGradient 1.5s linear infinite;
    }

    @keyframes animatedTextGradient {
        to {
            background-position: 200% center;
        }
    }
</style>
