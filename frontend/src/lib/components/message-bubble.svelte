<script lang="ts">
  import { getUserData } from "$lib/services/auth.svelte";
  import type { Message } from "$lib/services/message.svelte";
  import { CircleAlert, Copy, Ellipsis, Trash } from "@lucide/svelte";
  import Avatar from "./avatar.svelte";
  import { mount } from "svelte";
  import MessageOptionDropdown from "./message-option-dropdown.svelte";

  type MessageProps = {
    message: Message;
    onDelete: () => void;
  };
  const { message, onDelete }: MessageProps = $props();
  let dropdownOpened = $state(false);
  let dropdownButton: HTMLElement;
  const isFromUser = $derived(getUserData()?.id == message.user_id);

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
  };

  function openDropdown(e: Event) {
    dropdownOpened = true;
    const el = document.createElement("div");
    let x = dropdownButton.getBoundingClientRect().x;
    let y = dropdownButton.getBoundingClientRect().y;
    mount(MessageOptionDropdown, {
      props: {
        message,
        coordinates: {
          x: x,
          y: y,
        },
        onDelete: onDelete,
        onClose: () => {
          dropdownOpened = false;
          document.body.removeChild(el);
        },
      },
      target: el,
    });

    document.body.appendChild(el);
    e.stopPropagation();
  }
</script>

<div
  class="chat relative hover:bg-base-100 group"
  class:chat-end={isFromUser}
  class:chat-start={!isFromUser}
  class:bg-base-100={dropdownOpened}
>
  <div class="chat-image avatar">
    <Avatar src={message.user_image} size={40} />
  </div>
  <div class="chat-header">
    {message.username}
    <time class="text-xs opacity-50">{formatTime(message.created_at)}</time>
  </div>
  <div class="chat-bubble bg-primary text-white">{message.text}</div>

  <div
    class="absolute left-0 top-0 translate-y-[-50%] p-1 bg-base-100 rounded-xl transition-opacity message-options"
    class:opacity-0={!dropdownOpened}
    class:opacity-100={dropdownOpened}
    class:group-hover:opacity-100={!dropdownOpened}
  >
    <button
      bind:this={dropdownButton}
      class="btn btn-ghost rounded-sm w-[25px] h-[25px] p-0"
      onclick={openDropdown}
    >
      <Ellipsis size={18} />
    </button>
  </div>
</div>

<style>
  .message-options {
    border: 2px solid var(--color-base-300);
  }
</style>
