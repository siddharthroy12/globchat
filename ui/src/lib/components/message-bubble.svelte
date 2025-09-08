<script lang="ts">
  import {
    AuthenticationStatus,
    getAuthenticationStatus,
    getUserData,
  } from "$lib/services/auth.svelte";
  import type { Message } from "$lib/services/message.svelte";
  import { CircleAlert, Copy, Ellipsis, Trash } from "@lucide/svelte";
  import Avatar from "./avatar.svelte";
  import { mount } from "svelte";
  import MessageOptionDropdown from "./message-option-dropdown.svelte";
  import { extractLinks, linkify } from "$lib/helpers";

  type MessageProps = {
    message: Message;
    onDelete: () => void;
  };
  const { message, onDelete }: MessageProps = $props();
  let dropdownOpened = $state(false);
  let dropdownButton: HTMLElement | null = $state(null);
  let imageLoaded = $state(false);
  const isFromUser = $derived(getUserData()?.id == message.user_id);
  let firstLink = $derived(extractLinks(message.text)[0]);

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    const time = date.toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
    const day = date.toLocaleDateString("en-GB", {
      day: "2-digit",
      month: "2-digit",
      year: "2-digit",
    });

    return `${day} ${time}`;
  };

  function openDropdown(e: Event) {
    dropdownOpened = true;
    const el = document.createElement("div");
    let x = dropdownButton!.getBoundingClientRect().x;
    let y = dropdownButton!.getBoundingClientRect().y;
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
    <time class="text-xs opacity-50">{formatDateTime(message.created_at)}</time>
  </div>
  <div class="chat-bubble bg-primary text-white text-wrap break-all">
    {@html linkify(message.text)}

    {#if firstLink}
      <img
        src={firstLink}
        class="mt-2 rounded-lg max-w-xs block w-full"
        class:hidden={!imageLoaded}
        onload={() => (imageLoaded = true)}
        onerror={() => (firstLink = "")}
      />
    {/if}
  </div>

  {#if getAuthenticationStatus() === AuthenticationStatus.LoggedIn}
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
  {/if}
</div>

<style>
  .message-options {
    border: 2px solid var(--color-base-300);
  }
</style>
