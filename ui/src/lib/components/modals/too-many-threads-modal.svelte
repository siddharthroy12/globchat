<dialog id="too_many_threads_modal" class="modal">
  <div class="modal-box relative rounded-4xl">
    <form method="dialog">
      <button class="btn btn-sm btn-circle btn-soft absolute right-4 top-4"
        >✕</button
      >
    </form>
    <div class="flex flex-col gap-4">
      <div class="flex justify-center items-center">
        <h3 class="font-bold text-xl">Limit Reached</h3>
      </div>
      <p class="text-center">
        You can only create 10 threads on the map. If you want to add new ones
        you'll have to delete old ones.
      </p>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
