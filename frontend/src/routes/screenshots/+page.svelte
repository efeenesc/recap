<script lang="ts">
    import gsap from "gsap";
    import { onDestroy, onMount } from "svelte";
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
    import { GenerateReportFromScreenshotIds, DeleteScreenshotsById } from "$lib/wailsjs/go/app/AppMethods.js";
    import { addNewDialog } from "../../utils/dialog.ts";
    import { screenshotStore } from "$lib/stores/ScreenshotStore.ts";
    import { scrollStore } from "$lib/stores/ScrollStore.ts";
    import DoneAllIcon from "../../icons/DoneAllIcon.svelte";
    import DoneIcon from "../../icons/DoneIcon.svelte";
    import XIcon from './../../icons/XIcon.svelte';

    interface Data {
        streamed: {
            items: {
                subscribe: (
                    run: (value: DatedScreenshot) => void
                ) => () => void;
            };
        };
    }

    export let data: Data;

    let loadMoreDiv: Element;
    let loadMoreDivObserver: IntersectionObserver;
    let prevId: number | undefined;
    let selecting: boolean = false;
    const rcvScr = writable<DatedScreenshot | undefined>();
    let checkedItems: { [key: string]: boolean | undefined } = {};

    let scrollTop: number = 0;
    let titleBackgroundOpacity: boolean = false;

    let allScreenshotsLoaded: boolean = false;

    $: {
        if (scrollTop) {
            titleBackgroundOpacity = scrollTop > 100 ? true : false;
        }
    }

    async function getData(): Promise<DatedScreenshot> {
        return new Promise((resolve) => {
            data.streamed.items.subscribe((items) => {
                rcvScr.set(items);
                resolve(items);
            });
        });
    }

    async function animateLoad(id: string) {
        if (prevId && "s" + prevId === id) {
            const target = document.getElementById("s" + prevId)!;
            target.style.scale = "1";
            target.style.opacity = "1";
            return;
        }

        gsap.to("#" + id, {
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
            console.log("No existing screenshots found");
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
        }
    }

    onDestroy(() => {
        EventsOff("rcv:screenshotran");
        loadMoreDivObserver.disconnect();
    });

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

        return () => {
            unsubscribe();
        }; // Unsubscribe from scrollStore when destroying
    });

    $: if (loadMoreDiv) {
        loadMoreDivObserver.observe(loadMoreDiv);
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
        targetScr.querySelector('.icon')!.classList.add("exclude-transition");
        await nav.complete;
    });

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
                description: `${selectedIds.length} screenshot${selectedIds.length !== 1 ? 's' : ''} will be deleted. Are you sure you want to proceed?`,
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
            await DeleteScreenshotsById(ids);
            screenshotStore.update(prev => {
                const filtered = prev.filter(r => !ids.includes(r.CaptureID));
                return filtered;
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
    <div class="pb-2 gap-5 flex flex-col">
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
                        class="opacity-0 transition-all {selecting
                            ? 'opacity-100'
                            : ''} -tracking-wide text-xl px-4 p-2 bg-opacity-80 hover:bg-blue-300 active:scale-[99%] cursor-pointer bg-blue-400 text-black font-semibold rounded-lg"
                    >
                        Report
                    </div>
                    <div
                        on:click={deleteSelectedStep1}
                        class="opacity-0 transition-all {selecting
                            ? 'opacity-100'
                            : ''} -tracking-wide text-xl px-4 p-2 bg-opacity-80 cursor-pointer active:scale-[99%] hover:bg-red-300 bg-red-400 text-black font-semibold rounded-lg"
                    >
                        Delete
                    </div>
                {/if}

                <div
                    on:click={multiSelectClicked}
                    class="-tracking-wide transition-all text-xl px-4 p-2 bg-opacity-80 cursor-pointer active:scale-[99%] hover:bg-opacity-90 bg-white text-black font-semibold rounded-lg"
                >
                    {selecting ? "Cancel" : "Select"}
                </div>
            </div>
        </div>

        {#await getData()}
            Loading screenshots...
        {:then}
            {#if $rcvScr}
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
                                    on:checked={(e) =>
                                        selectAllFromDate(e, date)}
                                ></Checkbox>
                            </div>
                        </div>

                        <div class="my-4 grid grid-cols-2 gap-4">
                            {#each screenshots as s (s.CaptureID)}
                                <div
                                    on:click={(event) =>
                                        scrClicked(date, s.CaptureID, event)}
                                    class="m-0 p-0"
                                >
                                    <div
                                        id="s{s.CaptureID}"
                                        class="group cursor-pointer relative rounded-lg bg-neutral-800 outline overflow-hidden outline-1 outline-neutral-900 p-1 mr-5 shadow-2xl opacity-0 scale-95"
                                    >
                                        {#if selecting}
                                            <div
                                                class="flex items-center justify-center absolute z-30 transition-all bg-opacity-50 left-0 top-0 right-0 bottom-0 {s.Selected
                                                    ? 'bg-neutral-500'
                                                    : 'bg-neutral-800'}"
                                            >
                                                <div
                                                    class="transition-all opacity-0 scale-90 w-[50%] {s.Selected
                                                        ? 'opacity-90 scale-100'
                                                        : ''}"
                                                >
                                                    <Checkmark
                                                        strokeColor="#fff"
                                                    ></Checkmark>
                                                </div>
                                            </div>
                                        {/if}

                                        <div class="icon absolute top-2 right-2 w-8 h-8 z-40 p-1 rounded-full bg-blue-200">
                                            {#if s.Description}
                                                {#if s.ReportID}
                                                    <DoneAllIcon title="Screenshot was included in a report" strokeColor="#2966eb"></DoneAllIcon>
                                                {:else}
                                                    <DoneIcon title="Description was generated for screenshot" strokeColor="#2966eb"></DoneIcon>
                                                {/if}
                                            {:else}
                                            <XIcon title="No description generated yet" strokeColor="#3c424a"></XIcon>
                                            {/if}
                                        </div>

                                        <img
                                            alt="screenshot"
                                            on:load|once={() =>
                                                animateLoad("s" + s.CaptureID)}
                                            class="group-hover:scale-[99%] group-active:scale-[95%] transition-all flex rounded-md object-contain select-none pointer-events-none"
                                            loading="lazy"
                                            src={s.Screenshot}
                                        />

                                        <h3 class="flex-shrink-0 pt-1 pl-2">
                                            Snapped at {s.Time}
                                        </h3>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/each}
                {#if !allScreenshotsLoaded}
                    <div bind:this={loadMoreDiv}>Loading more screenshots...</div>
                {:else}
                    <div></div>
                {/if}
            {:else}
                No screenshots available.
            {/if}
        {:catch error}
            Error loading screenshots: {error.message}
        {/await}
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
</style>
