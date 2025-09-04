<script lang="ts">
  import Chat from "../lib/components/chat.svelte";
  import maplibregl from "maplibre-gl";
  import { mount, onMount } from "svelte";
  import Controls from "../lib/components/controls.svelte";
  import Addchat from "../lib/components/addchat.svelte";
  import type { Thread } from "$lib/services/threads.svelte";
  import {
    fetchRandomThread,
    fetchThread,
    fetchThreads,
  } from "$lib/services/threads.svelte";
  import { DoorOpen } from "@lucide/svelte";
  import {
    AuthenticationStatus,
    getAuthenticationStatus,
    getUserData,
  } from "$lib/services/auth.svelte";
  import ThreadTooCloseModal from "$lib/components/modals/thread-too-close-modal.svelte";

  let map: null | maplibregl.Map = null;
  let mountedComponents: Map<
    number,
    { marker: maplibregl.Marker; mount: ReturnType<typeof mount> }
  > = new Map();

  let addChatComponent: {
    marker: maplibregl.Marker;
    mount: ReturnType<typeof mount>;
  } | null = null;

  let mapZoomLevel = $state(0);
  let loadedThreads: Thread[] = [];
  let isLoadingThreads = false;
  let tooManyItems = $state(false);

  onMount(() => {
    map = new maplibregl.Map({
      style: "https://tiles.openfreemap.org/styles/liberty",
      center: [0, 0],
      zoom: 0,
      container: "map",
      attributionControl: false,
      doubleClickZoom: false,
    });

    map.on("zoom", handleMapUpdate);
    map.on("move", handleMapUpdate);
    map.on("moveend", handleMapUpdate);
    map.on("click", handleMapClick);

    handleMapUpdate();
    loadTheadFromQueryParameter();
  });

  async function loadTheadFromQueryParameter() {
    const urlParams = new URLSearchParams(window.location.search);
    const threadId = urlParams.get("threadId");
    if (threadId) {
      const thread = await fetchThread(+threadId);
      zoomTo(thread.long, thread.lat);
    }
  }

  function handleMapClick(e: maplibregl.MapMouseEvent) {
    if (!map) return;

    const { lng, lat } = e.lngLat;
    removeAddChatComponent();
    mountAddChatComponent(lng, lat);
  }

  function mountAddChatComponent(long: number, lat: number) {
    if (!map || addChatComponent) return;

    const componentDom = document.createElement("div");
    componentDom.addEventListener("mouseenter", () => {
      componentDom.style.zIndex = "9999";
    });
    componentDom.addEventListener("mouseleave", () => {
      componentDom.style.zIndex = "1";
    });

    const componentMount = mount(Addchat, {
      target: componentDom,
      props: {
        lat,
        long: long,
        onCreate: (thread: Thread) => {
          loadChatComponent(thread, true);
          loadedThreads = [...loadedThreads, thread];
          removeAddChatComponent();
        },
        onClose: () => {
          removeAddChatComponent();
        },
      },
    });

    const marker = new maplibregl.Marker({
      element: componentDom,
      anchor: "bottom",
    })
      .setLngLat([long, lat])
      .addTo(map);

    addChatComponent = { marker, mount: componentMount };
  }

  function removeAddChatComponent() {
    if (!addChatComponent) return;
    addChatComponent.marker.remove();
    addChatComponent = null;
  }

  function updateMapCursor() {
    if (!map) return;
    const canvas = map.getCanvas();
    canvas.style.cursor = "crosshair";
  }

  async function handleMapUpdate() {
    if (!map) return;

    mapZoomLevel = map.getZoom();
    updateMapCursor();

    await loadVisibleThreads();
  }

  async function loadVisibleThreads() {
    if (!map || isLoadingThreads) return;
    isLoadingThreads = true;

    try {
      const bounds = map.getBounds();
      const sw = bounds.getSouthWest();
      const ne = bounds.getNorthEast();

      const minLat = sw.lat;
      const maxLat = ne.lat;
      const minLong = sw.lng;
      const maxLong = ne.lng;

      const threads = await fetchThreads(minLat, maxLat, minLong, maxLong);

      if (threads === null) {
        unloadAllChatComponents();
        tooManyItems = true;
        return;
      } else {
        tooManyItems = false;
      }

      const visibleThreadIds = new Set<number>();

      threads.forEach((thread) => {
        visibleThreadIds.add(thread.id);
        if (!mountedComponents.has(thread.id)) {
          loadChatComponent(thread, false);
        }
      });

      mountedComponents.forEach((_, id) => {
        if (!visibleThreadIds.has(id)) {
          unloadChatComponent(id);
        }
      });

      loadedThreads = threads;
    } catch (error) {
      console.error("Error fetching threads:", error);
    } finally {
      isLoadingThreads = false;
    }
  }

  function loadChatComponent(thread: Thread, showAnimation: boolean) {
    if (!map || mountedComponents.has(thread.id)) return;

    const componentDom = document.createElement("div");
    componentDom.addEventListener("mouseenter", () => {
      componentDom.style.zIndex = "9999";
    });
    componentDom.addEventListener("mouseleave", () => {
      componentDom.style.zIndex = "0";
    });

    const componentMount = mount(Chat, {
      target: componentDom,
      props: {
        ...thread,
        showAnimation,
        onDelete: () => {
          unloadChatComponent(thread.id);
        },
      },
    });

    const marker = new maplibregl.Marker({
      element: componentDom,
      anchor: "bottom",
    })
      .setLngLat([thread.long, thread.lat])
      .addTo(map);

    mountedComponents.set(thread.id, { marker, mount: componentMount });
  }

  function unloadChatComponent(id: number) {
    const component = mountedComponents.get(id);
    if (!component) return;
    component.marker.remove();
    mountedComponents.delete(id);
  }

  function unloadAllChatComponents() {
    mountedComponents.forEach((_, id) => {
      unloadChatComponent(id);
    });
    loadedThreads = [];
  }

  function zoomTo(long: number, lat: number) {
    map!.flyTo({
      center: [long, lat],
      zoom: 15, // fixed zoom level
      duration: 2000,
      essential: true,
    });
  }

  async function zoomToRandomChat() {
    const thread = await fetchRandomThread();
    zoomTo(thread.long, thread.lat);
  }

  function zoomToMyLocation() {
    if (!map) return;
    if (!navigator.geolocation) {
      alert("Geolocation is not supported by your browser");
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (position) => {
        zoomTo(position.coords.longitude, position.coords.latitude);
      },
      (error) => {
        alert("Error getting location: " + error.message);
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 300000 }
    );
  }

  function zoomIn() {
    if (!map) return;
    const newZoom = Math.min(map.getZoom() + 1, map.getMaxZoom());
    map.flyTo({ zoom: newZoom, duration: 500, essential: true });
  }

  function zoomOut() {
    if (!map) return;
    const newZoom = Math.max(map.getZoom() - 1, map.getMinZoom());
    map.flyTo({ zoom: newZoom, duration: 500, essential: true });
  }

  function onDestroy() {
    unloadAllChatComponents();
    removeAddChatComponent();
    if (map) {
      map.off("zoom", handleMapUpdate);
      map.off("move", handleMapUpdate);
      map.off("moveend", handleMapUpdate);
      map.off("click", handleMapClick);
    }
  }
</script>

<svelte:window on:beforeunload={onDestroy} />

<div id="map" style="width: 100vw; height: 100vh"></div>

{#if tooManyItems}
  <div
    role="alert"
    class="alert fixed top-[10px] right-[50%] translate-x-[50%] rounded-full"
  >
    <span>Too many conversations in view. Zoom in more to see them.</span>
  </div>
{/if}

{#if getAuthenticationStatus() == AuthenticationStatus.LoggedOut}
  <div
    role="alert"
    class="alert fixed top-[10px] right-[50%] translate-x-[50%] rounded-full"
  >
    <DoorOpen />
    <span>Log In to join and create conversations </span>
  </div>
{/if}

<Controls {zoomToRandomChat} {zoomToMyLocation} {zoomIn} {zoomOut}></Controls>
<ThreadTooCloseModal />

<style>
  :global(.chat-marker) {
    z-index: 1;
  }

  /* Hover state for higher z-index */
  :global(.chat-marker:hover) {
    z-index: 9999 !important;
  }
</style>
