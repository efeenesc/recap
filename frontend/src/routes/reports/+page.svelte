<script lang="ts">
    import gsap from "gsap";
    import type { DatedReport } from "../../types/ExtendedReport.interface.ts";
    import { writable } from "svelte/store";
    import { EventsOff, EventsOn } from "$lib/wailsjs/runtime/runtime.js";
    import { addReportToStore } from "../../utils/report.ts";
    import { onDestroy, onMount } from "svelte";
    import { afterNavigate, goto } from "$app/navigation";
    import Checkbox from "../../components/checkbox/Checkbox.svelte";
    import Checkmark from "../../icons/Checkmark.svelte";
    import { addNewDialog } from "../../utils/dialog.ts";
    import MarkdownRenderer from "../../components/markdown-renderer/MarkdownRenderer.svelte";
    import { lex, parse } from "$lib/markdown/MarkdownParser.ts";

    interface Data {
        streamed: {
            items: Promise<DatedReport | undefined>;
        };
    }

    let title: HTMLDivElement;
    let prevId: number | undefined;
    let selecting: boolean = false;
    const rcvRep = writable<DatedReport | undefined>();
    let checkedItems: { [key: string]: boolean | undefined } = {};
    export let data: Data;

    function parseMd(content: string) {
        const parsed = parse(lex(content));
        return parsed;
    }

    async function getData(): Promise<DatedReport | undefined> {
        try {
            const items = await data.streamed.items;
            rcvRep.set(items);
            setTimeout(() => animateLoadForAllDivs(), 100);
            return items;
        } catch (err) {
            console.error(err);
            rcvRep.set(undefined);
        }
    }

    function animateLoadForAllDivs() {
        gsap.to(".report-div", {
            opacity: 1,
            scale: 1,
            duration: 1,
            ease: "expo.out",
        });
    }

    function subscribeToReportEvent() {
        EventsOn("rcv:llmran", async (lastId) => {
            await addReportToStore(lastId);
            addNewDialog({
                title: "New item added",
                description: `New report ID: ${lastId}`,
                primaryButtonCallback: () => goto(`/reports/${lastId}`),
                primaryButtonName: "Open report",
            });
        });
    }

    onDestroy(() => {
        EventsOff("rcv:llmran");
    });

    onMount(() => {
        subscribeToReportEvent();
    });

    function selectAllFromDate(event: any, date: string) {
        let checkAll: boolean = checkedItems[date] ? false : true;
        checkedItems[date] = !checkAll;

        rcvRep.update((prev) => {
            if (!prev) return prev;

            return {
                ...prev,
                [date]: prev[date].map((s) => ({
                    ...s,
                    Selected: !checkAll,
                })),
            };
        });
    }

    afterNavigate((nav) => {
        if (!nav.from) return;
        const targetScr = document.getElementById("s" + nav.from.params)!;
        if (!targetScr) return;
        targetScr.classList.add("transition-box-container");
        targetScr.children[1].classList.add("transition-box-content");
    });

    function repClicked(date: string, reportId: number, event: MouseEvent) {
        // If 'selecting' is true and a screenshot was clicked, it was selected, not routed to. Early return
        if (selecting) {
            rcvRep.update((prev) => {
                if (!prev || !prev[date]) return prev;

                const updatedDate = prev[date].map((r) =>
                    r.ReportID === reportId
                        ? { ...r, Selected: !r.Selected }
                        : r
                );

                // Check if all reports are selected
                const allSelected = updatedDate.every((rep) => rep.Selected);

                // Update the checkedItems state
                checkedItems[date] = allSelected;

                return {
                    ...prev,
                    [date]: updatedDate,
                };
            });
            return;
        }

        // Handle non-selecting clicks
        const clickTarget = event.target as HTMLDivElement | HTMLImageElement;
        const container: HTMLDivElement = clickTarget.id
            ? (clickTarget as HTMLDivElement)
            : ((clickTarget as HTMLImageElement)
                  .parentElement as HTMLDivElement);

        const image: HTMLImageElement = container
            .children[0] as HTMLImageElement;

        // Add classes for styling
        container.classList.add("transition-box-container");
        image.classList.add("transition-box-content");

        // Navigate to the report
        goto(`/reports/${reportId}`);
    }

    function multiSelectClicked() {
        selecting ? multiSelectStop() : multiSelectStart();
    }

    function multiSelectStart() {
        const targets = document.getElementsByClassName("checkbox")!;
        selecting = true;

        gsap.to(targets, {
            opacity: 1,
            duration: 0.5,
            ease: "expo.out",
            onComplete: () => {
                Array.from(targets).forEach((el) => {
                    (el as HTMLDivElement).style.pointerEvents = "auto";
                });
            },
        });
    }

    function resetAllSelections() {
        rcvRep.update((prev) => {
            if (!prev) return prev;

            const updatedReports = Object.entries(prev).reduce(
                (acc, [date, reports]) => {
                    acc[date] = reports.map((rep) => ({
                        ...rep,
                        Selected: false,
                    }));
                    return acc;
                },
                {} as typeof prev
            );

            return updatedReports;
        });
        Object.keys(checkedItems).forEach((key) => {
            delete checkedItems[key];
        });
    }

    function multiSelectStop() {
        const targets = document.getElementsByClassName("checkbox")!;
        selecting = false;
        Array.from(targets).forEach((el) => {
            (el as HTMLDivElement).style.pointerEvents = "";
        });

        resetAllSelections();

        gsap.to(targets, {
            opacity: 0,
            duration: 0.5,
            ease: "expo.out",
        });
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-max min-h-screen inline bypass-pad">
    <div class="pb-2 gap-5 flex flex-col">
        <div class="pt-10 sticky left-0 top-0 z-30 flex">
            <h1
                bind:this={title}
                class="page-title text-2xl -tracking-wide opacity-85 w-full"
            >
                <span class="z-50">
                    Reports
                </span>
                
            </h1>
            <div
                class="flex gap-2 justify-self-end -tracking-wide text-xl z-30 text-black font-semibold"
            >
                {#if selecting}
                    <div
                        class="opacity-0 transition-all delay-200 {selecting
                            ? 'opacity-100'
                            : ''} -tracking-wide text-xl px-4 p-2 bg-opacity-80 cursor-pointer bg-red-400 text-black font-semibold rounded-lg"
                    >
                        Delete
                    </div>
                {/if}

                <div
                    on:click={multiSelectClicked}
                    class="-tracking-wide text-xl px-4 p-2 bg-opacity-80 cursor-pointer bg-white text-black font-semibold rounded-lg"
                >
                    {selecting ? "Cancel" : "Select"}
                </div>
            </div>
        </div>

        {#await getData()}
            Loading reports...
        {:then}
            {#if $rcvRep}
                {#each Object.entries($rcvRep) as [date, reports]}
                    <div>
                        <div
                            class="flex gap-5 sticky items-center top-[74px] z-20"
                        >
                            <h2 class="text-3xl font-bold tracking-wider">
                                {date}
                            </h2>
                            <div class="checkbox opacity-0 pointer-events-none">
                                <Checkbox
                                    id={date}
                                    bind:checked={checkedItems[date]}
                                    on:checked={(e) =>
                                        selectAllFromDate(e, date)}
                                ></Checkbox>
                            </div>
                        </div>

                        <div class="my-4 flex flex-col gap-4">
                            {#each reports as r (r.ReportID)}
                                <div
                                    on:click={(event) =>
                                        repClicked(date, r.ReportID, event)}
                                    class="m-0 p-0"
                                >
                                    <div
                                        id="r{r.ReportID}"
                                        class="group report-div max-h-[200px] flex flex-col cursor-pointer relative rounded-lg w-fit bg-neutral-800 outline overflow-hidden outline-1 outline-neutral-900 p-1 mr-5 shadow-2xl opacity-0 scale-95"
                                    >
                                        {#if selecting}
                                            <div
                                                class="flex items-center justify-center absolute z-30 transition-all bg-opacity-50 left-0 top-0 right-0 bottom-0 {r.Selected
                                                    ? 'bg-neutral-500'
                                                    : 'bg-neutral-800'}"
                                            >
                                                <div
                                                    class="transition-all opacity-0 scale-90 w-[50%] {r.Selected
                                                        ? 'opacity-90 scale-100'
                                                        : ''}"
                                                >
                                                    <Checkmark
                                                        strokeColor="#fff"
                                                    ></Checkmark>
                                                </div>
                                            </div>
                                        {/if}

                                        <div
                                            class="flex flex-col flex-shrink overflow-hidden p-2 bg-neutral-900 transition-all rounded-lg object-contain select-none pointer-events-none"
                                        >
                                            <div class="-mt-4">
                                                <MarkdownRenderer
                                                    parsedContent={parseMd(r.Content)
                                                        .content}
                                                ></MarkdownRenderer>
                                            </div>
                                        </div>
                                        <h3 class="flex-shrink-0 pl-2 py-1">
                                            {r.Time}
                                        </h3>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/each}
            {:else}
                No reports yet
            {/if}
        {:catch error}
            Error loading reports: {error.message}
        {/await}
    </div>
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
