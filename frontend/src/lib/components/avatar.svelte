<script lang="ts">
  import { UserRound } from "@lucide/svelte";

  const { size = 20, iconSize }: { size?: number; iconSize?: number } =
    $props();
</script>

<div class="avatar" style={`width: ${size}px; height: ${size}px`}>
  <UserRound size={iconSize} />
</div>

<style>
  .avatar {
    background: var(--color-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50px;
  }
</style>
