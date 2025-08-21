<script lang="ts">
  import { UserRound } from "@lucide/svelte";

  const {
    size = 20,
    iconSize,
    src,
  }: { size?: number; iconSize?: number; src?: string } = $props();
</script>

<div
  class="avatar"
  style={`width: ${size}px; height: ${size}px; background: ${src ? "url(" + src + ")" : "var(--color-primary)"}`}
>
  <UserRound size={iconSize} />
</div>

<style>
  .avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50px;
  }
</style>
