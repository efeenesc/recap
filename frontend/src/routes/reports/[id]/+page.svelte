<script lang="ts">
	import { ConvertToHtmlTree } from "$lib/markdown/Markdown.ts"
    import { goto } from "$app/navigation";
    import BackArrow from "../../../icons/BackArrow.svelte";
    import type { ExtendedReport } from "../../../types/ExtendedReport.interface.ts";
    import MarkdownRenderer from "../../../components/markdown-renderer/MarkdownRenderer.svelte";

    interface Data {
        streamed: {
            items: Promise<ExtendedReport | undefined>;
        };
    }

    function parseMd(content: string) {
        const parsed = ConvertToHtmlTree(content);
        return parsed;
    }

    /** @type {import('./$types').PageData} */
    export let data: Data;

    async function goBack() {
        window.history.back();
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
        <div
            class="flex gap-2 overflow-hidden z-30 pt-10 sticky top-0 h-min items-center"
        >
            <div
                on:click={goBack}
                class="aspect-square flex-shrink h-8 cursor-pointer"
            >
                <BackArrow strokeColor="#fff"></BackArrow>
            </div>
            <h1
                class="page-title text-2xl -tracking-wide opacity-85 left-0 top-0 w-1/2 z-20"
            >
                Reports
            </h1>
        </div>
        {#await data.streamed.items then report}
            {#if report !== undefined}
                <div class="w-full h-max inline bypass-pad">
                    <div class="pb-2 gap-5 flex flex-col pt-2">
                        <div
                            id="r{report.ReportID}"
                            class="transition-box-container max-h-max relative rounded-lg w-fit justify-center items-center bg-neutral-800 outline overflow-hidden outline-1 outline-neutral-900 p-1 shadow-2xl"
                        >
                            <div
                                class="absolute left-0 top-0 right-0 bottom-0 bg-neutral-800 opacity-0 hidden"
                            ></div>
                            <div
                                class="flex flex-col flex-shrink overflow-hidden p-2 bg-neutral-900 transition-all rounded-lg object-contain select-none pointer-events-none"
                            >
                                <div class="-mt-4">
                                    <MarkdownRenderer
                                        parsedContent={report.ParsedMarkdown}
                                    ></MarkdownRenderer>
                                </div>
                            </div>
                            <h3 class="flex-shrink-0 pl-2 py-1">
                                Generated at {report.Time}
                            </h3>
                            <!-- <div
                                class="transition-box-content flex flex-col rounded-xl object-contain select-none pointer-events-none py-6 px-3"
                            >
                                <div>
                                    <MarkdownRenderer parsedContent={parseMd(report.Content).content}></MarkdownRenderer>
                                </div>
                                <h3 class="flex-shrink-0 pl-2 py-1">
                                    {report.Time}
                                </h3>
                            </div> -->
                        </div>
                    </div>
                    <div class="flex flex-col justify-end items-end gap-5 mt-4">
                        <div class="flex justify-center items-center gap-1">
                            <p class="whitespace-pre">
                                Generated with <span
                                    class="whitespace-pre animated-text"
                                    >{report.GenWithModel}</span
                                >
                            </p>
                            <img
                                class="h-4 aspect-square"
                                alt="Gemini logo"
                                src="/Gemini_Logo.png"
                            />
                        </div>
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
        background: radial-gradient(
            circle at 100%,
            #1071ee,
            #ed30a2 50%,
            #807ae3 75%
        );
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
