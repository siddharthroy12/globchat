<script>
  import Messages from "./messages.svelte";
  import Reports from "./reports.svelte";
  import Users from "./users.svelte";

  let selectedTab = "Users"; // Default to Users since it's checked
</script>

<div class="navbar bg-base-100 shadow-sm bg-primary">
  <div class="navbar-start">
    <a class="btn btn-ghost text-xl">Admin Panel</a>
  </div>
  <div class="navbar-center">
    <!-- name of each tab group should be unique -->
    <div class="tabs tabs-box">
      <input
        type="radio"
        name="my_tabs_1"
        class="tab"
        aria-label="Messages"
        bind:group={selectedTab}
        value="Messages"
      />
      <input
        type="radio"
        name="my_tabs_1"
        class="tab"
        aria-label="Users"
        bind:group={selectedTab}
        value="Users"
        checked
      />
      <input
        type="radio"
        name="my_tabs_1"
        class="tab"
        aria-label="Reports"
        bind:group={selectedTab}
        value="Reports"
      />
    </div>
  </div>
  <div class="navbar-end">
    <a class="btn" href="/">Open Map</a>
  </div>
</div>

<div class="p-6">
  {#if selectedTab === "Messages"}
    <Messages />
  {:else if selectedTab === "Users"}
    <Users />
  {:else if selectedTab === "Reports"}
    <Reports />
  {/if}
</div>
