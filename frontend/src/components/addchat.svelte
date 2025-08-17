<script lang="ts">
  import { onMount, mount } from "svelte";
  import Conversation from "./conversation.svelte";
  let wrapper: HTMLElement;
  type AddChatProps = {
    onClose: () => void;
  };

  let { onClose }: AddChatProps = $props();

  function showConversation() {
    const el = document.createElement("div");
    mount(Conversation, {
      props: {
        coordinates: {
          x: wrapper.getBoundingClientRect().x + 40,
          y: wrapper.getBoundingClientRect().y - 30,
        },
        create: true,
        onClose: () => {
          onClose();
          document.body.removeChild(el);
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
