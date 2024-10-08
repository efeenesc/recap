<script lang="ts">
    import { db } from "./../lib/wailsjs/go/models.ts";
    import Carousel from "../components/carousel/Carousel.svelte";
    import {
        GetScreenshots,
        GetReports,
    } from "../lib/wailsjs/go/app/AppMethods.js";
    import { EventsOn } from "../lib/wailsjs/runtime/runtime.js";
    import { timeSinceUNIXSeconds } from "../utils/timeSince.js";
    import gsap from "gsap";
    import MarkdownRenderer from "../components/markdown-renderer/MarkdownRenderer.svelte";
    import type { ExtendedReport } from "../types/ExtendedReport.interface.ts";
    import type { ExtendedScreenshot } from "../types/ExtendedScreenshot.interface.ts";
    import { onMount } from "svelte";
    import { ConvertToHtmlTree } from "$lib/markdown/Markdown.ts";

    // EventsOn(
    //     "rcv:greet",
    //     (msg) => ()
    // );

    let screenshots: ExtendedScreenshot[] = [];
    let reports: ExtendedReport[] = [];
    let noScreenshots: boolean = false;
    let noReports: boolean = false;

    async function getScreenshots() {
        const res = await GetScreenshots(10);
        if (!res) {
            noScreenshots = true;
            return;
        }
        noScreenshots = false;
        screenshots = res.map((s: any) => {
            s.Date = timeSinceUNIXSeconds(s.Timestamp);
            return s;
        });
    }

    async function getReports() {
        const res = await GetReports(10);
        if (!res) {
            noReports = true;
            return;
        }
        noReports = false;
        reports = res.map((r: any) => {
            r.Date = timeSinceUNIXSeconds(r.Timestamp);
            r.ParsedMarkdown = parseMd(r.Content).content;
            return r;
        });
    }

    async function animateLoad(id: string) {
        gsap.to("#" + id, {
            opacity: 1,
            scale: 1,
            duration: 1,
            ease: "expo.out",
        });
    }

    function parseMd(content: string) {
        return ConvertToHtmlTree(content);
    }

    onMount(() => {
        getScreenshots();
        getReports();
    });
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-full bypass-padding">
    <div class="pb-2 gap-5 flex items-end">
        <h1
            class="text-2xl w-fit -tracking-wide opacity-85 text-[#1f1f1f] dark:text-white"
        >
            Home
        </h1>
    </div>

    <section class="relative">
        <h2 class="text-3xl font-bold tracking-wider">My screenshots</h2>
        <div class="my-4 h-[300px] flex">
            <div class="absolute" style="display: unset">
                {#if !noScreenshots}
                    <Carousel let:onLoad>
                        {#each screenshots as s (s.CaptureID)}
                            <div
                                id="s{s.CaptureID}"
                                class="rounded-lg h-[300px] bg-neutral-200 outline-neutral-300 dark:bg-neutral-800 dark:outline-neutral-900 
                                outline w-max overflow-hidden outline-1 p-1 mr-5 shadow-2xl opacity-0 scale-95"
                            >
                                <img
                                    alt="screenshot"
                                    loading="lazy"
                                    on:load|once={() => {
                                        onLoad();
                                        animateLoad("s" + s.CaptureID);
                                    }}
                                    class="flex rounded-md h-[90%] flex-shrink object-contain select-none pointer-events-none"
                                    src={s.Screenshot}
                                />
                                <h3 class="flex-shrink-0 pl-2 py-1 self-end">
                                    Snapped {s.Date}
                                </h3>
                            </div>
                        {/each}
                    </Carousel>
                {:else}
                    <div class="flex flex-grow flex-col w-max h-[300px]">
                        <h1 class="text-2xl">
                            Your screenshots will appear here.
                        </h1>
                        <h1 class="text-2xl">
                            Turn scheduled screenshots on to have the app get
                            them for you 📸
                        </h1>
                    </div>
                {/if}
            </div>
        </div>
    </section>

    <section class="w-full pt-6">
        <h2 class="text-3xl font-bold tracking-wider">My reports</h2>
        <div class="w-full my-4 h-[300px] flex">
            <div class="absolute w-max">
                {#if !noReports}
                    <Carousel let:onLoad>
                        {#each reports as r (r.ReportID)}
                            <div
                                id="s{r.ReportID}"
                                class="max-h-[300px] flex flex-col max-w-[400px] relative rounded-lg w-fit bg-neutral-800 outline overflow-hidden outline-1 outline-neutral-900 p-1 mr-5 shadow-2xl"
                            >
                                <div
                                    on:load|once={() => {
                                        onLoad();
                                        animateLoad("s" + r.ReportID);
                                    }}
                                    class="flex flex-col flex-shrink overflow-hidden p-2 bg-neutral-900 transition-all rounded-lg object-contain select-none pointer-events-none"
                                >
                                    <div class="-mt-4">
                                        <MarkdownRenderer
                                            parsedContent={r.ParsedMarkdown}
                                        ></MarkdownRenderer>
                                    </div>
                                </div>
                                <h3 class="flex-shrink-0 pl-2 py-1">
                                    Generated {r.Date}
                                </h3>
                            </div>
                        {/each}
                    </Carousel>
                {:else}
                    <div class="flex flex-grow flex-col w-max h-[300px]">
                        <h1 class="text-2xl">
                            Your generated reports will appear here.
                        </h1>
                        <h1 class="text-2xl">
                            Go to the screenshots page, select your screenshots, and click
                            'Report'! 🪄
                        </h1>
                    </div>
                {/if}
            </div>
        </div>
    </section>
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
