<script lang="ts">
    import { afterUpdate, onMount } from "svelte";
    import gsap from "gsap";

    export let bodyInner: number;
    export let bodyHeight: number;
    export let bodyScroll: number;
    export let _class: string | undefined;
    export let style: string | undefined = undefined;

    export { _class as class };

    let thisHeight: number;
    let scrollbarHeight: number;
    let scrollbarPosition: number;
    let virtualScrollHeight: number;
    let hideTimeout: number | null;
    let hidden: boolean = true;

    function calculateScrollHeight(): number {
        const ratio = (thisHeight / bodyHeight) * 100;
        return ratio;
    }

    function calculateScrollbarPosition(): number {
        const scrollPos =
            ((bodyScroll + bodyInner) / bodyHeight) * 100 - virtualScrollHeight;

        return scrollPos;
    }

    function startHideTimeout() {
        if (hideTimeout) {
            clearTimeout(hideTimeout);
        }
        hideTimeout = setTimeout(() => {
            gsap.to(".scrollbar", {
                opacity: 0,
                duration: 0.2,
                onComplete: () => {
                    hidden = true;
                },
            });
        }, 500);
    }

    function showScrollbar(hide: boolean = true) {
        if (hidden) {
            gsap.to(".scrollbar", {
                opacity: 1,
                duration: 0.2,
                onComplete: () => {
                    hidden = false;
                },
            });
        }

        if (hide) startHideTimeout();
    }

    onMount(() => {
        scrollbarHeight = (thisHeight / bodyHeight) * thisHeight;
    });

    afterUpdate(() => {
        showScrollbar();
        scrollbarPosition = calculateScrollbarPosition();
        virtualScrollHeight = calculateScrollHeight();
    });

    $: if (bodyScroll !== undefined) {
        gsap.to(".scrollbar", {
            top: scrollbarPosition + "%",
            duration: 0.1,
            ease: "power2.out",
        });
    }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    bind:clientHeight={thisHeight}
    on:mouseenter={() => showScrollbar(false)}
    on:mouseleave={startHideTimeout}
    class="{_class} block"
    style={style ? style : ''}
>
    <div
        bind:clientHeight={scrollbarHeight}
        class="scrollbar relative w-full rounded-full bg-neutral-800 dark:bg-neutral-200"
        style="height: {virtualScrollHeight}%;"
    ></div>
</div>

<style lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
