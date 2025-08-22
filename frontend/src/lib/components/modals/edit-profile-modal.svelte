<script lang="ts">
  import { Pencil } from "@lucide/svelte";
  import Avatar from "../avatar.svelte";
  import {
    getUserData,
    updateUserImageAndUsername,
  } from "$lib/services/auth.svelte";

  let username = $state("");
  let image = $state("");
  let isUploading = $state(false);
  let uploadMessage = $state("");
  let messageType = $state<"success" | "error" | "">("");

  // File input reference
  let fileInput: HTMLInputElement;
  let selectedFile: File | null = null;

  $effect(() => {
    username = getUserData()!.username;
    image = getUserData()!.image;
  });

  function uploadPicture() {
    fileInput.click();
  }

  async function handleFileSelect(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];

    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      showMessage("Please select a valid image file", "error");
      return;
    }

    // Validate file size (10MB limit)
    if (file.size > 10 * 1024 * 1024) {
      showMessage("Image size must be less than 10MB", "error");
      return;
    }

    selectedFile = file;

    // Create preview URL for immediate display
    const previewUrl = URL.createObjectURL(file);
    image = previewUrl;

    // Clear file input
    target.value = "";
  }

  function showMessage(message: string, type: "success" | "error") {
    uploadMessage = message;
    messageType = type;

    // Clear message after 5 seconds
    setTimeout(() => {
      uploadMessage = "";
      messageType = "";
    }, 5000);
  }

  async function updateInfo() {
    if (!username.trim()) {
      showMessage("Username is required", "error");
      return;
    }

    if (username.trim().length > 16) {
      showMessage("Username must be 16 characters or less", "error");
      return;
    }

    isUploading = true;
    uploadMessage = "Updating profile...";
    messageType = "";

    try {
      const response = await updateUserImageAndUsername(
        selectedFile,
        username.trim()
      );

      if (response.ok) {
        // Clear the selected file
        selectedFile = null;

        showMessage("Profile updated successfully!", "success");

        // Close modal after successful update
        setTimeout(() => {
          closeModal();
        }, 1500);
      } else {
        const errorData = await response.json().catch(() => ({}));
        showMessage(errorData.error || "Failed to update profile", "error");
      }
    } catch (error) {
      console.error("Update error:", error);
      showMessage("Failed to update profile", "error");
    } finally {
      isUploading = false;
    }
  }

  function closeModal() {
    // Clean up preview URL if exists
    if (image.startsWith("blob:")) {
      URL.revokeObjectURL(image);
      // Reset to original image
      image = getUserData()!.image;
    }

    // Clear selected file
    selectedFile = null;

    // Clear any messages when closing
    uploadMessage = "";
    messageType = "";

    // Reset username to original
    username = getUserData()!.username;

    // @ts-ignore
    document.getElementById("edit-profile-modal")?.close();
  }
</script>

<!-- Hidden file input -->
<input
  bind:this={fileInput}
  type="file"
  accept="image/*"
  style="display: none;"
  onchange={handleFileSelect}
/>

<dialog id="edit-profile-modal" class="modal rounded-4xl">
  <div class="modal-box relative rounded-4xl">
    <form method="dialog">
      <button
        class="btn btn-sm btn-circle btn-soft absolute right-4 top-4"
        onclick={closeModal}>✕</button
      >
    </form>
    <div class="flex flex-col gap-4">
      <div class="flex justify-center items-center">
        <h3 class="font-bold text-xl">Edit Profile</h3>
      </div>
    </div>

    <!-- Message Display -->
    {#if uploadMessage}
      <div
        class="alert rounded-full alert-soft my-3 {messageType === 'success'
          ? 'alert-success'
          : messageType === 'error'
            ? 'alert-error'
            : 'alert-info'}"
      >
        <span>{uploadMessage}</span>
      </div>
    {/if}

    <div class="flex gap-6 items-center">
      <div class="relative">
        <Avatar size={76} iconSize={30} src={image} />
        <button
          class="btn btn-soft btn-circle absolute bottom-[-5px] right-[-5px] scale-[90%] {isUploading
            ? 'loading'
            : ''}"
          onclick={uploadPicture}
          disabled={isUploading}
        >
          {#if !isUploading}
            <Pencil size={15} />
          {/if}
        </button>
      </div>
      <div class="flex flex-col gap-0 w-full">
        <label class="input w-full">
          <span class="label">Username</span>
          <input
            type="text"
            placeholder="Username"
            class="w-full"
            bind:value={username}
            maxlength="16"
          />
        </label>
      </div>
    </div>
    <div class="flex items-center justify-end gap-3">
      <button class="btn btn-soft rounded-full" onclick={closeModal}
        >Cancel</button
      >
      <button
        class="btn btn-primary rounded-full"
        onclick={updateInfo}
        disabled={isUploading}
      >
        {isUploading ? "Saving..." : "Save"}
      </button>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button onclick={closeModal}>close</button>
  </form>
</dialog>
