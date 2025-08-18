<script lang="ts">
  import "../app.css";
  import "../../node_modules/maplibre-gl/src/css/maplibre-gl.css";
  import favicon from "$lib/assets/favicon.webp";
  import { onMount } from "svelte";
  import { PUBLIC_GOOGLE_CLIENT_ID } from "$env/static/public";
  let { children } = $props();

  onMount(() => {
    //@ts-ignore
    google.accounts.id.initialize({
      client_id: PUBLIC_GOOGLE_CLIENT_ID,
      callback: (e: any) => {
        //@ts-ignore
        if (window.googlelogincallbacks) {
          //@ts-ignore
          window.googlelogincallbacks.forEach((callback) => {
            callback(e);
          });
        }
      },
    });
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
