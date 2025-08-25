<script lang="ts">
  import { getMessageById, type Message } from "$lib/services/message.svelte";

  const { id }: { id: number } = $props();

  let fetching = $state(false);
  let message: Message | null = $state(null);

  async function fetchMessage() {
    fetching = true;
    try {
      message = await getMessageById(id);
    } catch (e) {
      console.error("Failed to fetch message", e);
    } finally {
      fetching = false;
    }
  }

  $effect(() => {
    fetchMessage();
  });
</script>

<dialog id="message_preview_modal" class="modal">
  <div class="modal-box relative rounded-4xl">
    <form method="dialog">
      <button class="btn btn-sm btn-circle btn-soft absolute right-4 top-4">
        ✕
      </button>
    </form>

    <div class="flex flex-col gap-4">
      <div class="flex justify-center items-center">
        <h3 class="font-bold text-xl">Message</h3>
      </div>

      {#if fetching}
        <div class="flex justify-center items-center py-8">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
      {:else if message}
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-3">
            <img
              src={message.user_image}
              alt={message.username}
              class="w-10 h-10 rounded-full"
            />
            <div>
              <p class="font-semibold">{message.username}</p>
              <p class="text-xs text-gray-500">
                {new Date(message.created_at).toLocaleString()}
              </p>
            </div>
          </div>

          {#if message.image}
            <div class="flex justify-center">
              <img
                src={message.image}
                alt="Message attachment"
                class="rounded-xl max-h-64 object-contain"
              />
            </div>
          {/if}

          <p class="text-base">{message.text}</p>
        </div>
      {:else}
        <p class="text-center text-gray-500">No message found</p>
      {/if}
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
