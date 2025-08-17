<script>
  import logo from "$lib/assets/icon.webp";
  import { X } from "@lucide/svelte";
</script>

<dialog id="login-modal" class="modal">
  <div class="modal-box flex flex-col gap-4 relative">
    <form method="dialog">
      <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
        >✕</button
      >
    </form>

    <div class="flex justify-center items-center">
      <img src={logo} width={60} height={60} alt="logo" />
      <h3 class="font-bold text-xl">GlobeChat</h3>
    </div>
    <button class="btn btn-primary w-full rounded-full"
      >Login With Google</button
    >
  </div>
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
