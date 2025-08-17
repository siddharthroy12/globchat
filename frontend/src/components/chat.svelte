<script lang="ts">
  import Avatar from "./avatar.svelte";
  import { mount } from "svelte";
  import Conversation from "./conversation.svelte";

  let conversationOpen = $state(false);
  let wrapper: HTMLElement;

  function showConversation() {
    conversationOpen = true;
    const el = document.createElement("div");
    mount(Conversation, {
      props: {
        coordinates: {
          x: wrapper.getBoundingClientRect().x + 40,
          y: wrapper.getBoundingClientRect().y - 30,
        },
        create: true,
        onClose: () => {
          console.log("called");
          conversationOpen = false;
          document.body.removeChild(el);
        },
      },
      target: el,
    });

    document.body.appendChild(el);
  }
</script>

<div class="conversation-bubble-wrapper" bind:this={wrapper}>
  <button
    class="conversation-bubble"
    class:active={conversationOpen}
    onclick={showConversation}
  >
    <div class="avatar-wrapper">
      <Avatar />
    </div>
    <div class="content">
      <p class="username">Siddharth Roy <span class="time">1 hr. ago</span></p>
      <p class="last-chat"></p>
      <p>I live here</p>
      <p class="time">1 reply</p>
    </div>
  </button>
</div>

<style>
  .conversation-bubble-wrapper {
    position: relative;
  }
  .conversation-bubble {
    cursor: pointer;
    position: absolute;
    bottom: 0;
    left: 0;
    background: var(--color-base-300);
    padding: 5px;
    color: var(--color-primary-content);
    width: 30px;
    max-height: 30px;
    border-radius: 50px;
    border-bottom-left-radius: 0;
    display: flex;
    justify-content: start;
    text-align: start;
    overflow: hidden;
    align-items: start;
    gap: 10px;
    transition:
      width 200ms,
      max-height 200ms,
      padding 200ms;
  }
  .active {
    box-shadow: 0 0 0 2px var(--color-primary);
  }
  .avatar-wrapper {
    flex-shrink: 0;
  }
  .conversation-bubble:not(.active):hover {
    z-index: 999;
    width: 200px;
    max-height: 200px;
    padding: 10px;
    border-radius: 20px;
    border-bottom-left-radius: 0;
  }
  .content {
    display: none;
  }
  .conversation-bubble:not(.active):hover .content {
    display: block;
  }
  .username {
    font-weight: bold;
    text-wrap: nowrap;
  }
  .last-chat {
    width: 180px;
  }
  .time {
    font-weight: 400;
    color: var(--color-base-content);
    opacity: 0.8;
  }
</style>
