<script lang="ts">
  import { ArrowUp, Ellipsis, Image, X } from "@lucide/svelte";
  import Avatar from "./avatar.svelte";
  import { createThread, type Thread } from "$lib/services/threads.svelte";
  type ConversationProps = {
    coordinates: {
      x: number;
      y: number;
    };
    lat: number;
    long: number;
    onClose: () => void;
    onCreate: (thread: Thread) => void;
    create: boolean;
  };
  let { coordinates, onClose, create, onCreate, lat, long }: ConversationProps =
    $props();

  let inputValue = $state("");
  let isSendButtonDisabled = $derived(inputValue.trim() == "");

  async function onCreateThread() {
    const thread = await createThread(lat, long, inputValue);
    onCreate(thread);
  }

  function onSend() {
    if (create) {
      onCreateThread();
    } else {
    }
  }
</script>

<div
  class="w-[100vw] height-[100dvh] fixed top-0 left-0 right-0 bottom-0 overflow-hidden z-[99]"
  onclick={onClose}
>
  <div
    class="container shadow-md fixed rounded-md w-[336px]"
    style={`top:${coordinates.y}px;left:${coordinates.x}px`}
    onclick={(e) => {
      e.stopPropagation();
    }}
  >
    <div
      class="p-3 py-2 flex items-center justify-between w-full border-b border-[#445160]"
    >
      <div class="flex items-center">
        <p>{create ? "Start Thread" : "Thread"}</p>
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

    <div class="bottom p-3 flex gap-2 py-4">
      <div>
        <Avatar iconSize={10} />
      </div>
      <div class="flex flex-col w-full">
        <textarea
          class="textarea w-full text-input"
          placeholder="Write"
          bind:value={inputValue}
        ></textarea>
        <div class="flex items-center justify-between p-2 bottom-buttons">
          <div class="flex items-center gap2">
            <button class="icon-btn">
              <Image size={16} />
            </button>
          </div>
          <button
            class="btn btn-circle btn-primary w-[24px] h-[24px] p-1"
            disabled={isSendButtonDisabled}
            onclick={onSend}
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
    cursor: pointer;
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
