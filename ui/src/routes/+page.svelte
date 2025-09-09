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
  import { setMap } from "$lib/services/map.svelte";

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

  // Debouncing variables
  const DEBOUNCE_DURATION = 100; // milliseconds - you can adjust this
  let debounceTimeout: number | null = null;

  onMount(() => {
    map = new maplibregl.Map({
      style: "https://tiles.openfreemap.org/styles/liberty",
      center: [0, 0],
      zoom: 0,
      container: "map",
      attributionControl: false,
      doubleClickZoom: false,
      renderWorldCopies: false,
    });

    setMap(map);

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
      map?.on("load", () => {
        zoomTo(thread.long, thread.lat);
        setTimeout(() => {
          loadChatComponent(thread, false, true);
        }, 2000);
      });
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

  function getBox() {
    if (map) {
      const canvas = map.getCanvas();
      const width = canvas.width;
      const height = canvas.height;

      const corners = [
        [0, 0], // top-left
        [width, 0], // top-right
        [width, height], // bottom-right
        [0, height], // bottom-left
      ];

      // Convert screen coordinates to LngLat
      //@ts-ignore
      const lngLats = corners.map((pt) => map!.unproject(pt));

      // Optionally, compute the bounding box
      let west = lngLats[0].lng,
        south = lngLats[0].lat;
      let east = lngLats[0].lng,
        north = lngLats[0].lat;

      for (const coord of lngLats) {
        if (coord.lng < west) west = coord.lng;
        if (coord.lng > east) east = coord.lng;
        if (coord.lat < south) south = coord.lat;
        if (coord.lat > north) north = coord.lat;
      }

      const accurateBoundsArray = [
        [west, south], // SW
        [east, north], // NE
      ];

      return accurateBoundsArray;
    }
    return null;
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

    // Clear existing debounce timeout
    if (debounceTimeout !== null) {
      clearTimeout(debounceTimeout);
    }

    // Set new debounced call
    debounceTimeout = setTimeout(async () => {
      await loadVisibleThreads();
      debounceTimeout = null;
    }, DEBOUNCE_DURATION);
  }

  async function loadVisibleThreads() {
    if (!map || isLoadingThreads) return;
    isLoadingThreads = true;

    try {
      // Use getBox() instead of map.getBounds()
      const boundsArray = getBox();
      if (!boundsArray) return;

      const [[west, south], [east, north]] = boundsArray;

      const minLat = south;
      const maxLat = north;
      const minLong = west;
      const maxLong = east;

      let threads: Thread[] | null = null;

      // Check if we're crossing the International Date Line
      if (minLong > maxLong) {
        console.log("Crossing International Date Line", { minLong, maxLong });

        // Split the query into two parts:
        // 1. From minLong to 180 (eastern side)
        // 2. From -180 to maxLong (western side)
        const threadsEast = await fetchThreads(minLat, maxLat, minLong, 180);
        const threadsWest = await fetchThreads(minLat, maxLat, -180, maxLong);

        // Handle the case where either query returns null (too many items)
        if (threadsEast === null || threadsWest === null) {
          threads = null;
        } else {
          // Combine results from both sides
          threads = [...threadsEast, ...threadsWest];

          // Remove duplicates if any (threads exactly on the date line might appear twice)
          const uniqueThreads = new Map<number, Thread>();
          threads.forEach((thread) => {
            uniqueThreads.set(thread.id, thread);
          });
          threads = Array.from(uniqueThreads.values());
        }
      } else {
        // Normal case - no date line crossing
        threads = await fetchThreads(minLat, maxLat, minLong, maxLong);
      }

      // Handle too many items case
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

      // Remove components that are no longer visible
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

  function loadChatComponent(
    thread: Thread,
    showAnimation: boolean,
    startOpen: boolean = false
  ) {
    if (!map || mountedComponents.has(thread.id)) {
      if (startOpen) {
        unloadChatComponent(thread.id);
      } else {
        return;
      }
    }

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
        startOpen,
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
    // Clear any pending debounce timeout
    if (debounceTimeout !== null) {
      clearTimeout(debounceTimeout);
      debounceTimeout = null;
    }

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
