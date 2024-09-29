<script lang="ts">
    import { createEventDispatcher } from "svelte";
    export let checked: boolean = false;
    export let id: string;
    let _class: string;
    export {_class as class}
    export let neu: boolean = false;

    const dispatch = createEventDispatcher();

    function handleChange(event: any) {
        dispatch("checked", { checked: event.target.checked });
    }
</script>

<label class="inline-flex items-center cursor-pointer {_class}">
  <input type="checkbox" value="" class="sr-only peer" {id} bind:checked={checked} on:change={handleChange}>
  <div class="{neu ? 'neu' : '' } relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 rounded-full peer dark:bg-neutral-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-green-600"></div>
</label>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;

    .neu {
        box-shadow:  5px 5px 10px #bebebe,
                    -5px -5px 10px #ffffff;
    }
</style>
