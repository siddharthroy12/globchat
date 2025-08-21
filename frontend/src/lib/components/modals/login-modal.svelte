<script>
  import logo from "$lib/assets/icon.webp";
  import { googleLogin } from "$lib/google-login";
  import { loginUsingGoogle } from "$lib/services/auth.svelte";
  import GoogleLogo from "../logos/google-logo.svelte";
  function handleGoogleLogin() {
    googleLogin(async (res) => {
      await loginUsingGoogle(res.credential);
      // @ts-ignore
      document.getElementById("login-modal")?.close();
    });
  }
</script>

<dialog id="login-modal" class="modal">
  <div class="modal-box relative">
    <form method="dialog">
      <button class="btn btn-sm btn-circle btn-soft absolute right-2 top-2"
        >✕</button
      >
    </form>
    <div class="flex flex-col gap-4">
      <div class="flex justify-center items-center">
        <img src={logo} width={60} height={60} alt="logo" />
        <h3 class="font-bold text-xl">GlobeChat</h3>
      </div>

      <button
        class="btn bg-white text-black w-full rounded-full"
        onclick={handleGoogleLogin}><GoogleLogo /> Login With Google</button
      >
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
