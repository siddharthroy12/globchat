<script lang="ts">
  import { ArrowUp, Copy, Ellipsis, Image, Trash, X } from "@lucide/svelte";
  import {
    createThread,
    deleteThread,
    type Thread,
  } from "$lib/services/threads.svelte";
  import {
    createMessage,
    getMessages,
    type Message,
  } from "$lib/services/message.svelte";
  import { onMount, tick } from "svelte";
  import MessageBubble from "./message-bubble.svelte";
  import { getUserData } from "$lib/services/auth.svelte";
  import { joinThread } from "$lib/services/websocket";

  type ConversationProps = {
    coordinates: {
      x: number;
      y: number;
    };
    lat: number;
    long: number;
    onClose: () => void;
    onCreate: (thread: Thread) => void;
    onDelete: () => void;
    create: boolean;
    threadId?: number;
    threadUserId?: number;
  };

  let {
    coordinates,
    onClose,
    create,
    onCreate,
    lat,
    long,
    threadId,
    threadUserId,
    onDelete,
  }: ConversationProps = $props();

  let inputValue = $state("");
  let messages: Message[] = $state([]);
  let showFullLoading = $state(false);
  let showLoadingMoreMessages = $state(false);
  let showSendLoading = $state(false);
  let closeConnection = $state(() => {});
  let isSendButtonDisabled = $derived(inputValue.trim() == "");

  // References for scroll management
  // svelte-ignore non_reactive_update
  let messagesContainer: HTMLDivElement | undefined;
  let isUserScrolling = $state(false);
  let hasFirstMessage = $derived(messages.some((m) => m.is_first));

  onMount(() => {
    if (!create) {
      firstLoadMessage();
      openConnection();
    }
  });

  // Auto-scroll to bottom function
  async function scrollToBottom() {
    await tick();

    if (messagesContainer && !isUserScrolling) {
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

  function openConnection() {
    if (threadId) {
      closeConnection = joinThread(
        threadId,
        () => {
          // TODO: load more newer messages
        },
        () => {
          // If connection got disconnected for some reason connect again
          openConnection();
        }
      );
    }
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

  function openDeleteModal() {
    // @ts-ignore
    delete_confimation_modal.showModal();
  }

  async function onDeleteConfirmation() {
    if (threadId != undefined) {
      await deleteThread(threadId);
      onDelete();
    }
  }
  // Then add this function in your <script> section:
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      onClose();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<dialog
  id="delete_confimation_modal"
  class="modal modal-bottom sm:modal-middle"
>
  <div class="modal-box">
    <h3 class="text-lg font-bold">Are you sure?</h3>
    <p class="py-4">Are you sure you want to delete this thread?</p>
    <div class="modal-action">
      <form method="dialog">
        <button class="btn">Close</button>
        <button class="btn btn-error" onclick={onDeleteConfirmation}>Yes</button
        >
      </form>
    </div>
  </div>
</dialog>
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
            <div class="dropdown">
              <div tabindex="0" role="button" class="icon-btn m-1">
                <Ellipsis />
              </div>
              <ul
                tabindex="0"
                class="dropdown-content menu bg-base-100 rounded-box z-1 w-[160px] p-2 shadow-sm"
              >
                {#if threadUserId == getUserData()?.id}
                  <li class="text-error" onclick={openDeleteModal}>
                    <a><Trash size={14} />Delete Thread</a>
                  </li>
                {/if}
                <li><a><Copy size={14} />Copy Link</a></li>
              </ul>
            </div>
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

          {#each messages as message, index (message.id)}
            <MessageBubble
              {message}
              onDelete={() => {
                messages = messages.splice(index, 1);
              }}
            />
          {/each}
        </div>
      {/if}

      <div class="bottom p-3 flex gap-2 py-4">
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
