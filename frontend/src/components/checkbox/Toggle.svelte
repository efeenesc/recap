<script lang="ts">
    import { createEventDispatcher } from "svelte";
    export let checked: boolean = false;
    export let id: string;
    let _class: string;
    export { _class as class };

    const dispatch = createEventDispatcher();

    function handleChange(event: any) {
        console.log(event);
        dispatch("checked", { checked: event.target.checked });
    }
</script>

<!-- Toggle style by Flowbite. License below. -->
<!-- 
Copyright (c) Themesberg (Bergside Inc.)
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
-->
<label class="flex items-center cursor-pointer {_class}">
    <input
        type="checkbox"
        class="sr-only peer"
        {id}
        bind:checked
        on:change={handleChange}
    />
    <div
        class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 rounded-full peer dark:bg-neutral-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-500"
    ></div>
</label>

<style global lang="postcss">
    @tailwind utilities;
    @tailwind components;
    @tailwind base;
</style>
