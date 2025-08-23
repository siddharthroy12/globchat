<script lang="ts">
  import { getUserData } from "$lib/services/auth.svelte";
  import type { Message } from "$lib/services/message.svelte";
  import Avatar from "./avatar.svelte";

  type MessageProps = {
    message: Message;
  };
  const { message }: MessageProps = $props();

  const isFromUser = $derived(getUserData()?.id == message.user_id);

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
  };
</script>

<div class="chat" class:chat-end={isFromUser} class:chat-start={!isFromUser}>
  <div class="chat-image avatar">
    <Avatar src={message.user_image} size={40} />
  </div>
  <div class="chat-header">
    {message.username}
    <time class="text-xs opacity-50">{formatTime(message.created_at)}</time>
  </div>
  <div class="chat-bubble bg-primary text-white">{message.text}</div>
</div>
