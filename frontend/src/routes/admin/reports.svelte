<script lang="ts">
  import { onMount } from "svelte";
  import MessagePreviewModal from "$lib/components/modals/message-preview-modal.svelte";
  import {
    queryReports,
    deleteReport,
    resolveReport,
    type ReportQueryResult,
  } from "../../lib/services/reports.svelte";

  let searchQuery = "";
  let currentPage = 0;
  let pageSize = 20;
  let queryResult: ReportQueryResult = { total: 0, count: 0, reports: [] };
  let loading = false;
  let searchTimeout: number;
  let selectedMessageId = 0;

  function openMessageModal(messageId: number) {
    selectedMessageId = messageId;
    // @ts-ignore
    message_preview_modal?.showModal();
  }

  // Calculate total pages
  $: totalPages = Math.ceil(queryResult.total / pageSize);
  $: hasNextPage = currentPage < totalPages - 1;
  $: hasPrevPage = currentPage > 0;

  async function loadReports() {
    if (loading) return;

    loading = true;
    try {
      queryResult = await queryReports(
        searchQuery.trim() || undefined,
        pageSize,
        currentPage
      );
    } catch (error) {
      console.error("Failed to load reports:", error);
      queryResult = { total: 0, count: 0, reports: [] };
    } finally {
      loading = false;
    }
  }

  function handleSearch() {
    // Debounce search to avoid too many API calls
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      currentPage = 0; // Reset to first page when searching
      loadReports();
    }, 300);
  }

  function goToPage(page: number) {
    if (page >= 0 && page < totalPages) {
      currentPage = page;
      loadReports();
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

  async function handleDeleteReport(reportId: number) {
    if (!confirm("Are you sure you want to delete this report?")) return;

    try {
      await deleteReport(reportId);
      await loadReports(); // Reload the list
    } catch (error) {
      console.error("Failed to delete report:", error);
      alert("Failed to delete report. Please try again.");
    }
  }

  async function handleResolveReport(reportId: number) {
    if (!confirm("Are you sure you want to resolve this report?")) return;

    try {
      await resolveReport(reportId);
      await loadReports(); // Reload the list
    } catch (error) {
      console.error("Failed to resolve report:", error);
      alert("Failed to resolve report. Please try again.");
    }
  }

  onMount(() => {
    loadReports();
  });
</script>

<div class="flex flex-col gap-6 items-end">
  <div class="flex justify-between w-full">
    <h1 class="text-2xl font-bold">Reports</h1>
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
        placeholder="Search reports..."
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
            <td>Reason</td>
            <td>Reporter ID</td>
            <td>Message ID</td>
            <td>Created At</td>
            <th class="bg-base-300">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-base-300">
          {#each queryResult.reports as report (report.id)}
            <tr class="hover:bg-base-200">
              <th>{report.id}</th>
              <td class="max-w-xs">
                <div class="tooltip" data-tip={report.reason}>
                  {truncateText(report.reason)}
                </div>
              </td>
              <td>
                <span class="font-mono text-sm">{report.reporter_id}</span>
              </td>
              <td>
                <button
                  class="link link-primary"
                  on:click={() => openMessageModal(report.message_id)}
                >
                  {report.message_id}
                </button>
              </td>
              <td class="text-sm">{formatDate(report.created_at)}</td>
              <th>
                <div class="flex gap-2">
                  <button
                    class="btn btn-sm btn-success"
                    on:click={() => handleResolveReport(report.id)}
                    title="Resolve Report"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M5 13l4 4L19 7"
                      />
                    </svg>
                  </button>
                  <button
                    class="btn btn-sm btn-error"
                    on:click={() => handleDeleteReport(report.id)}
                    title="Delete Report"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>
              </th>
            </tr>
          {:else}
            <tr>
              <td colspan="6" class="text-center py-8 text-gray-500">
                {searchQuery
                  ? "No reports found matching your search."
                  : "No reports found."}
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
      Showing {queryResult.count} of {queryResult.total} reports
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

<MessagePreviewModal id={selectedMessageId} />
