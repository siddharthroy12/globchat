<script lang="ts">
  import { getUserData } from "$lib/services/auth.svelte";
  import { deleteMessage, type Message } from "$lib/services/message.svelte";
  import { createReport } from "$lib/services/reports.svelte";
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
  let reason: string = $state("");
  let reportStatus: "idle" | "submitting" | "already_exists" | "success" =
    $state("idle");

  function showDeleteConfirmation() {
    deleteModal.showModal();
  }

  function showReportConfirmation() {
    // Reset state when opening modal
    reason = "";
    reportStatus = "idle";
    reportModal.showModal();
  }

  function onConfirmDelete() {
    deleteMessage(message.id);
    deleteModal.close();
    onDelete();
    onClose();
  }

  async function onConfirmReport() {
    if (!reason.trim()) {
      return; // Don't submit if reason is empty
    }

    reportStatus = "submitting";

    try {
      const res = await createReport(message.id, reason.trim());
      if (res.includes("already exists")) {
        reportStatus = "already_exists";
      } else {
        reportStatus = "success";
        // Auto close after success
        setTimeout(() => {
          reportModal.close();
          onClose();
        }, 1500);
      }
    } catch (error) {
      console.log(error);
      reportStatus = "idle"; // Reset on error
      console.error("Failed to create report:", error);
    }
  }

  function onCancelDelete() {
    deleteModal.close();
  }

  function onCancelReport() {
    onClose();
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

<!-- Report Modal with Reason Input -->
<dialog bind:this={reportModal} class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg text-warning">Report Message</h3>

    {#if reportStatus === "idle" || reportStatus === "submitting"}
      <p class="py-2 mb-2">
        Please provide a reason for reporting this message. This will help
        moderators review the content appropriately.
      </p>

      <div class="form-control w-full flex flex-col">
        <label class="label" for="report-reason mb-2">
          <span class="label-text">Reason for report</span>
        </label>
        <textarea
          id="report-reason"
          bind:value={reason}
          class="textarea textarea-bordered h-24 resize-none w-full mt-2"
          placeholder="e.g., Spam, inappropriate content, harassment..."
          disabled={reportStatus === "submitting"}
        ></textarea>
      </div>

      <div class="modal-action">
        <button
          class="btn"
          onclick={onCancelReport}
          disabled={reportStatus === "submitting"}
        >
          Cancel
        </button>
        <button
          class="btn btn-warning"
          onclick={onConfirmReport}
          disabled={reportStatus === "submitting" || !reason.trim()}
        >
          {#if reportStatus === "submitting"}
            <span class="loading loading-spinner loading-sm"></span>
            Submitting...
          {:else}
            <CircleAlert size={16} />
            Report
          {/if}
        </button>
      </div>
    {:else if reportStatus === "already_exists"}
      <div class="alert alert-info mt-6">
        <CircleAlert size={16} />
        <span
          >This message has already been reported. Moderators are aware and will
          review it.</span
        >
      </div>

      <div class="modal-action">
        <button class="btn" onclick={onCancelReport}> Close </button>
      </div>
    {:else if reportStatus === "success"}
      <div class="alert alert-success mt-6">
        <CircleAlert size={16} />
        <span
          >Report submitted successfully. Moderators will review this message.</span
        >
      </div>

      <div class="modal-action">
        <button class="btn" onclick={onCancelReport}> Close </button>
      </div>
    {/if}
  </div>
  <form method="dialog" class="modal-backdrop">
    <button onclick={onCancelReport}>close</button>
  </form>
</dialog>
