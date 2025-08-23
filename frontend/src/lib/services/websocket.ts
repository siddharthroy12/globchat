import type { Message } from "./message.svelte";

export function joinThread(
  threadId: number,
  onNewMessage: (message: Message) => void,
  onDisconnect: () => void
) {
  let origin = window.origin;
  origin = origin.replace("https://", "wss://");
  origin = origin.replace("http://", "ws://");

  const socket = new WebSocket(`${origin}/api/v1/ws`);

  // Connection opened
  socket.addEventListener("open", (event) => {
    // Join room
    socket.send(
      JSON.stringify({
        room_id: threadId,
        type: "join",
      })
    );
  });

  // Listen for messages
  socket.addEventListener("message", (event) => {
    const json = JSON.parse(event.data);
    if (json["type"] == "new-message") {
      onNewMessage(json["message"]);
    }
  });

  socket.onclose = onDisconnect;

  return () => {
    socket.close();
  };
}
