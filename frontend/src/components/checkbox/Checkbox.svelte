<script lang="ts">
	import DoneIcon from './../../icons/DoneIcon.svelte';
    import { createEventDispatcher } from "svelte";
    export let checked: boolean = false;
    export let id: string;
    export const setState = (newState: boolean) => {
        checked = newState;
    };

    const dispatch = createEventDispatcher();

    function handleChange(event: any) {
        dispatch("checked", { checked: event.target.checked });
    }
</script>

<div>
    <input type="checkbox" {id} bind:checked on:change={handleChange} 
    class="peer hidden" />
    <label
        class="bg-neutral-300 dark:bg-white border block border-neutral-50 rounded-full relative transition-all cursor-pointer w-8 h-8
        peer-checked:bg-green-600 peer-checked:border-green-600 peer-checked:after:opacity-100 peer-active:scale-95"
        for={id}
    >
        <DoneIcon title="" strokeColor="#fff"></DoneIcon>
    </label>
</div>

<style>
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
