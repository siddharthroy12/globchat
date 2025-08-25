<dialog id="thread_too_close_modal" class="modal">
  <div class="modal-box relative rounded-4xl">
    <form method="dialog">
      <button class="btn btn-sm btn-circle btn-soft absolute right-4 top-4"
        >✕</button
      >
    </form>
    <div class="flex flex-col gap-4">
      <div class="flex justify-center items-center">
        <h3 class="font-bold text-xl">Too Close!</h3>
      </div>
      <p class="text-center">
        You are trying to create a chat too close to others. Please choose a
        different location nearby.
      </p>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
