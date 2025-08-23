<script lang="ts">
  import { ArrowUp, Ellipsis, Image, X } from "@lucide/svelte";
  import Avatar from "./avatar.svelte";
  import { createThread, type Thread } from "$lib/services/threads.svelte";
  import {
    createMessage,
    getMessages,
    type Message,
  } from "$lib/services/message.svelte";
  import { onMount, tick } from "svelte";
  import MessageBubble from "./message-bubble.svelte";

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
    threadId?: number;
  };

  let {
    coordinates,
    onClose,
    create,
    onCreate,
    lat,
    long,
    threadId,
  }: ConversationProps = $props();

  let inputValue = $state("");
  let messages: Message[] = $state([]);
  let showFullLoading = $state(false);
  let showLoadingMoreMessages = $state(false);
  let showSendLoading = $state(false);
  let isSendButtonDisabled = $derived(inputValue.trim() == "");

  // References for scroll management
  let messagesContainer: HTMLDivElement | undefined;
  let isUserScrolling = $state(false);
  let hasFirstMessage = $derived(messages.some((m) => m.is_first));

  onMount(() => {
    console.log(messagesContainer); // Here it is not null

    if (!create) {
      firstLoadMessage();
    }
  });

  // Auto-scroll to bottom function
  async function scrollToBottom() {
    await tick();

    if (messagesContainer && !isUserScrolling) {
      console.log("scrolling");
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }

  // Handle scroll events for infinite loading
  function handleScroll(event: Event) {
    const target = event.target as HTMLDivElement;

    // Check if scrolled to top for loading older messages
    if (
      target.scrollTop === 0 &&
      !showLoadingMoreMessages &&
      !hasFirstMessage
    ) {
      const oldestMessage = messages[0];
      if (oldestMessage) {
        loadOlderMessages(oldestMessage.id);
      }
    }

    // Track if user is actively scrolling (prevent auto-scroll during user interaction)
    isUserScrolling = true;
    clearTimeout(scrollTimeout);
    scrollTimeout = setTimeout(() => {
      isUserScrolling = false;
    }, 150);
  }

  let scrollTimeout: number;

  async function firstLoadMessage() {
    showFullLoading = true;
    console.log(messagesContainer); //  Here it is not null

    try {
      if (threadId) {
        let m = await getMessages(threadId, 10);
        m = m.reverse();
        messages = m;
      }
    } catch {}

    showFullLoading = false;
    await scrollToBottom();
  }

  async function loadOlderMessages(oldestMessageId: number) {
    if (showLoadingMoreMessages || hasFirstMessage) return;

    showLoadingMoreMessages = true;
    const previousScrollHeight = messagesContainer!.scrollHeight;

    try {
      if (threadId) {
        let m = await getMessages(threadId, 10, oldestMessageId);
        m = m.reverse();
        messages = [...m, ...messages]; // Prepend older messages

        // Maintain scroll position after adding older messages
        await tick();
        const newScrollHeight = messagesContainer!.scrollHeight;
        messagesContainer!.scrollTop = newScrollHeight - previousScrollHeight;
      }
    } catch {}

    showLoadingMoreMessages = false;
  }

  async function sendMessage() {
    try {
      if (threadId) {
        const m = await createMessage(threadId, inputValue);
        messages = [...messages, m]; // Add new message
        // await scrollToBottom(); // Auto-scroll to show new message
      }
    } catch {}
  }

  async function onCreateThread() {
    const thread = await createThread(lat, long, inputValue);
    onCreate(thread);
  }

  async function onSend() {
    showSendLoading = true;
    if (create) {
      await onCreateThread();
    } else {
      await sendMessage();
      scrollToBottom();
    }
    inputValue = "";
    showSendLoading = false;
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
    {#if showFullLoading}
      <div class="flex justify-center items-center h-[300px]">
        <span class="loading loading-bars loading-xl"></span>
      </div>
    {:else}
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
        <div
          class="body p-3 h-[212px] overflow-y-auto"
          bind:this={messagesContainer}
          onscroll={handleScroll}
        >
          {#if showLoadingMoreMessages}
            <div class="flex justify-center py-2">
              <span class="loading loading-spinner loading-sm"></span>
            </div>
          {/if}

          {#if hasFirstMessage}
            <div class="text-center text-sm text-gray-500 mb-4 py-2">
              This is the start of the conversation
            </div>
          {/if}

          {#each messages as message}
            <MessageBubble {message} />
          {/each}
        </div>
      {/if}

      <div class="bottom p-3 flex gap-2 py-4">
        <div>
          <Avatar iconSize={10} />
        </div>
        <div class="flex flex-col w-full">
          <textarea
            class="textarea w-full text-input"
            placeholder="Write"
            disabled={showSendLoading}
            bind:value={inputValue}
          ></textarea>
          <div class="flex items-center justify-between p-2 bottom-buttons">
            <div class="flex items-center gap2">
              <!-- <button class="icon-btn">
              <Image size={16} />
            </button> -->
            </div>
            <button
              class="btn btn-circle btn-primary w-[24px] h-[24px] p-1"
              disabled={isSendButtonDisabled || showSendLoading}
              onclick={onSend}
            >
              {#if showSendLoading}
                <span class="loading loading-spinner loading-xs"></span>
              {:else}
                <ArrowUp />
              {/if}
            </button>
          </div>
        </div>
      </div>
    {/if}
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
