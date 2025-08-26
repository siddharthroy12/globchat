<script lang="ts">
  import "../app.css";
  import "../../node_modules/maplibre-gl/src/css/maplibre-gl.css";
  import favicon from "$lib/assets/favicon.webp";
  import { onMount } from "svelte";
  import { PUBLIC_GOOGLE_CLIENT_ID } from "$env/static/public";
  import { checkAuthenticationStatus } from "$lib/services/auth.svelte";
  let { children } = $props();

  onMount(() => {
    checkAuthenticationStatus();
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
  <meta
    name="description"
    content="GlobeChat is a platform where you can start conversation anywhere on the world map and chat with people who are interested in the same location as you. Start topics anywhere and let other people find it and join the conversation."
  />
  <meta
    name="keywords"
    content="Map, Globe, World Map, Chats, Messages, Threads, Conversations"
  />
  <title>GlobeChat</title>
</svelte:head>

{@render children?.()}
