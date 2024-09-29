<script lang="ts">
    import gsap from "gsap";
    import { onDestroy, onMount } from "svelte";
    import type { DatedScreenshot } from "../../types/ExtendedScreenshot.interface.ts";
    import { goto, afterNavigate } from "$app/navigation";
    import Checkbox from "../../components/checkbox/Checkbox.svelte";
    import { get, writable } from "svelte/store";
    import Checkmark from "../../icons/Checkmark.svelte";
    import { EventsOff, EventsOn } from "$lib/wailsjs/runtime/runtime.js";
    import { addScreenshotToStore } from "../../utils/screenshot.ts";
    import { dialogStore } from "$lib/stores/DialogStore.ts";
    import { GenerateReportFromScreenshotIds } from "$lib/wailsjs/go/app/AppMethods.js";
    import { addNewDialog } from "../../utils/dialog.ts";

    interface Data {
        streamed: {
            items: Promise<DatedScreenshot | undefined>;
        };
    }

    let title: HTMLDivElement;
    let prevId: number | undefined;
    let selecting: boolean = false;
    const rcvScr = writable<DatedScreenshot | undefined>();
    let checkedItems: { [key: string]: boolean | undefined } = {};
    export let data: Data;

    async function getData(): Promise<DatedScreenshot | undefined> {
        try {
            const items = await data.streamed.items;
            rcvScr.set(items);
            return items;
        } catch (err) {
            console.error(err);
            throw err;
        }
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

    function subscribeToScreenshotEvent() {
        EventsOn("rcv:screenshotran", async (lastId) => {
            await addScreenshotToStore(lastId);
            addNewDialog({
                title: "New item added",
                description: `New screenshot ID: ${lastId}`,
            });
        });
    }

    onDestroy(() => {
        EventsOff("rcv:screenshotran");
    });

    onMount(() => {
        subscribeToScreenshotEvent();

        const observer = new IntersectionObserver(
            (asd) => {
                console.log(asd);
            },
            {
                rootMargin: "-1px 0px 0px 0px",
                threshold: [1],
            }
        );

        observer.observe(title);
    });

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

    afterNavigate((nav) => {
        if (!nav.from) return;
        const targetScr = document.getElementById("s" + nav.from.params)!;
        if (!targetScr) return;
        targetScr.classList.add("transition-box-container");
        targetScr.children[1].classList.add("transition-box-content");
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

        console.log(selectedIds);

        try {
            const reportId: number | undefined =
                await GenerateReportFromScreenshotIds(selectedIds);

            console.log(reportId);
            if (!reportId) throw "No report ID was found";
            addNewDialog({
                title: "Report generated",
                description: `A new report was generated!`,
                buttonLink: `/reports/${reportId}`,
                buttonLinkDescription: "Open report",
            });
        } catch (err) {
            console.error(err);
            addNewDialog({
                title: "Report generation failed",
                description: `The following error was received: ${err}`,
            });
        }
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

                                        <img
                                            alt="screenshot"
                                            on:load|once={() =>
                                                animateLoad("s" + s.CaptureID)}
                                            class="group-hover:scale-[99%] group-active:scale-[95%] transition-all flex rounded-md object-contain select-none pointer-events-none"
                                            loading="lazy"
                                            src={s.Screenshot}
                                        />
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/each}
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
</style>
