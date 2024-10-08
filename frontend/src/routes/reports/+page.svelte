<script lang="ts">
	import type { DatedReport } from './../../types/ExtendedReport.interface.ts';
	import { reportStore } from '$lib/stores/ReportStore.ts';
    import gsap from "gsap";
    import { get, writable } from "svelte/store";
    import { EventsOff, EventsOn } from "$lib/wailsjs/runtime/runtime.js";
    import { addReportToStore, getReportsNewerThan, getReportsOlderThan } from "../../utils/report.ts";
    import { onDestroy, onMount } from "svelte";
    import { afterNavigate, goto } from "$app/navigation";
    import Checkbox from "../../components/checkbox/Checkbox.svelte";
    import Checkmark from "../../icons/Checkmark.svelte";
    import { addNewDialog } from "../../utils/dialog.ts";
    import MarkdownRenderer from "../../components/markdown-renderer/MarkdownRenderer.svelte";
    import { ConvertToHtmlTree } from "$lib/markdown/Markdown.ts"
    import { scrollStore } from "$lib/stores/ScrollStore.ts";
    import { DeleteReportsById } from "$lib/wailsjs/go/app/AppMethods.js"
    import type { Snapshot } from '@sveltejs/kit';

    interface Data {
        data: DatedReport | undefined
    }

    export let data: Data;

    let selecting: boolean = false;
    const rcvRep = writable<DatedReport | undefined>();
    let checkedItems: { [key: string]: boolean | undefined } = {};
    let loadMoreDiv: Element;
    let loadMoreDivObserver: IntersectionObserver;
    let loadMoreDivObserverTimeout: number | undefined;
    
    let scrollTop: number = 0;

    // Variable assigned to during snapshot.recover. If assigned, scroll this element to previous Y position. 
    let scrollTopSnapshot: number | undefined;
    let titleBackgroundOpacity: boolean = false;

    let allReportsLoaded: boolean = false;
    let noReports: boolean = false;

    $: rcvRep.set(data.data);

    $: {
        if (scrollTop) {
            titleBackgroundOpacity = scrollTop > 100 ? true : false;
        }
    }

    $: {
        if (scrollTopSnapshot) {
            const scroller = document.getElementsByClassName('scroller')[0];
            setTimeout(() => {
                scroller.scroll(0, scrollTopSnapshot!)
            }, 0)
        }
    }

    export const snapshot: Snapshot<number> = {
        capture: () => {
            return scrollTop
        },
        restore: (snapshot) => {
            scrollTopSnapshot = snapshot;
        },
    };

    function animateLoad(id: string) {
        gsap.fromTo('#' + id, {
            opacity: 0,
            scale: 0.9
        }, {
            opacity: 1,
            scale: 1,
            duration: 1,
            ease: "expo.out",
        });
    }

    async function addNewReports() {
        let newestKnown = 0;
        const allScr = get(rcvRep);

        if (allScr) {
            const newestKey = Object.keys(allScr)[0];
            allScr[newestKey].forEach((scr) => {
                if (scr.ReportID > newestKnown) newestKnown = scr.ReportID;
            });
        }

        const newReports = await getReportsNewerThan(newestKnown);
        if (!newReports) return;

        reportStore.update((prev) => {
            if (!prev) return prev;

            return [...newReports, ...prev];
        });
    }

    function subscribeToReportEvent() {
        EventsOn("rcv:llmran", async () => {
            await addNewReports();
        });
    }

    async function getOlderReports() {
        let oldestKnownId: number | undefined;
        const allReports = get(rcvRep);

        if (allReports && typeof allReports === "object") {
            const reportArrays = Object.values(allReports);
            if (reportArrays.length > 0) {
                const flatScreenshots = reportArrays.flat();
                oldestKnownId = Math.min(
                    ...flatScreenshots.map((scr) => scr.ReportID)
                );
            }
        }

        if (oldestKnownId === undefined) {
            console.log("No existing reports found");
            allReportsLoaded = false;
            return;
        }

        try {
            const newReports = await getReportsOlderThan(
                oldestKnownId,
                30
            );

            if (newReports === null || newReports.length === 0) {
                return;
            }

            reportStore.update((prev) => {
                if (!Array.isArray(prev)) return newReports;
                return [...prev, ...newReports];
            });
        } catch (error) {
            console.error("Error fetching older screenshots:", error);
        }
    }

    $: if (loadMoreDiv) {
        loadMoreDivObserverTimeout = setTimeout(() => {
            loadMoreDivObserver.observe(loadMoreDiv);
        }, 100)
    }

    onMount(() => {
        const unsubscribe = scrollStore.subscribe(scrollPos => scrollTop = scrollPos);
        subscribeToReportEvent();
        loadMoreDivObserver = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        getOlderReports();
                    }
                });
            },
            { threshold: 0.1 }
        ); // Trigger when 10% of the element is visible

        return () => {
            unsubscribe();
            loadMoreDivObserver.disconnect();
            clearTimeout(loadMoreDivObserverTimeout);
            EventsOff("rcv:llmran");
        }
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

    /**
     * Ask the user if they're sure they want to delete N number of selected screenshots. If they proceed, call deleteSelectedStep2 
     */
    async function deleteSelectedStep1() {
        const allScrs = get(rcvRep);
        if (!allScrs) return;

        const selectedIds: number[] = [];
        Object.keys(allScrs).forEach((key) =>
            allScrs[key].map((scr) => {
                if (scr.Selected === true) {
                    selectedIds.push(scr.ReportID);
                }
            })
        );

        try {
            addNewDialog({
                title: "Confirmation",
                description: `${selectedIds.length} report${selectedIds.length !== 1 ? 's' : ''} will be deleted. Are you sure you want to proceed?`,
                primaryButtonName: "Delete",
                primaryButtonCallback: async () => await deleteSelectedStep2(selectedIds),
                secondaryButtonName: "Cancel",
                secondaryButtonCallback: () => console.log("Cancelled")
            });
        } catch (err) {
            console.error(err);
            addNewDialog({
                title: "Report generation failed",
                description: `The following error was received: ${err}`,
            });
        }
    }

    async function deleteSelectedStep2(ids: number[]) {
        try {
            await DeleteReportsById(ids);
            reportStore.update(prev => {
                const filtered = prev.filter(r => !ids.includes(r.ReportID))
                return filtered
            })
        } catch (err: any) {
            addNewDialog({
                title: "Error",
                description: `Could not delete screenshots with the following error: ${err}`
            })
        }
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-max min-h-screen inline bypass-pad">
    <div class="pb-2 flex flex-col">
        <div class="top-gradient-bg { titleBackgroundOpacity ? "after:opacity-100" : "after:opacity-0" } pt-10 sticky left-0 top-0 z-30 flex justify-between">
            <h1
                class="page-title text-2xl -tracking-wide opacity-85 w-1/2 z-40 pointer-events-none"
            >
                Reports
            </h1>
            <div
                class="flex gap-2 justify-self-end -tracking-wide text-xl z-30 text-black font-semibold"
            >
                {#if selecting}
                    <div
                        on:click={deleteSelectedStep1}
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

        {#if $rcvRep !== undefined}
            {#each Object.entries($rcvRep) as [date, reports]}
                <div>
                    <div
                        class="flex gap-5 sticky items-center top-16 z-40 pointer-events-none"
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
                                data-intersect
                                on:intersect={(e) => {r.Visible = e.detail.isIntersecting; animateLoad("r" + r.ReportID); }}
                                class="m-0 p-0 {r.Visible ? '' : 'invisible'}"
                            >
                                <div
                                    id="r{r.ReportID}"
                                    class="group report-div max-h-[200px] flex flex-col cursor-pointer relative rounded-lg w-fit bg-neutral-800 outline overflow-hidden outline-1 outline-neutral-900 p-1 mr-5 shadow-2xl"
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
                                        class="group-hover:scale-[99%] group-active:scale-[95%] flex flex-col flex-shrink overflow-hidden p-2 bg-neutral-900 transition-all rounded-lg object-contain select-none pointer-events-none"
                                    >
                                        <div class="-mt-4">
                                            <MarkdownRenderer
                                                parsedContent={r.ParsedMarkdown}
                                            ></MarkdownRenderer>
                                        </div>
                                    </div>
                                    <h3 class="flex-shrink-0 pl-2 py-1">
                                        Generated at {r.Time}
                                    </h3>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/each}
            {#if !allReportsLoaded}
                <div class="h-10 w-full" bind:this={loadMoreDiv}></div>
            {:else}
                <div></div>
            {/if}
        {:else}
            No reports yet. Create a report from your screenshots to get started!
        {/if}
    </div>
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;

    .top-gradient-bg::after {
        display: block;
        content: '';
        position: absolute;
        top: 0;
        left: -100px;
        right: -100px;
        height: 200px;
        z-index: 1;
        pointer-events: none;
        background: linear-gradient(180deg, rgb(0, 0, 0) 0%, rgba(0, 0, 0, 0) 100%);
    }
</style>
