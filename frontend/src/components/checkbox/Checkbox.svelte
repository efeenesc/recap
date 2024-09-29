<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  export let checked: boolean = false;
  export let id: string;
  export const setState = (newState: boolean) => {
    checked = newState;
  }

  const dispatch = createEventDispatcher();

  function handleChange(event: any) {
    dispatch("checked", { checked: event.target.checked });
  }
</script>

<div class="checkbox-wrapper-18">
  <div class="round">
    <input type="checkbox" {id} bind:checked on:change={handleChange} />
    <label for={id}></label>
  </div>
</div>

<style>
  .checkbox-wrapper-18 .round {
    position: relative;
  }

  .checkbox-wrapper-18 .round label {
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 50%;
    cursor: pointer;
    height: 28px;
    width: 28px;
    display: block;
  }

  .checkbox-wrapper-18 .round label:after {
    border: 2px solid #fff;
    border-top: none;
    border-right: none;
    content: "";
    height: 6px;
    left: 8px;
    opacity: 0;
    position: absolute;
    top: 9px;
    transform: rotate(-45deg);
    width: 12px;
  }

  .checkbox-wrapper-18 .round input[type="checkbox"] {
    transition: all 0.1s;
    visibility: hidden;
    display: none;
    opacity: 0;
  }

  .checkbox-wrapper-18 .round input[type="checkbox"]:checked + label {
    transition: all 0.1s;
    background-color: #66bb6a;
    border-color: #66bb6a;
  }

  .checkbox-wrapper-18 .round input[type="checkbox"]:checked + label:after {
    transition: all 0.1s;
    opacity: 1;
  }
</style>