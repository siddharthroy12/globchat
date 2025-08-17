<script lang="ts">
  import { ArrowUp, Ellipsis, X } from "@lucide/svelte";
  import Avatar from "./avatar.svelte";
  type ConversationProps = {
    coordinates: {
      x: number;
      y: number;
    };
    onClose: () => void;
    create: boolean;
  };
  let { coordinates, onClose, create }: ConversationProps = $props();

  let inputValue = $state("");
  let isSendButtonDisabled = $derived(inputValue.trim() == "");
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="w-[100vw] height-[100dvh] fixed top-0 left-0 right-0 bottom-0 overflow-hidden z-[99]"
  onclick={onClose}
>
  <!-- svelte-ignore a11y_no_noninteractive_tabindex -->
  <div
    class="container shadow-md fixed rounded-md w-[336px]"
    style={`top:${coordinates.y}px;left:${coordinates.x}px`}
    onclick={(e) => {
      e.stopPropagation();
    }}
  >
    <div
      class="p-3 flex items-center justify-between w-full border-b border-[#445160]"
    >
      <div class="flex items-center">
        <p>{create ? "Start Thread" : "Tread"}</p>
      </div>
      <div class="flex items-center gap-2">
        {#if !create}
          <button class="icon-btn">
            <Ellipsis />
          </button>
        {/if}

        <button class="icon-btn" onclick={onClose}>
          <X />
        </button>
      </div>
    </div>
    {#if !create}
      <div class="body p-3"></div>
    {/if}

    <div class="bottom p-3 flex gap-2">
      <div>
        <Avatar />
      </div>
      <div class="flex flex-col w-full">
        <textarea
          class="textarea w-full text-input"
          placeholder="Write"
          bind:value={inputValue}
        ></textarea>
        <div class="flex items-center justify-end p-2 bottom-buttons">
          <button
            class="btn btn-circle btn-primary w-[24px] h-[24px] p-1"
            disabled={isSendButtonDisabled}
          >
            <ArrowUp />
          </button>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .container {
    background: var(--color-base-300);
  }

  .icon-btn {
    width: 27px;
    height: 27px;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0;
  }

  .icon-btn:hover {
    background: var(--color-base-100);
  }

  .text-input {
    border-bottom-left-radius: 0px;
    border-bottom-right-radius: 0px;
    resize: none;
  }

  .bottom-buttons {
    border: 1px solid color-mix(in oklab, var(--color-base-content) 20%, #0000);
    background-color: var(--color-base-100);
    border-top: 0px;

    border-radius: var(--radius-field);
    border-top-left-radius: 0px;
    border-top-right-radius: 0px;
  }
</style>
