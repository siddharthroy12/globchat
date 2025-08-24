<script>
  import logo from "$lib/assets/icon.webp";
  import { Gavel, Info } from "@lucide/svelte";
</script>

<dialog id="about_modal" class="modal">
  <div
    class="modal-box relative sm:max-h-11/12 max-sm:size-full max-sm:px-4 sm:max-w-5xl rounded-4xl"
  >
    <form method="dialog">
      <button class="btn btn-sm btn-circle btn-soft absolute right-4 top-4"
        >✕</button
      >
    </form>
    <div class="flex flex-col gap-4">
      <div class="flex justify-center items-center">
        <img src={logo} width={90} height={90} alt="logo" />
        <h3 class="font-bold text-2xl">GlobeChat</h3>
      </div>
    </div>
    <p class="text-center text-sm">
      Map powered by: <a
        href="https://maplibre.org/"
        target="_blank"
        class="link">MapLibre</a
      >
      |
      <a href="https://openfreemap.org/" target="_blank" class="link"
        >OpenFreeMap</a
      >
    </p>
    <div class="flex flex-col gap-6">
      <section>
        <h3 class="text-lg font-semibold flex items-center gap-3 mb-2">
          <Info />Overview
        </h3>
        <p>
          GlobeChat is a platform where you can start conversation anywhere on
          the world map and chat with people who are interested in the same
          location as you. Start topics anywhere and let other people find it
          and join the conversation.
        </p>
      </section>
      <section>
        <h3 class="text-lg font-semibold flex items-center gap-3 mb-2">
          <Gavel />Rules
        </h3>
        <div class="flex gap-2 flex-col">
          <div role="alert" class="alert alert-error alert-soft">
            <span
              >🚨 Do not post illegal content. Your email will be included in a
              special list if you do it.</span
            >
          </div>
          <div role="alert" class="alert alert-error alert-soft">
            <span>🤖 Use of bots is not allowed.</span>
          </div>
          <div role="alert" class="alert alert-error alert-soft">
            <span>🤬 Hate speech is not allowed.</span>
          </div>
          <div role="alert" class="alert alert-error alert-soft">
            <span
              >👥 Disclosing other's personal information inluding location is
              not allowed.</span
            >
          </div>
          <div role="alert" class="alert alert-warning alert-soft">
            <span
              >✋ Disobeying the rules will lead to temporary or permanent ban</span
            >
          </div>
          <div role="alert" class="alert alert-success alert-soft">
            <span>😊 Keep conversations civil and friendly</span>
          </div>
          <div role="alert" class="alert alert-success alert-soft">
            <span>‼️ Report messages that do not follow the rules</span>
          </div>
        </div>
      </section>
      <section>
        <h3 class="text-lg font-semibold flex items-center gap-3 mb-2">
          <Gavel /> FAQs
        </h3>
        <div
          class="join join-vertical bg-base-100 w-full bg-base-200 rounded-md"
        >
          <div
            class="collapse collapse-arrow join-item border-base-300 border rounded-md"
          >
            <input type="radio" name="my-accordion-4" />
            <div class="collapse-title font-semibold">Is my info public?</div>
            <div class="collapse-content text-sm">
              <p>
                Only your Name, Profile Picture and DB ID is public. Your email
                will not be visible to others.
              </p>
              <p>
                I recommend to not use your own name and image for your profile.
              </p>
            </div>
          </div>
          <div class="collapse collapse-arrow join-item border-base-300 border">
            <input type="radio" name="my-accordion-4" />
            <div class="collapse-title font-semibold">
              How do I create account?
            </div>
            <div class="collapse-content text-sm">
              Click on Login button after closing this window. Registrations and
              Login both are done using Google OAuth. There are no other ways.
            </div>
          </div>
          <div class="collapse collapse-arrow join-item border-base-300 border">
            <input type="radio" name="my-accordion-4" />
            <div class="collapse-title font-semibold">Is this site safe?</div>
            <div class="collapse-content text-sm">
              As long as you don't reveal your true indentity. You are safe.
            </div>
          </div>
          <div class="collapse collapse-arrow join-item border-base-300 border">
            <input type="radio" name="my-accordion-4" />
            <div class="collapse-title font-semibold">
              Is this a copy of wplace?
            </div>
            <div class="collapse-content text-sm">
              This website indeed is inspired by wplace. I really liked the
              concept of communicating with the world using pixel art. And I
              thought why not have conversations on the world map too.
            </div>
          </div>
          <div class="collapse collapse-arrow join-item border-base-300 border">
            <input type="radio" name="my-accordion-4" />
            <div class="collapse-title font-semibold">
              How do I report someone?
            </div>
            <div class="collapse-content text-sm">
              Hover over the message that you want to report. You'll see an
              option to report the message. The moderators will review and take
              necessary actions.
            </div>
          </div>
          <div class="collapse collapse-arrow join-item border-base-300 border">
            <input type="radio" name="my-accordion-4" />
            <div class="collapse-title font-semibold">Is this Open Source?</div>
            <div class="collapse-content text-sm">
              Not at the moment. It's no certain if I'm going to release the
              source code in future.
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
