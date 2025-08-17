<script lang="ts">
  import "../app.css";
  import "../../node_modules/maplibre-gl/src/css/maplibre-gl.css";
  import favicon from "$lib/assets/favicon.webp";

  let { children } = $props();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
