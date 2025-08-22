<script lang="ts">
  interface CropArea {
    x: number;
    y: number;
    size: number;
  }

  interface Props {
    imageUrl: string;
    onCropComplete: (croppedFile: File) => void;
    onCancel: () => void;
    show: boolean;
  }

  let { imageUrl, onCropComplete, onCancel, show }: Props = $props();

  // Cropping states
  let cropCanvas: HTMLCanvasElement;
  let cropCtx: CanvasRenderingContext2D | null = null;
  let cropImage = new Image();
  let isDragging = $state(false);
  let dragType = $state<"move" | "resize" | "">("");
  let activeHandle = $state<number>(-1); // -1 for move, 0-7 for resize handles
  let cropArea = $state<CropArea>({
    x: 50,
    y: 50,
    size: 200,
  });

  function initializeCropper() {
    if (!cropCanvas || !show) return;

    cropCtx = cropCanvas.getContext("2d");
    if (!cropCtx) return;

    cropImage.onload = () => {
      // Set canvas size
      const maxSize = 400;
      let { width, height } = cropImage;

      if (width > height) {
        if (width > maxSize) {
          height = (height * maxSize) / width;
          width = maxSize;
        }
      } else {
        if (height > maxSize) {
          width = (width * maxSize) / height;
          height = maxSize;
        }
      }

      cropCanvas.width = width;
      cropCanvas.height = height;

      // Initialize crop area in center
      const minSize = Math.min(width, height);
      cropArea.size = Math.min(200, minSize * 0.8);
      cropArea.x = (width - cropArea.size) / 2;
      cropArea.y = (height - cropArea.size) / 2;

      drawCropper();
    };

    cropImage.src = imageUrl;
  }

  function drawCropper() {
    if (!cropCtx || !cropImage.complete) return;

    const { width, height } = cropCanvas;

    // Clear canvas
    cropCtx.clearRect(0, 0, width, height);

    // Draw image
    cropCtx.drawImage(cropImage, 0, 0, width, height);

    // Draw overlay
    cropCtx.fillStyle = "rgba(0, 0, 0, 0.5)";
    cropCtx.fillRect(0, 0, width, height);

    // Clear crop area
    cropCtx.globalCompositeOperation = "destination-out";
    cropCtx.beginPath();
    cropCtx.arc(
      cropArea.x + cropArea.size / 2,
      cropArea.y + cropArea.size / 2,
      cropArea.size / 2,
      0,
      Math.PI * 2
    );
    cropCtx.fill();

    // Draw crop border
    cropCtx.globalCompositeOperation = "source-over";
    cropCtx.strokeStyle = "#3b82f6";
    cropCtx.lineWidth = 2;
    cropCtx.beginPath();
    cropCtx.arc(
      cropArea.x + cropArea.size / 2,
      cropArea.y + cropArea.size / 2,
      cropArea.size / 2,
      0,
      Math.PI * 2
    );
    cropCtx.stroke();

    // Draw multiple handles around the circle
    const handleSize = 12;
    const centerX = cropArea.x + cropArea.size / 2;
    const centerY = cropArea.y + cropArea.size / 2;
    const radius = cropArea.size / 2;

    // Define 8 handle positions around the circle
    const handlePositions = [
      { angle: 0, label: "right" }, // Right
      { angle: Math.PI / 2, label: "bottom" }, // Bottom
      { angle: Math.PI, label: "left" }, // Left
      { angle: (3 * Math.PI) / 2, label: "top" }, // Top
    ];

    // Draw all resize handles
    cropCtx.fillStyle = "#3b82f6";
    handlePositions.forEach((pos, index) => {
      const handleX = centerX + Math.cos(pos.angle) * radius;
      const handleY = centerY + Math.sin(pos.angle) * radius;

      cropCtx!.beginPath();
      cropCtx!.arc(handleX, handleY, handleSize / 2, 0, Math.PI * 2);
      cropCtx!.fill();

      // Add white border to handles for better visibility
      cropCtx!.strokeStyle = "white";
      cropCtx!.lineWidth = 2;
      cropCtx!.stroke();
      cropCtx!.strokeStyle = "#3b82f6";
      cropCtx!.lineWidth = 2;
    });
  }

  function getHandleAtPosition(x: number, y: number): number {
    const centerX = cropArea.x + cropArea.size / 2;
    const centerY = cropArea.y + cropArea.size / 2;
    const radius = cropArea.size / 2;

    // Check center handle for moving
    if (
      Math.sqrt(Math.pow(x - centerX, 2) + Math.pow(y - centerY, 2)) <
      radius - 10
    ) {
      return -1; // Move handle
    }

    // Check resize handles
    const handlePositions = [
      { angle: 0 }, // Right
      { angle: Math.PI / 2 }, // Bottom
      { angle: Math.PI }, // Left
      { angle: (3 * Math.PI) / 2 }, // Top
    ];

    for (let i = 0; i < handlePositions.length; i++) {
      const handleX = centerX + Math.cos(handlePositions[i].angle) * radius;
      const handleY = centerY + Math.sin(handlePositions[i].angle) * radius;

      if (Math.sqrt(Math.pow(x - handleX, 2) + Math.pow(y - handleY, 2)) < 15) {
        return i; // Resize handle index
      }
    }

    return -2; // No handle
  }

  function getCursorForHandle(handleIndex: number): string {
    if (handleIndex === -1) return "move";

    const cursors = [
      "e-resize", // Right
      "s-resize", // Bottom
      "w-resize", // Left
      "n-resize", // Top
    ];

    return cursors[handleIndex] || "default";
  }

  function handleMouseDown(event: MouseEvent) {
    const rect = cropCanvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    activeHandle = getHandleAtPosition(x, y);

    if (activeHandle >= -1) {
      isDragging = true;
      dragType = activeHandle === -1 ? "move" : "resize";
      cropCanvas.style.cursor = getCursorForHandle(activeHandle);
    }
  }

  function handleMouseMove(event: MouseEvent) {
    const rect = cropCanvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    if (!isDragging) {
      // Update cursor based on hover position
      const handleIndex = getHandleAtPosition(x, y);
      cropCanvas.style.cursor =
        handleIndex >= -1 ? getCursorForHandle(handleIndex) : "default";
      return;
    }

    if (dragType === "move") {
      // Move crop area
      cropArea.x = Math.max(
        0,
        Math.min(x - cropArea.size / 2, cropCanvas.width - cropArea.size)
      );
      cropArea.y = Math.max(
        0,
        Math.min(y - cropArea.size / 2, cropCanvas.height - cropArea.size)
      );
    } else if (dragType === "resize") {
      // Resize crop area from any handle
      const centerX = cropArea.x + cropArea.size / 2;
      const centerY = cropArea.y + cropArea.size / 2;

      // Calculate distance from center to mouse position
      const distanceFromCenter = Math.sqrt(
        Math.pow(x - centerX, 2) + Math.pow(y - centerY, 2)
      );

      const newSize = distanceFromCenter * 2;

      // Ensure crop area stays within canvas bounds
      const maxSize = Math.min(
        cropCanvas.width - Math.max(0, centerX - newSize / 2) * 2,
        cropCanvas.height - Math.max(0, centerY - newSize / 2) * 2,
        (centerX + newSize / 2 <= cropCanvas.width
          ? cropCanvas.width
          : centerX) * 2,
        (centerY + newSize / 2 <= cropCanvas.height
          ? cropCanvas.height
          : centerY) * 2
      );

      cropArea.size = Math.min(Math.max(newSize, 50), maxSize);

      // Recalculate position to keep centered
      cropArea.x = Math.max(
        0,
        Math.min(centerX - cropArea.size / 2, cropCanvas.width - cropArea.size)
      );
      cropArea.y = Math.max(
        0,
        Math.min(centerY - cropArea.size / 2, cropCanvas.height - cropArea.size)
      );
    }

    drawCropper();
  }

  function handleMouseUp() {
    isDragging = false;
    dragType = "";
    activeHandle = -2;
    cropCanvas.style.cursor = "default";
  }

  function applyCrop() {
    if (!cropCtx || !cropImage.complete) return;

    // Create a new canvas for the cropped image
    const croppedCanvas = document.createElement("canvas");
    const croppedCtx = croppedCanvas.getContext("2d");
    if (!croppedCtx) return;

    croppedCanvas.width = 200;
    croppedCanvas.height = 200;

    // Calculate scaling factors
    const scaleX = cropImage.naturalWidth / cropCanvas.width;
    const scaleY = cropImage.naturalHeight / cropCanvas.height;

    // Create circular clipping path
    croppedCtx.beginPath();
    croppedCtx.arc(100, 100, 100, 0, Math.PI * 2);
    croppedCtx.clip();

    // Draw the cropped portion
    croppedCtx.drawImage(
      cropImage,
      cropArea.x * scaleX,
      cropArea.y * scaleY,
      cropArea.size * scaleX,
      cropArea.size * scaleY,
      0,
      0,
      200,
      200
    );

    // Convert to blob and create file
    croppedCanvas.toBlob(
      (blob) => {
        if (blob) {
          const croppedFile = new File([blob], "cropped-image.png", {
            type: "image/png",
          });
          onCropComplete(croppedFile);
        }
      },
      "image/png",
      0.9
    );
  }

  $effect(() => {
    if (show && imageUrl) {
      initializeCropper();
    }
  });
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
  >
    <div class="bg-base-100 rounded-2xl p-6 max-w-lg w-full mx-4">
      <div class="flex justify-between items-center mb-4">
        <h3 class="font-bold text-lg">Crop Your Image</h3>
        <button class="btn btn-sm btn-circle btn-ghost" onclick={onCancel}
          >✕</button
        >
      </div>

      <div class="flex justify-center mb-4">
        <canvas
          bind:this={cropCanvas}
          class="border border-base-300 rounded-lg cursor-default"
          onmousedown={handleMouseDown}
          onmousemove={handleMouseMove}
          onmouseup={handleMouseUp}
          onmouseleave={handleMouseUp}
        ></canvas>
      </div>

      <div class="flex justify-end gap-3">
        <button class="btn btn-ghost" onclick={onCancel}>Cancel</button>
        <button class="btn btn-primary" onclick={applyCrop}>Apply Crop</button>
      </div>
    </div>
  </div>
{/if}
