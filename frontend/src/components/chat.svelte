<script>
  import Avatar from "./avatar.svelte";
</script>

<div class="conversation-bubble-wrapper">
  <div class="conversation-bubble">
    <div class="avatar-wrapper">
      <Avatar />
    </div>
    <div class="content">
      <p class="username">Siddharth Roy <span class="time">1 hr. ago</span></p>
      <p class="last-chat"></p>
      <p>I live here</p>
    </div>
  </div>
</div>

<style>
  .conversation-bubble-wrapper {
    position: relative;
  }
  .conversation-bubble {
    cursor: pointer;
    position: absolute;
    bottom: 0;
    left: 0;
    background: var(--color-base-300);
    padding: 5px;
    color: var(--color-primary-content);
    width: 30px;
    max-height: 30px;
    border-radius: 50px;
    border-bottom-left-radius: 0;
    display: flex;
    overflow: hidden;
    align-items: start;
    gap: 10px;
    transition:
      width 200ms,
      max-height 200ms,
      padding 200ms;
  }
  .avatar-wrapper {
    flex-shrink: 0;
  }
  .conversation-bubble:hover {
    z-index: 999;
    width: 200px;
    max-height: 200px;
    padding: 10px;
    border-radius: 20px;
    border-bottom-left-radius: 0;
  }
  .content {
    display: none;
  }
  .conversation-bubble:hover .content {
    display: block;
  }
  .username {
    font-weight: bold;
    text-wrap: nowrap;
  }
  .last-chat {
    width: 180px;
  }
  .time {
    font-weight: 100;
    color: aliceblue;
  }
</style>
