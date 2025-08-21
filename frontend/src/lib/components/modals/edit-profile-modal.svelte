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
  let uploadMessage = $state("a");
  let messageType = $state<"success" | "error" | "">("");

  // File input reference
  let fileInput: HTMLInputElement;

  $effect(() => {
    username = getUserData()!.username;
    image = getUserData()!.image;
  });

  async function uploadToImgur(file: File): Promise<string> {
    const formData = new FormData();
    formData.append("image", file);

    const response = await fetch("https://api.imgur.com/3/image", {
      method: "POST",
      headers: {
        Authorization: "Client-ID YOUR_IMGUR_CLIENT_ID", // Replace with your Imgur Client ID
      },
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.statusText}`);
    }

    const data = await response.json();

    if (!data.success) {
      throw new Error(data.data?.error || "Upload failed");
    }

    return data.data.link;
  }

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

    isUploading = true;
    uploadMessage = "Uploading image...";
    messageType = "";

    try {
      const imageUrl = await uploadToImgur(file);
      image = imageUrl;
      showMessage("Image uploaded successfully!", "success");
    } catch (error) {
      console.error("Upload error:", error);
      showMessage("Upload failed", "error");
    } finally {
      isUploading = false;
      // Clear file input
      target.value = "";
    }
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
    try {
      await updateUserImageAndUsername(username, image);
      showMessage("Profile updated successfully!", "success");

      // Close modal after successful update
      setTimeout(() => {
        closeModal();
      }, 1500);
    } catch (error) {
      showMessage("Failed to update profile", "error");
    }
  }

  function closeModal() {
    // Clear any messages when closing
    uploadMessage = "";
    messageType = "";
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
        Save
      </button>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
