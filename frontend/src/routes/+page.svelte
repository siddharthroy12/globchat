<script lang="ts">
  import Chat from "../lib/components/chat.svelte";
  import maplibregl from "maplibre-gl";
  import { mount, onMount } from "svelte";
  import Controls from "../lib/components/controls.svelte";
  import Addchat from "../lib/components/addchat.svelte";
  import type { Thread } from "$lib/services/threads.svelte";
  import {
    fetchRandomThread,
    fetchThreads,
  } from "$lib/services/threads.svelte";
  import { DoorOpen, ZoomIn } from "@lucide/svelte";
  import {
    AuthenticationStatus,
    getAuthenticationStatus,
    getUserData,
  } from "$lib/services/auth.svelte";

  let map: null | maplibregl.Map = null;
  let mountedComponents: Map<
    number,
    { marker: maplibregl.Marker; mount: ReturnType<typeof mount> }
  > = new Map();

  // Track AddChat component
  let addChatComponent: {
    marker: maplibregl.Marker;
    mount: ReturnType<typeof mount>;
  } | null = null;

  // Configuration
  const MIN_ZOOM_LEVEL = 10; // Minimum zoom level to show chat components

  let mapZoomLevel = $state(0);
  let zoomedInEnough = $derived(mapZoomLevel >= MIN_ZOOM_LEVEL);

  // Track loaded threads to avoid duplicates
  let loadedThreads: Thread[] = [];
  let isLoadingThreads = false;

  onMount(() => {
    map = new maplibregl.Map({
      style: "https://tiles.openfreemap.org/styles/liberty",
      center: [0, 0],
      zoom: 0,
      container: "map",
      attributionControl: false,
      doubleClickZoom: false,
    });

    // Listen for zoom and move events
    map.on("zoom", handleMapUpdate);
    map.on("move", handleMapUpdate);
    map.on("moveend", handleMapUpdate);

    // Listen for click events on the map
    map.on("click", handleMapClick);

    // Initial load
    handleMapUpdate();
  });

  function handleMapClick(e: maplibregl.MapMouseEvent) {
    if (!map) return;

    if (!zoomedInEnough) {
      return;
    }

    const { lng, lat } = e.lngLat;

    // Remove existing AddChat component if it exists
    removeAddChatComponent();

    // Mount new AddChat component at clicked location
    mountAddChatComponent(lng, lat);
  }

  function mountAddChatComponent(long: number, lat: number) {
    if (!map || addChatComponent) return;

    const componentDom = document.createElement("div");

    // Add hover event listeners for z-index management
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
          // Add to loaded threads
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

    // Remove marker from map
    addChatComponent.marker.remove();

    // Clear the reference
    addChatComponent = null;
  }

  async function handleMapUpdate() {
    if (!map) return;

    mapZoomLevel = map.getZoom();

    if (zoomedInEnough) {
      await loadVisibleThreads();
    } else {
      unloadAllChatComponents();
    }
  }

  async function loadVisibleThreads() {
    if (!map || isLoadingThreads) return;

    isLoadingThreads = true;

    try {
      const center = map.getCenter();
      const zoom = map.getZoom();
      const kmRadius = 50;

      // Fetch threads from API
      const threads = await fetchThreads(center.lat, center.lng, kmRadius);

      // Get currently visible thread IDs
      const visibleThreadIds = new Set<number>();

      // Process each thread
      threads.forEach((thread) => {
        const threadId = thread.id;
        visibleThreadIds.add(threadId);

        // Load component if not already loaded
        if (!mountedComponents.has(threadId)) {
          loadChatComponent(thread, false);
        }
      });

      // Unload components that are no longer in the fetched results
      mountedComponents.forEach((component, id) => {
        if (!visibleThreadIds.has(id)) {
          unloadChatComponent(id);
        }
      });

      // Update loaded threads
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

    // Add hover event listeners for z-index management
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

    // Remove marker from map
    component.marker.remove();

    // Remove from our tracking
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
      zoom: MIN_ZOOM_LEVEL + 5,
      duration: 2000, // 2 second animation
      essential: true, // This animation is essential with respect to prefers-reduced-motion
    });
  }

  // New function to zoom to a random chat location from loaded threads
  async function zoomToRandomChat() {
    const thread = await fetchRandomThread();
    zoomTo(thread.long, thread.lat);
  }

  function zoomToMyLocation() {
    if (!map) return;

    // Check if geolocation is supported
    if (!navigator.geolocation) {
      console.error("Geolocation is not supported by this browser");
      // Optionally show a user-friendly message
      alert("Geolocation is not supported by your browser");
      return;
    }

    // Show loading state (optional - you might want to update UI)

    // Get current position
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { longitude, latitude } = position.coords;

        zoomTo(longitude, latitude);
      },
      (error) => {
        console.error("Error getting location:", error);

        // Handle different error types
        let errorMessage = "Unable to get your location";
        switch (error.code) {
          case error.PERMISSION_DENIED:
            errorMessage = "Location access denied by user";
            break;
          case error.POSITION_UNAVAILABLE:
            errorMessage = "Location information unavailable";
            break;
          case error.TIMEOUT:
            errorMessage = "Location request timed out";
            break;
        }

        console.error(errorMessage);
        // Optionally show user-friendly error message
        alert(errorMessage);
      },
      {
        enableHighAccuracy: true, // Use GPS if available
        timeout: 10000, // 10 second timeout
        maximumAge: 300000, // Accept cached position up to 5 minutes old
      }
    );
  }

  function zoomIn() {
    if (!map) return;

    const currentZoom = map.getZoom();
    const newZoom = Math.min(currentZoom + 1, map.getMaxZoom()); // Increment by 1, respect max zoom

    map.flyTo({
      zoom: newZoom,
      duration: 500, // Smooth 0.5 second animation
      essential: true,
    });
  }

  function zoomOut() {
    if (!map) return;

    const currentZoom = map.getZoom();
    const newZoom = Math.max(currentZoom - 1, map.getMinZoom()); // Decrement by 1, respect min zoom

    map.flyTo({
      zoom: newZoom,
      duration: 500, // Smooth 0.5 second animation
      essential: true,
    });
  }

  // Cleanup on component destroy
  function onDestroy() {
    unloadAllChatComponents();
    removeAddChatComponent(); // Clean up AddChat component
    if (map) {
      map.off("zoom", handleMapUpdate);
      map.off("move", handleMapUpdate);
      map.off("moveend", handleMapUpdate);
      map.off("click", handleMapClick); // Remove click listener
    }
  }
</script>

<svelte:window on:beforeunload={onDestroy} />

<div id="map" style="width: 100vw; height: 100vh"></div>

{#if !zoomedInEnough}
  <div
    role="alert"
    class="alert fixed top-[10px] right-[50%] translate-x-[50%] rounded-full"
  >
    <ZoomIn />
    <span>Zoom in more to view and create conversations </span>
  </div>
{/if}

{#if zoomedInEnough && getAuthenticationStatus() == AuthenticationStatus.LoggedOut}
  <div
    role="alert"
    class="alert fixed top-[10px] right-[50%] translate-x-[50%] rounded-full"
  >
    <DoorOpen />
    <span>Log In to join and create conversations </span>
  </div>
{/if}

<Controls {zoomToRandomChat} {zoomToMyLocation} {zoomIn} {zoomOut}></Controls>

<style>
  :global(.chat-marker) {
    z-index: 1;
  }

  /* Hover state for higher z-index */
  :global(.chat-marker:hover) {
    z-index: 9999 !important;
  }
</style>
