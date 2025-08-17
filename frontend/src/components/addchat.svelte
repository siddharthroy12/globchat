<script lang="ts">
  import { onMount, mount } from "svelte";

  let inputValue = $state("");
  let inputEl: HTMLElement;

  onMount(() => {
    inputEl.focus();
  });
</script>

<div class="add-chat-wrapper">
  <div class="add-chat">
    <div class="pointer shadow-md"></div>
    <div class="input-container shadow-md">
      <div class="text-area-wrapper">
        <textarea
          bind:this={inputEl}
          bind:value={inputValue}
          onclick={(e) => {
            e.stopPropagation();
            //@ts-ignore
            e.target.focus();
          }}
        ></textarea>
      </div>
      <div class="buttons-wrapper"></div>
    </div>
  </div>
</div>

<style>
  .add-chat-wrapper {
    position: relative;
  }
  .add-chat {
    position: absolute;
    bottom: 0;
    left: 0;
    display: flex;
    gap: 10px;
    align-items: center;
  }
  .pointer {
    border-radius: 50px;
    border-bottom-left-radius: 0;
    width: 30px;
    height: 30px;
    background-color: var(--color-primary);
  }
  .input-container {
    background-color: var(--color-base-100);
    border-radius: 10px;
    color: var(--color-primary-content);
  }
  .text-area-wrapper {
    padding: 10px;
  }
</style>
