<script lang="ts">
  import { Dices, LocateFixed, Minus, Plus } from "@lucide/svelte";
  type ControlsProps = {
    zoomToRandomChat: () => void;
    zoomToMyLocation: () => void;
    zoomIn: () => void;
    zoomOut: () => void;
  };
  let { zoomToRandomChat, zoomToMyLocation, zoomIn, zoomOut }: ControlsProps =
    $props();
</script>

<div class="toolbar">
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
