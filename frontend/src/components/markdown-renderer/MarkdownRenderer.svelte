<script lang="ts">
  import type { MdNode } from '$lib/markdown/Markdown.interface.ts';
  import { BrowserOpenURL } from '$lib/wailsjs/runtime/runtime.js'

  export let parsedContent: MdNode[] | string;
  let parsedNode: MdNode[];

  function linkClicked(link: string | undefined) {
    if (link)
      BrowserOpenURL(link);
  }

  $: {
    if (typeof parsedContent === 'string') {
      parsedNode = [{ type: 'text', content: parsedContent }];
    } else {
      parsedNode = parsedContent;
    }
  }
</script>

{#if parsedNode}
  {#each parsedNode as node (node)}
    {#if node.type === 'text'}
      <span>{node.content}</span>
    {:else if node.type === 'p'}
      <p>
        <svelte:self parsedContent={node.content} />
      </p>
    {:else if node.type === 'ul'}
      <ul>
        <svelte:self parsedContent={node.content} />
      </ul>
    {:else if node.type === 'ol'}
      <ol>
        <svelte:self parsedContent={node.content} />
      </ol>
    {:else if node.type === 'li'}
      <li>
        <svelte:self parsedContent={node.content} />
      </li>
    {:else if node.type === 'h1'}
      <h1>
        <svelte:self parsedContent={node.content} />
      </h1>
    {:else if node.type === 'h2'}
      <h2>
        <svelte:self parsedContent={node.content} />
      </h2>
    {:else if node.type === 'h3'}
      <h3>
        <svelte:self parsedContent={node.content} />
      </h3>
    {:else if node.type === 'h4'}
      <h4>
        <svelte:self parsedContent={node.content} />
      </h4>
    {:else if node.type === 'h5'}
      <h5>
        <svelte:self parsedContent={node.content} />
      </h5>
    {:else if node.type === 'h6'}
      <h6>
        <svelte:self parsedContent={node.content} />
      </h6>
    {:else if node.type === 'code'}
      <code><svelte:self parsedContent={node.content} /></code>
    {:else if node.type === 'bq'}
      <blockquote>
        <svelte:self parsedContent={node.content} />
      </blockquote>
    {:else if node.type === 'b'}
      <strong>
        <svelte:self parsedContent={node.content} />
      </strong>
    {:else if node.type === 'i'}
      <em>
        <svelte:self parsedContent={node.content} />
      </em>
    {:else if node.type === 'bi'}
      <strong>
        <em>
          <svelte:self parsedContent={node.content} />
        </em>
      </strong>
    {:else if node.type === 'st'}
      <del>
        <svelte:self parsedContent={node.content} />
      </del>
    {:else if node.type === 'br'}
      <br />
    
    {:else if node.type === 'a'}
      <!-- svelte-ignore a11y-missing-attribute -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <a on:click={() => linkClicked(node.url)} target="_blank" class="flex">
        <svelte:self parsedContent={node.content}></svelte:self>
      </a>
    {/if}
  {/each}
{/if}

<slot></slot>

<style global lang="postcss">
  * {
    font-size: 1em
  }
  
  strong {
    font-weight: 400;
  }

  h2 {
    margin-bottom: 0.5em;
    font-size: 2em;
  }

  ul {
    margin-top: 12px;
    margin-bottom: 4px;
    display: block;
    list-style-type: disc;
    margin-top: 1em;
    margin-bottom: 1 em;
    margin-left: 0;
    margin-right: 0;
    padding-left: 40px;
  }

  code {
    margin-left: 6px;
    margin-right: 0px;
  }

  a {
    text-decoration: underline;
    text-decoration-color: white;
    cursor: pointer;
    margin: 0px;
    padding: 0px;
  }

  br {
    content: "";
    display: block;
    height: 10px;
    width: 100px;
  }
</style>
