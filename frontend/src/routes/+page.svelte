<script lang="ts">
  import Chat from "../components/chat.svelte";
  import maplibregl from "maplibre-gl";
  import { mount, onMount } from "svelte";
  import Controls from "../components/controls.svelte";
  import Addchat from "../components/addchat.svelte";

  let map: null | maplibregl.Map = null;
  let mountedComponents: Map<
    string,
    { marker: maplibregl.Marker; mount: ReturnType<typeof mount> }
  > = new Map();

  // Configuration
  const MIN_ZOOM_LEVEL = 5; // Minimum zoom level to show chat components
  const RANDOM_ZOOM_LEVEL = 16; // Zoom level when focusing on a random chat
  const CHAT_LOCATIONS = [
    { id: "chat1", lng: 13.3775, lat: 52.516 },
    { id: "chat2", lng: 13.4, lat: 52.52 },
    { id: "chat3", lng: 13.42, lat: 52.51 },
    { id: "chat4", lng: 13.35, lat: 52.505 },
    { id: "chat5", lng: 13.39, lat: 52.535 },
    // Add more chat locations as needed
  ];

  onMount(() => {
    map = new maplibregl.Map({
      style: "https://tiles.openfreemap.org/styles/liberty",
      center: [13.4, 52.5],
      zoom: 10,
      container: "map",
      attributionControl: false,
    });

    // Listen for zoom and move events
    map.on("zoom", handleMapUpdate);
    map.on("move", handleMapUpdate);
    map.on("moveend", handleMapUpdate);

    // Initial load
    handleMapUpdate();
  });

  function handleMapUpdate() {
    if (!map) return;

    const currentZoom = map.getZoom();

    if (currentZoom >= MIN_ZOOM_LEVEL) {
      loadVisibleChatComponents();
    } else {
      unloadAllChatComponents();
    }
  }

  function loadVisibleChatComponents() {
    if (!map) return;

    const bounds = map.getBounds();
    const visibleChatIds = new Set<string>();

    CHAT_LOCATIONS.forEach((location) => {
      if (bounds.contains([location.lng, location.lat])) {
        visibleChatIds.add(location.id);

        // Load component if not already loaded
        if (!mountedComponents.has(location.id)) {
          loadChatComponent(location);
        }
      }
    });

    // Unload components that are no longer visible
    mountedComponents.forEach((component, id) => {
      if (!visibleChatIds.has(id)) {
        unloadChatComponent(id);
      }
    });
  }

  function loadChatComponent(location: {
    id: string;
    lng: number;
    lat: number;
  }) {
    if (!map || mountedComponents.has(location.id)) return;

    const componentDom = document.createElement("div");

    // Add hover event listeners for z-index management
    componentDom.addEventListener("mouseenter", () => {
      componentDom.style.zIndex = "9999";
    });

    componentDom.addEventListener("mouseleave", () => {
      componentDom.style.zIndex = "1";
    });

    const componentMount = mount(Chat, {
      target: componentDom,
      props: {},
    });

    const marker = new maplibregl.Marker({
      element: componentDom,
      anchor: "bottom",
    })
      .setLngLat([location.lng, location.lat])
      .addTo(map);

    mountedComponents.set(location.id, { marker, mount: componentMount });

    console.log(`Loaded chat component: ${location.id}`);
  }

  function unloadChatComponent(id: string) {
    const component = mountedComponents.get(id);
    if (!component) return;

    // Remove marker from map
    component.marker.remove();

    // Remove from our tracking
    mountedComponents.delete(id);

    console.log(`Unloaded chat component: ${id}`);
  }

  function unloadAllChatComponents() {
    mountedComponents.forEach((_, id) => {
      unloadChatComponent(id);
    });
  }

  // New function to zoom to a random chat location
  function zoomToRandomChat() {
    if (!map || CHAT_LOCATIONS.length === 0) return;

    const randomIndex = Math.floor(Math.random() * CHAT_LOCATIONS.length);
    const randomChat = CHAT_LOCATIONS[randomIndex];

    map.flyTo({
      center: [randomChat.lng, randomChat.lat],
      zoom: RANDOM_ZOOM_LEVEL,
      duration: 2000, // 2 second animation
      essential: true, // This animation is essential with respect to prefers-reduced-motion
    });
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
    console.log("Getting your location...");

    // Get current position
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { longitude, latitude } = position.coords;

        console.log(`Found location: ${latitude}, ${longitude}`);

        // Zoom to user's location
        map!.flyTo({
          center: [longitude, latitude],
          zoom: RANDOM_ZOOM_LEVEL, // Use same zoom level as random chat
          duration: 2000, // 2 second animation
          essential: true, // This animation is essential with respect to prefers-reduced-motion
        });
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
    if (map) {
      map.off("zoom", handleMapUpdate);
      map.off("move", handleMapUpdate);
      map.off("moveend", handleMapUpdate);
    }
  }
</script>

<svelte:window on:beforeunload={onDestroy} />

<div id="map" style="width: 100vw; height: 100vh"></div>

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
