<script lang="ts">
	import { addNewDialog } from './../../utils/dialog.ts';
    import gsap from "gsap";
    import { onMount } from "svelte";
    import type { DatedScreenshot } from "../../types/ExtendedScreenshot.interface.ts";
    import { goto, afterNavigate, beforeNavigate } from "$app/navigation";
    import Checkbox from "../../components/checkbox/Checkbox.svelte";
    import { get, writable } from "svelte/store";
    import Checkmark from "../../icons/Checkmark.svelte";
    import { EventsOff, EventsOn } from "$lib/wailsjs/runtime/runtime.js";
    import {
        getScreenshotsNewerThan,
        getScreenshotsOlderThan,
    } from "../../utils/screenshot.ts";
    import {
        GenerateReportFromScreenshotIds,
        DeleteScreenshotsById,
    } from "$lib/wailsjs/go/app/AppMethods.js";
    import {
        processedScreenshotStore,
        screenshotStore,
    } from "$lib/stores/ScreenshotStore.ts";
    import { scrollStore } from "$lib/stores/ScrollStore.ts";
    import DoneAllIcon from "../../icons/DoneAllIcon.svelte";
    import DoneIcon from "../../icons/DoneIcon.svelte";
    import XIcon from "./../../icons/XIcon.svelte";
    import type { Snapshot } from "@sveltejs/kit";

    interface Data {
        data: DatedScreenshot | undefined;
    }

    export let data: Data;

    let loadMoreDiv: Element;
    let loadMoreDivObserver: IntersectionObserver;
    let loadMoreDivObserverTimeout: number;
    let prevId: number | undefined;
    let selecting: boolean = false;
    const rcvScr = writable<DatedScreenshot | undefined>();
    let checkedItems: { [key: string]: boolean | undefined } = {};
    let isInitialLoad = true;

    $: rcvScr.set(data.data);

    let scrollTop: number = 0;

    // Variable assigned to during snapshot.recover. If assigned, scroll this element to previous Y position.
    let scrollTopSnapshot: number | undefined;
    let titleBackgroundOpacity: boolean = false;

    let allScreenshotsLoaded: boolean = false;

    $: {
        if (scrollTop !== undefined) {
            titleBackgroundOpacity = scrollTop > 100 ? true : false;
        }
    }

    $: {
        if (scrollTopSnapshot !== undefined) {
            const scroller = document.getElementsByClassName("scroller")[0];
            setTimeout(() => {
                scroller.scroll(0, scrollTopSnapshot!);
            }, 0);
        }
    }

    async function animateLoad(id: string) {
        gsap.fromTo("#" + id, {
            opacity: 0,
            scale: 0.8,
        }, {
            opacity: 1,
            scale: 1,
            duration: 1,
            ease: "expo.out",
        });
    }

    async function addNewScreenshots() {
        let newestKnown = 0;
        const allScr = get(rcvScr);

        if (allScr) {
            const newestKey = Object.keys(allScr)[0];
            allScr[newestKey].forEach((scr) => {
                if (scr.CaptureID > newestKnown) newestKnown = scr.CaptureID;
            });
        }

        const newScreenshots = await getScreenshotsNewerThan(newestKnown);
        if (!newScreenshots) return;

        screenshotStore.update((prev) => {
            if (!prev) return prev;

            return [...newScreenshots, ...prev];
        });
    }

    function subscribeToScreenshotEvent() {
        EventsOn("rcv:screenshotran", async (lastId) => {
            addNewScreenshots();
        });
    }

    async function getOlderScreenshots() {
        let oldestKnownId: number | undefined;
        const allScreenshots = get(rcvScr);

        if (allScreenshots && typeof allScreenshots === "object") {
            const screenshotArrays = Object.values(allScreenshots);
            if (screenshotArrays.length > 0) {
                const flatScreenshots = screenshotArrays.flat();
                oldestKnownId = Math.min(
                    ...flatScreenshots.map((scr) => scr.CaptureID)
                );
            }
        }

        if (oldestKnownId === undefined) {
            // console.log("No existing screenshots found");
            return;
        }

        try {
            const newScreenshots = await getScreenshotsOlderThan(
                oldestKnownId,
                30
            );

            if (newScreenshots === null || newScreenshots.length === 0) {
                allScreenshotsLoaded = true;
                return;
            }

            screenshotStore.update((prev) => {
                if (!Array.isArray(prev)) return newScreenshots;
                return [...prev, ...newScreenshots];
            });
        } catch (error) {
            console.error("Error fetching older screenshots:", error);
            allScreenshotsLoaded = true;
        }
    }


    onMount(() => {
        const unsubscribe = scrollStore.subscribe(
            (scrollPos) => (scrollTop = scrollPos)
        );

        loadMoreDivObserver = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        getOlderScreenshots();
                    }
                });
            },
            { threshold: 0.1 }
        ); // Trigger when 10% of the element is visible

        subscribeToScreenshotEvent();

        isInitialLoad = false;

        const scrUnsubscribe = processedScreenshotStore.subscribe((value) => {
            if (!isInitialLoad) {
                rcvScr.set(value);
            }
        });

        return () => {
            unsubscribe();
            scrUnsubscribe();
            clearTimeout(loadMoreDivObserverTimeout);
            EventsOff("rcv:screenshotran");
            loadMoreDivObserver.disconnect();
        }; // Unsubscribe from scrollStore when destroying. Clear timeout for loadMoreDivObserverTimeout, so that the timeout doesn't run in another page
    });

    $: if (loadMoreDiv) {
        loadMoreDivObserverTimeout = setTimeout(() => {
            loadMoreDivObserver.observe(loadMoreDiv);
        }, 100) as unknown as number;;
    }

    function selectAllFromDate(event: any, date: string) {
        let checkAll: boolean = checkedItems[date] ? false : true;
        checkedItems[date] = !checkAll;

        rcvScr.update((prev) => {
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

    beforeNavigate(async (nav) => {
        const targetScr = document.getElementById("s" + nav.to!.params!.id)!;
        if (!targetScr) return;
        targetScr.classList.add("transition-box-container");
        targetScr.children[1].classList.add("transition-box-content");
        targetScr.querySelector(".icon")!.classList.add("exclude-transition");
        await nav.complete;
    });

    export const snapshot: Snapshot<number> = {
        capture: () => {
            return scrollTop;
        },
        restore: (snapshot) => {
            scrollTopSnapshot = snapshot;
        },
    };

    function scrClicked(date: string, captureId: number, event: MouseEvent) {
        // If 'selecting' is true and a screenshot was clicked, it was selected, not routed to. Early return
        if (selecting) {
            rcvScr.update((prev) => {
                if (!prev || !prev[date]) return prev;

                const updatedDate = prev[date].map((s) =>
                    s.CaptureID === captureId
                        ? { ...s, Selected: !s.Selected }
                        : s
                );

                // Check if all screenshots are selected
                const allSelected = updatedDate.every((scr) => scr.Selected);

                // Update the checkedItems state
                checkedItems[date] = allSelected;

                return {
                    ...prev,
                    [date]: updatedDate,
                };
            });
            return;
        }

        // Navigate to the screenshot
        goto(`/screenshots/${captureId}`);
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
        rcvScr.update((prev) => {
            if (!prev) return prev;

            const updatedScreenshots = Object.entries(prev).reduce(
                (acc, [date, screenshots]) => {
                    acc[date] = screenshots.map((scr) => ({
                        ...scr,
                        Selected: false,
                    }));
                    return acc;
                },
                {} as typeof prev
            );

            return updatedScreenshots;
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

    async function generateReport() {
        const allScrs = get(rcvScr);
        if (!allScrs) return;

        const selectedIds: number[] = [];
        Object.keys(allScrs).forEach((key) =>
            allScrs[key].map((scr) => {
                if (scr.Selected === true) {
                    selectedIds.push(scr.CaptureID);
                }
            })
        );

        try {
            const reportId: number | undefined =
                await GenerateReportFromScreenshotIds(selectedIds);

            if (!reportId) throw "No report ID was found";
            addNewDialog({
                title: "Report generated",
                description: `A new report was generated!`,
                primaryButtonCallback: () => goto(`/reports/${reportId}`),
                primaryButtonName: "Open report",
            });
        } catch (err) {
            console.error(err);
            addNewDialog({
                title: "Report generation failed",
                description: `The following error was received: ${err}`,
            });
        }
    }

    /**
     * Ask the user if they're sure they want to delete N number of selected screenshots. If they proceed, call deleteSelectedStep2
     */
    async function deleteSelectedStep1() {
        const allScrs = get(rcvScr);
        if (!allScrs) return;

        const selectedIds: number[] = [];
        Object.keys(allScrs).forEach((key) =>
            allScrs[key].map((scr) => {
                if (scr.Selected === true) {
                    selectedIds.push(scr.CaptureID);
                }
            })
        );

        try {
            addNewDialog({
                title: "Confirmation",
                description: `${selectedIds.length} screenshot${selectedIds.length !== 1 ? "s" : ""} will be deleted. Are you sure you want to proceed?`,
                primaryButtonName: "Delete",
                primaryButtonCallback: async () =>
                    await deleteSelectedStep2(selectedIds),
                secondaryButtonName: "Cancel",
                secondaryButtonCallback: () => console.log("Cancelled"),
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
            await DeleteScreenshotsById(ids);
            screenshotStore.update((prev) => {
                const filtered = prev.filter((r) => !ids.includes(r.CaptureID));
                return filtered;
            });
        } catch (err: any) {
            addNewDialog({
                title: "Error",
                description: `Could not delete screenshots with the following error: ${err}`,
            });
        }
    }
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
                class="page-title text-2xl -tracking-wide opacity-85 w-1/2 z-40 pointer-events-none"
            >
                Screenshots
            </h1>
            <div
                class="flex gap-2 justify-self-end -tracking-wide text-xl z-30 text-black font-semibold"
            >
                {#if selecting}
                    <div
                        on:click={generateReport}
                        class="transition-all -tracking-wide text-xl px-4 p-2 bg-opacity-80 hover:bg-blue-300 active:scale-[99%] cursor-pointer bg-blue-400 text-black font-semibold rounded-lg"
                    >
                        Report
                    </div>
                    <div
                        on:click={deleteSelectedStep1}
                        class="transition-all -tracking-wide text-xl px-4 p-2 bg-opacity-80 cursor-pointer active:scale-[99%] hover:bg-red-300 bg-red-400 text-black font-semibold rounded-lg"
                    >
                        Delete
                    </div>
                {/if}

                <div
                    on:click={multiSelectClicked}
                    class="-tracking-wide transition-all text-xl px-4 p-2 bg-opacity-80 cursor-pointer active:scale-[99%] hover:bg-opacity-90 bg-neutral-200 dark:bg-white text-black font-semibold rounded-lg"
                >
                    {selecting ? "Cancel" : "Select"}
                </div>
            </div>
        </div>

        {#if $rcvScr !== undefined}
            {#each Object.entries($rcvScr) as [date, screenshots]}
                <div>
                    <div
                        class="flex gap-5 sticky items-center top-16 z-30 pointer-events-none"
                    >
                        <h2 class="text-3xl font-bold tracking-wider">
                            {date}
                        </h2>
                        <div class="checkbox opacity-0 pointer-events-none">
                            <Checkbox
                                id={date}
                                bind:checked={checkedItems[date]}
                                on:checked={(e) => selectAllFromDate(e, date)}
                            ></Checkbox>
                        </div>
                    </div>

                    <div class="my-4 grid grid-cols-2 gap-4">
                        {#each screenshots as s (s.CaptureID)}
                            <div
                                on:click={(event) =>
                                    scrClicked(date, s.CaptureID, event)}
                                data-intersect
                                on:intersect={(e) => {s.Visible = e.detail.isIntersecting}}
                                class="m-0 p-0 {s.Visible ? '' : 'invisible'} aspect-video"
                            >
                                <div
                                    id="s{s.CaptureID}"
                                    class="group cursor-pointer relative rounded-lg bg-neutral-200 dark:bg-neutral-800 outline overflow-hidden outline-1 outline-neutral-300 dark:outline-neutral-900 p-1 shadow-2xl"
                                >
                                    {#if selecting}
                                        <div
                                            class="flex items-center justify-center absolute z-30 transition-all opacity-50 left-0 top-0 right-0 bottom-0 {s.Selected
                                                ? 'bg-neutral-200 dark:bg-neutral-500'
                                                : 'bg-neutral-400 dark:bg-neutral-800'}"
                                        >
                                            <div
                                                class="transition-all opacity-0 scale-90 w-[50%] {s.Selected
                                                    ? 'opacity-90 scale-100'
                                                    : ''}"
                                            >
                                                <Checkmark strokeColor="#fff"
                                                ></Checkmark>
                                            </div>
                                        </div>
                                    {/if}

                                    <div
                                        class="icon absolute top-2 right-2 w-8 h-8 z-40 p-1 rounded-full bg-blue-200"
                                    >
                                        {#if s.Description}
                                            {#if s.ReportID}
                                                <DoneAllIcon
                                                    title="Screenshot was included in a report"
                                                    strokeColor="#2966eb"
                                                ></DoneAllIcon>
                                            {:else}
                                                <DoneIcon
                                                    title="Description was generated for screenshot"
                                                    strokeColor="#2966eb"
                                                ></DoneIcon>
                                            {/if}
                                        {:else}
                                            <XIcon
                                                title="No description generated yet"
                                                strokeColor="#3c424a"
                                            ></XIcon>
                                        {/if}
                                    </div>

                                    <img
                                        alt="screenshot"
                                        on:load|once={() =>
                                            animateLoad("s" + s.CaptureID)}
                                        class="group-hover:scale-[99%] group-active:scale-[95%] transition-all flex rounded-md object-contain select-none pointer-events-none"
                                        loading="lazy"
                                        src={s.Visible ? s.Screenshot : ''}
                                    />

                                    <h3 class="flex-shrink-0 pt-1 pl-2 text-black dark:text-white">
                                        Snapped at {s.Time}
                                    </h3>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/each}
            {#if !allScreenshotsLoaded}
                <div class="h-10 w-full" bind:this={loadMoreDiv}></div>
            {:else}
                <div></div>
            {/if}
        {:else}
            No screenshots available yet. Turn on the screenshot schedule to get started!
        {/if}
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
        pointer-events: none;
        background: linear-gradient(
            180deg,
            rgb(0, 0, 0) 0%,
            rgba(0, 0, 0, 0) 100%
        );
    }

    @media (prefers-color-scheme: light) {
        .top-gradient-bg::after {
            background: linear-gradient(
                180deg,
                rgb(255, 255, 255) 0%,
                rgba(255, 255, 255, 0) 100%
            )
        }
    }
</style>
