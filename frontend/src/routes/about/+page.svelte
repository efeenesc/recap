<script lang="ts">
    import MarkdownRenderer from '../../components/markdown-renderer/MarkdownRenderer.svelte';
import type { ExtendedPageContent } from '../../types/BasicPageContent.interface.ts'
    import gsap from "gsap";

    interface Data {
        data: ExtendedPageContent[]
    }

    export let data;

    async function animateLoad(id: string) {
        gsap.to("#" + id, {
            opacity: 1,
            scale: 1,
            duration: 1,
            ease: "expo.out",
        });
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="w-full h-full bypass-padding">
    <div class="pb-2 flex items-end">
        <h1
            class="text-2xl w-fit -tracking-wide opacity-85 text-[#1f1f1f] dark:text-white"
        >
            About
        </h1>
    </div>

    {#if data}
        {#each data.data as s}
            <section class="relative">
                {#if s.Title}
                    <h2 class="text-3xl font-bold tracking-wider mb-2">{s.Title}</h2>
                {/if}

                <div class="my-4">
                    <MarkdownRenderer parsedContent={s.ParsedMarkdown}></MarkdownRenderer>
                </div>
            </section>
        {/each}
    {/if}
</div>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
