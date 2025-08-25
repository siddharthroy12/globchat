<script lang="ts">
  import { onMount } from "svelte";
  import {
    queryUsers,
    type UserQueryResponse,
  } from "../../lib/services/users.svelte";
  import Avatar from "$lib/components/avatar.svelte";

  let searchQuery = "";
  let currentPage = 0;
  let pageSize = 20;
  let queryResult: UserQueryResponse = {
    users: [],
    pagination: {
      total: 0,
      count: 0,
      page: 0,
      page_size: 20,
      total_pages: 0,
      has_next: false,
      has_prev: false,
    },
  };
  let loading = false;
  let searchTimeout: number;

  // Calculate total pages from pagination
  $: totalPages = queryResult.pagination.total_pages;
  $: hasNextPage = queryResult.pagination.has_next;
  $: hasPrevPage = queryResult.pagination.has_prev;

  async function loadUsers() {
    if (loading) return;

    loading = true;
    try {
      queryResult = await queryUsers(
        searchQuery.trim() || undefined,
        pageSize,
        currentPage
      );
      console.log(queryResult);
    } catch (error) {
      console.error("Failed to load users:", error);
      queryResult = {
        users: [],
        pagination: {
          total: 0,
          count: 0,
          page: 0,
          page_size: 20,
          total_pages: 0,
          has_next: false,
          has_prev: false,
        },
      };
    } finally {
      loading = false;
    }
  }

  function handleSearch() {
    // Debounce search to avoid too many API calls
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      currentPage = 0; // Reset to first page when searching
      loadUsers();
    }, 300);
  }

  function goToPage(page: number) {
    if (page >= 0 && page < totalPages) {
      currentPage = page;
      loadUsers();
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

  onMount(() => {
    loadUsers();
  });
</script>

<div class="flex flex-col gap-6 items-end">
  <div class="flex justify-between w-full">
    <h1 class="text-2xl font-bold">Users</h1>
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
        placeholder="Search users..."
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
            <th class="bg-base-300">ID</th>
            <td>Avatar</td>
            <td>Username</td>
            <td>Email</td>
            <td>Messages</td>
            <th class="bg-base-300">Created At</th>
          </tr>
        </thead>
        <tbody class="bg-base-300">
          {#each queryResult.users as user (user.id)}
            <tr class="hover:bg-base-200">
              <th class="">{user.id}</th>
              <td>
                <Avatar src={user.image} size={40} />
              </td>
              <td>
                <div class="font-medium">{user.username}</div>
              </td>
              <td class="max-w-xs">
                <div class="tooltip" data-tip={user.email}>
                  {user.email}
                </div>
              </td>
              <td>
                <div class="badge badge-primary">
                  {user.messages}
                </div>
              </td>
              <th class="text-sm">
                {formatDate(user.created_at)}
              </th>
            </tr>
          {:else}
            <tr>
              <td colspan="6" class="text-center py-8 text-gray-500">
                {searchQuery
                  ? "No users found matching your search."
                  : "No users found."}
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
      Showing {queryResult.pagination.count} of {queryResult.pagination.total} users
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
