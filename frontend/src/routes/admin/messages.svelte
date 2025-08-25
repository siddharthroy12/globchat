<script lang="ts">
  import { onMount } from "svelte";
  import {
    queryMessages,
    type MessageQueryResult,
  } from "../../lib/services/message.svelte";

  let searchQuery = "";
  let currentPage = 0;
  let pageSize = 20;
  let queryResult: MessageQueryResult = { total: 0, count: 0, messages: [] };
  let loading = false;
  let searchTimeout: number;

  // Calculate total pages
  $: totalPages = Math.ceil(queryResult.total / pageSize);
  $: hasNextPage = currentPage < totalPages - 1;
  $: hasPrevPage = currentPage > 0;

  async function loadMessages() {
    if (loading) return;

    loading = true;
    try {
      queryResult = await queryMessages(
        searchQuery.trim() || undefined,
        pageSize,
        currentPage
      );
    } catch (error) {
      console.error("Failed to load messages:", error);
      queryResult = { total: 0, count: 0, messages: [] };
    } finally {
      loading = false;
    }
  }

  function handleSearch() {
    // Debounce search to avoid too many API calls
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      currentPage = 0; // Reset to first page when searching
      loadMessages();
    }, 300);
  }

  function goToPage(page: number) {
    if (page >= 0 && page < totalPages) {
      currentPage = page;
      loadMessages();
    }
  }

  function nextPage() {
    if (hasNextPage) {
      goToPage(currentPage + 1);
    }
  }

  function prevPage() {
    if (hasPrevPage) {
      goToPage(currentPage - 1);
    }
  }

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  function truncateText(text: string, maxLength = 50) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + "...";
  }

  onMount(() => {
    loadMessages();
  });
</script>

<div class="flex flex-col gap-6 items-end">
  <div class="flex justify-between w-full">
    <h1 class="text-2xl font-bold">Messages</h1>
    <label class="input input-bordered flex items-center gap-2">
      <svg
        class="h-[1em] opacity-50"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
      >
        <g
          stroke-linejoin="round"
          stroke-linecap="round"
          stroke-width="2.5"
          fill="none"
          stroke="currentColor"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.3-4.3"></path>
        </g>
      </svg>
      <input
        type="search"
        placeholder="Search messages..."
        bind:value={searchQuery}
        on:input={handleSearch}
        class="grow"
      />
    </label>
  </div>

  {#if loading}
    <div class="flex justify-center w-full py-8">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    <div
      class="overflow-x-auto bg-base-300 h-[calc(100dvh-250px)] border rounded-lg border-base-300 w-full"
    >
      <table class="table table-pin-rows bg-base-300">
        <thead class="bg-base-300">
          <tr class="bg-base-300">
            <th>ID</th>
            <td>Text</td>
            <td>Image</td>
            <td>Thread ID</td>
            <td>Is First</td>
            <td>User</td>
            <th class="bg-base-300">Created At</th>
          </tr>
        </thead>
        <tbody class="bg-base-300">
          {#each queryResult.messages as message (message.id)}
            <tr class="hover:bg-base-200">
              <th>{message.id}</th>
              <td class="max-w-xs">
                <div class="tooltip" data-tip={message.text}>
                  {truncateText(message.text)}
                </div>
              </td>
              <td>
                {#if message.image}
                  <div class="avatar">
                    <div class="w-12 h-12 rounded">
                      <img src={message.image} alt="Message attachment" />
                    </div>
                  </div>
                {:else}
                  <span class="text-gray-500">None</span>
                {/if}
              </td>
              <td>
                <a
                  href="/?threadId={message.thread_id}"
                  class="link link-primary"
                  target="_blank"
                >
                  {message.thread_id}
                </a>
              </td>

              <td>
                <div
                  class="badge {message.is_first
                    ? 'badge-primary'
                    : 'badge-outline'}"
                >
                  {message.is_first ? "Yes" : "No"}
                </div>
              </td>
              <td>
                <div class="flex items-center gap-2">
                  {#if message.user_image}
                    <div class="avatar">
                      <div class="w-8 h-8 rounded-full">
                        <img src={message.user_image} alt={message.username} />
                      </div>
                    </div>
                  {/if}
                  <div>
                    <div class="font-medium">{message.username}</div>
                    <div class="text-sm text-gray-500">
                      ID: {message.user_id}
                    </div>
                  </div>
                </div>
              </td>
              <th class="text-sm">{formatDate(message.created_at)}</th>
            </tr>
          {:else}
            <tr>
              <td colspan="8" class="text-center py-8 text-gray-500">
                {searchQuery
                  ? "No messages found matching your search."
                  : "No messages found."}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}

  <!-- Pagination -->
  <div class="flex items-center gap-4">
    <div class="text-sm text-gray-500">
      Showing {queryResult.count} of {queryResult.total} messages
      {#if totalPages > 0}
        (Page {currentPage + 1} of {totalPages})
      {/if}
    </div>

    <div class="join">
      <button
        class="join-item btn btn-sm"
        class:btn-disabled={!hasPrevPage}
        on:click={prevPage}
        disabled={!hasPrevPage}
      >
        «
      </button>

      <!-- Show page numbers -->
      {#each Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
        const startPage = Math.max(0, Math.min(currentPage - 2, totalPages - 5));
        return startPage + i;
      }) as pageNum}
        <button
          class="join-item btn btn-sm"
          class:btn-active={pageNum === currentPage}
          on:click={() => goToPage(pageNum)}
        >
          {pageNum + 1}
        </button>
      {/each}

      <button
        class="join-item btn btn-sm"
        class:btn-disabled={!hasNextPage}
        on:click={nextPage}
        disabled={!hasNextPage}
      >
        »
      </button>
    </div>
  </div>
</div>
