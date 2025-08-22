<script lang="ts">
  import { onMount, mount, onDestroy } from "svelte";
  import Conversation from "./conversation.svelte";
  import type { Thread } from "$lib/services/threads.svelte";
  let wrapper: HTMLElement;
  type AddChatProps = {
    onClose: () => void;
    onCreate: (thread: Thread) => void;
    lat: number;
    long: number;
  };

  let { onClose, onCreate, lat, long }: AddChatProps = $props();

  let el: Element | null = $state(null);

  function showConversation() {
    el = document.createElement("div");
    mount(Conversation, {
      props: {
        lat,
        long,
        onCreate: (thread) => {
          onCreate(thread);
          if (el) document.body.removeChild(el);
        },
        coordinates: {
          x: wrapper.getBoundingClientRect().x + 40,
          y: wrapper.getBoundingClientRect().y - 30,
        },
        create: true,
        onClose: () => {
          onClose();
          if (el) document.body.removeChild(el);
        },
      },
      target: el,
    });
    document.body.appendChild(el);
  }

  onMount(() => {
    showConversation();
  });
</script>

<div class="add-chat-wrapper" bind:this={wrapper}>
  <div class="add-chat">
    <div class="pointer shadow-md"></div>
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
</style>
