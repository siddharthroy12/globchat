<script lang="ts">
  import {
    Dices,
    DoorOpen,
    Info,
    LocateFixed,
    Minus,
    Plus,
  } from "@lucide/svelte";
  import LoginModal from "./modals/login-modal.svelte";
  import {
    AuthenticationStatus,
    getAuthenticationStatus,
    getUserData,
  } from "$lib/services/auth.svelte";
  import Avatar from "./avatar.svelte";
  import AccountInfoCard from "./account-info-card.svelte";
  import EditProfileModal from "./modals/edit-profile-modal.svelte";
  import AboutModal from "./modals/about-modal.svelte";
  type ControlsProps = {
    zoomToRandomChat: () => void;
    zoomToMyLocation: () => void;
    zoomIn: () => void;
    zoomOut: () => void;
    showOnlyUserThreads: boolean;
  };
  let {
    zoomToRandomChat,
    zoomToMyLocation,
    zoomIn,
    zoomOut,
    showOnlyUserThreads = $bindable(),
  }: ControlsProps = $props();

  function openLoginModal() {
    // @ts-ignore
    document.getElementById("login-modal")?.showModal();
  }

  function openAboutModal() {
    // @ts-ignore
    about_modal.showModal();
    // @ts-ignore
    document.querySelector("#about_modal .modal-box").scrollTop = 0;
  }
</script>

{#if getAuthenticationStatus() == AuthenticationStatus.LoggedIn}
  <EditProfileModal />
{/if}
<AboutModal />
<LoginModal />
<div class="toolbar">
  <div class="toolbar__container shadow-md">
    <button
      class="btn btn-primary btn-circle"
      title="App info"
      onclick={openAboutModal}
    >
      <Info />
    </button>
    {#if getAuthenticationStatus() == AuthenticationStatus.Unknown}
      <button
        class="btn btn-primary btn-circle"
        disabled
        title="Loading auth status"
        aria-label="Loading auth status"
      >
        <span class="loading loading-spinner loading-xl"></span>
      </button>
    {:else if getAuthenticationStatus() == AuthenticationStatus.LoggedIn}
      <div class="dropdown dropdown-top dropdown-center">
        <div tabindex="0" role="button" class="">
          <Avatar size={38} src={getUserData()?.image} />
        </div>
        <AccountInfoCard />
      </div>
    {:else}
      <button
        class="btn btn-primary rounded-full"
        title="Zoom to random chat"
        onclick={openLoginModal}
      >
        <DoorOpen />
        Log in
      </button>{/if}
  </div>
  <div class="toolbar__container shadow-md">
    <button
      class="btn btn-primary btn-circle"
      onclick={zoomToRandomChat}
      title="Zoom to random chat"
    >
      <Dices />
    </button>
    <button
      class="btn btn-primary btn-circle"
      title="Zoom to my location"
      onclick={zoomToMyLocation}
    >
      <LocateFixed />
    </button>
  </div>
  <div class="toolbar__container shadow-md">
    <button class="btn btn-primary btn-circle" onclick={zoomIn} title="Zoom In">
      <Plus />
    </button>
    <button
      class="btn btn-primary btn-circle"
      title="Zoom out"
      onclick={zoomOut}
    >
      <Minus />
    </button>
  </div>
  {#if getAuthenticationStatus() == AuthenticationStatus.LoggedIn}
    <div class="toolbar__container shadow-md">
      <div class="h-[40px] flex items-center gap-2 px-2">
        <input
          type="checkbox"
          class="toggle toggle-md"
          bind:checked={showOnlyUserThreads}
        />
        Show only mine
      </div>
    </div>
  {/if}
</div>

<style>
  .btn-circle {
    border-radius: 50px;
  }
  .toolbar {
    padding: 10px;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    display: flex;
    justify-content: center;
    gap: 10px;
    align-items: center;
  }
  .toolbar__container {
    padding: 10px;
    border-radius: 50px;
    background: var(--color-base-300);
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
  }
</style>
