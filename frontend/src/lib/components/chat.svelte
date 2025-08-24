<script lang="ts">
  import Avatar from "./avatar.svelte";
  import { mount, onDestroy, onMount } from "svelte";
  import Conversation from "./conversation.svelte";
  import type { Thread } from "$lib/services/threads.svelte";
  import { getTimeAgo } from "$lib/helpers";

  let conversationOpen = $state(false);
  let showContentAnimation = $state(false);
  let wrapper: HTMLElement;

  const {
    lat,
    long,
    user_image,
    username,
    created_at,
    message,
    replies,
    showAnimation,
    id,
    user_id,
    expires_at,
    onDelete,
  }: Thread & { showAnimation: boolean; onDelete: () => void } = $props();

  // Calculate days left until expiration
  function getDaysLeft(): number {
    if (!expires_at) return 0;

    const now = new Date();
    const expirationDate = new Date(expires_at);
    const timeDiff = expirationDate.getTime() - now.getTime();
    const daysLeft = Math.ceil(timeDiff / (1000 * 3600 * 24));

    return Math.max(0, daysLeft); // Don't show negative days
  }

  function getDaysLeftText(): string {
    const days = getDaysLeft();
    if (days === 0) return "Expires today";
    if (days === 1) return "1 day left";
    return `${days} days left`;
  }

  function isExpired(): boolean {
    if (!expires_at) return false;
    return new Date() > new Date(expires_at);
  }

  onMount(() => {
    if (showAnimation) {
      // Show the bubble with content that appears on hover with animation for few seconds
      showContentAnimation = true;

      // Hide the animation after 3 seconds
      setTimeout(() => {
        showContentAnimation = false;
      }, 1000);
    }
  });

  function showConversation(e: Event) {
    conversationOpen = true;
    const el = document.createElement("div");
    const conversationBoxWidth = 336;
    const conversationBoxHeight = 417;
    let x = wrapper.getBoundingClientRect().x;
    let y = wrapper.getBoundingClientRect().y;
    if (x + conversationBoxWidth <= window.innerWidth) {
      x += 40;
    } else {
      x = x - (conversationBoxWidth + 10);
    }
    if (y + conversationBoxHeight <= window.innerHeight) {
      y -= 30;
    } else {
      y -= conversationBoxHeight;
    }
    mount(Conversation, {
      props: {
        lat,
        expires_at: expires_at,
        long,
        onDelete: () => {
          document.body.removeChild(el);
          onDelete();
        },
        onCreate: () => {},
        threadId: id,
        threadUserId: user_id,
        coordinates: {
          x: x,
          y: y,
        },
        create: false,
        onClose: () => {
          conversationOpen = false;
          try {
            document.body.removeChild(el);
          } catch {}
        },
      },
      target: el,
    });

    document.body.appendChild(el);
    e.stopPropagation();
  }
</script>

<div class="conversation-bubble-wrapper" bind:this={wrapper}>
  <button
    class="conversation-bubble"
    class:active={conversationOpen}
    class:animate-expand={showContentAnimation}
    class:expired={isExpired()}
    onclick={showConversation}
  >
    <div class="avatar-wrapper">
      <Avatar iconSize={10} src={user_image} />
    </div>
    <div class="content">
      <p class="username">
        {username} <span class="time">{getTimeAgo(created_at)}</span>
      </p>
      <p class="last-chat"></p>
      <p>{message}</p>
      {#if replies > 0}
        <p class="time">
          {replies}
          {#if replies > 1}replies{:else}reply{/if}
        </p>
      {/if}
      {#if expires_at}
        <p class="expiration" class:expired-text={isExpired()}>
          {#if isExpired()}
            Expired
          {:else}
            {getDaysLeftText()}
          {/if}
        </p>
      {/if}
    </div>
  </button>
</div>

<style>
  .conversation-bubble-wrapper {
    position: relative;
  }
  .conversation-bubble {
    z-index: 1;
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
      width 70ms,
      max-height 200ms,
      padding 200ms;
  }
  .active {
    box-shadow: 0 0 0 2px var(--color-primary);
  }
  .expired {
    opacity: 0.6;
    filter: grayscale(50%);
  }
  .avatar-wrapper {
    flex-shrink: 0;
  }
  .conversation-bubble:not(.active):hover,
  .conversation-bubble.animate-expand {
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
  .conversation-bubble:not(.active):hover .content,
  .conversation-bubble.animate-expand .content {
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
  .expiration {
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--color-warning);
    margin-top: 2px;
  }
  .expired-text {
    color: var(--color-error);
  }
</style>
