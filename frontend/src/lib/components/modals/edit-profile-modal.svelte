<script lang="ts">
  import { Pencil } from "@lucide/svelte";
  import Avatar from "../avatar.svelte";
  import { onMount } from "svelte";
  import { getUserData } from "$lib/services/auth.svelte";

  let username = $state("");

  $effect(() => {
    username = getUserData()!.username;
  });

  function uploadPicture() {}

  function closeModal() {
    // @ts-ignore
    document.getElementById("edit-profile-modal")?.close();
  }
</script>

<dialog id="edit-profile-modal" class="modal">
  <div class="modal-box relative">
    <form method="dialog">
      <button
        class="btn btn-sm btn-circle btn-soft absolute right-2 top-2"
        onclick={closeModal}>✕</button
      >
    </form>
    <div class="flex flex-col gap-4">
      <div class="flex justify-center items-center">
        <h3 class="font-bold text-xl">Edit Profile</h3>
      </div>
    </div>
    <div class="flex gap-6 items-center">
      <div class="relative">
        <Avatar size={76} iconSize={30} />
        <button
          class="btn btn-soft btn-circle absolute bottom-[-5px] right-[-5px] scale-[90%]"
          onclick={uploadPicture}
        >
          <Pencil size={15} />
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
      <button class="btn btn-primary rounded-full">Save</button>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
