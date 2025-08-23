<script lang="ts">
  import { getUserData } from "$lib/services/auth.svelte";
  import {
    deleteMessage,
    reportMessage,
    type Message,
  } from "$lib/services/message.svelte";
  import { CircleAlert, Trash } from "@lucide/svelte";

  type MessageOptionDropdownProps = {
    coordinates: {
      x: number;
      y: number;
    };
    onClose: () => void;
    onDelete: () => void;
    message: Message;
  };

  const {
    coordinates,
    onClose,
    message,
    onDelete,
  }: MessageOptionDropdownProps = $props();

  let deleteModal: HTMLDialogElement;
  let reportModal: HTMLDialogElement;

  function showDeleteConfirmation() {
    deleteModal.showModal();
  }

  function showReportConfirmation() {
    reportModal.showModal();
  }

  function onConfirmDelete() {
    deleteMessage(message.id);
    deleteModal.close();
    onDelete();
    onClose();
  }

  function onConfirmReport() {
    reportMessage(message.id);
    reportModal.close();
    onClose();
  }

  function onCancelDelete() {
    deleteModal.close();
  }

  function onCancelReport() {
    reportModal.close();
  }
</script>

<div
  class="w-[100vw] height-[100dvh] fixed top-0 left-0 right-0 bottom-0 overflow-hidden z-[999]"
  onclick={onClose}
>
  <div
    class="bg-base-100 shadow-md menu fixed rounded-md w-[166px] p-1 border border-base-300"
    style={`top:${coordinates.y}px;left:${coordinates.x}px`}
    onclick={(e) => {
      e.stopPropagation();
    }}
  >
    {#if message.user_id == getUserData()?.id && !message.is_first}
      <li class="text-error">
        <button onclick={showDeleteConfirmation}>
          <Trash size={14} />Delete Message
        </button>
      </li>
    {/if}
    <li>
      <button onclick={showReportConfirmation}>
        <CircleAlert size={14} />Report Message
      </button>
    </li>
  </div>
</div>

<!-- Delete Confirmation Modal -->
<dialog bind:this={deleteModal} class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg text-error">Delete Message</h3>
    <p class="py-4">
      Are you sure you want to delete this message? This action cannot be
      undone.
    </p>
    <div class="modal-action">
      <button class="btn" onclick={onCancelDelete}>Cancel</button>
      <button class="btn btn-error" onclick={onConfirmDelete}>
        <Trash size={16} />
        Delete
      </button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button onclick={onCancelDelete}>close</button>
  </form>
</dialog>

<!-- Report Confirmation Modal -->
<dialog bind:this={reportModal} class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg text-warning">Report Message</h3>
    <p class="py-4">
      Are you sure you want to report this message? This will notify moderators
      for review.
    </p>
    <div class="modal-action">
      <button class="btn" onclick={onCancelReport}>Cancel</button>
      <button class="btn btn-warning" onclick={onConfirmReport}>
        <CircleAlert size={16} />
        Report
      </button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button onclick={onCancelReport}>close</button>
  </form>
</dialog>
