<script>
  import { getUserData, logout } from "$lib/services/auth.svelte";
  import { Pencil } from "@lucide/svelte";
  import Avatar from "./avatar.svelte";

  function closeModal() {
    // @ts-ignore
    document.activeElement.blur();
  }
  function onClickLogout() {
    logout();
    closeModal();
  }

  function showEditProfile() {
    // @ts-ignore
    document.getElementById("edit-profile-modal")?.showModal();
    closeModal();
  }
</script>

<div
  tabindex="0"
  class="dropdown-content card-body w-[366px] px-5 py-5 bg-base-100 rounded-4xl mb-6 flex flex-col gap-3 z-[999]"
>
  <button
    class="btn btn-sm btn-circle btn-soft absolute right-4 top-4"
    onclick={closeModal}>✕</button
  >
  <div class="flex gap-4 items-center">
    <div class="relative">
      <Avatar size={76} iconSize={30} />
      <button
        class="btn btn-soft btn-circle absolute bottom-[-5px] right-[-5px] scale-[90%]"
        onclick={showEditProfile}
      >
        <Pencil size={15} />
      </button>
    </div>
    <div class="flex flex-col gap-0">
      <p class="text-xl">
        {getUserData()?.username}
        <span class="text-secondary">#{getUserData()?.id}</span>
      </p>
      <p>Chats: <span class="text-primary">0</span></p>
    </div>
  </div>
  <button class="btn btn-soft rounded-full btn-error" onclick={onClickLogout}
    >Logout</button
  >
</div>
