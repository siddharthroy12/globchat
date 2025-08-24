import type { Message } from "./message.svelte";
import type { Thread } from "./threads.svelte";

type JoinThreadInputs = {
  threadId: number;
  onNewMessage: (message: Message) => void;
  onDeleteMessage: (message: Message) => void;
  onDeleteThread: (thread: Thread) => void;
  onDisconnect: () => void;
};

export function joinThread(inputs: JoinThreadInputs) {
  let origin = window.origin;
  origin = origin.replace("https://", "wss://");
  origin = origin.replace("http://", "ws://");

  const socket = new WebSocket(`${origin}/api/v1/ws`);

  // Connection opened
  socket.addEventListener("open", (event) => {
    // Join room
    socket.send(
      JSON.stringify({
        room_id: inputs.threadId,
        type: "join",
      })
    );
  });

  // Listen for messages
  socket.addEventListener("message", (event) => {
    const json = JSON.parse(event.data);
    switch (json["type"]) {
      case "new-message":
        if (json["data"].thread_id === inputs.threadId)
          inputs.onNewMessage(json["data"]);
        break;
      case "delete-message":
        if (json["data"].thread_id === inputs.threadId)
          inputs.onDeleteMessage(json["data"]);
        break;
      case "delete-thread":
        if (json["data"].id === inputs.threadId)
          inputs.onDeleteThread(json["data"]);
        break;
    }
  });

  socket.onclose = inputs.onDisconnect;

  return () => {
    socket.close();
  };
}
