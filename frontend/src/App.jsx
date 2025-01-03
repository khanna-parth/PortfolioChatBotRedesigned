import { useEffect, useState } from "react";
import Header from "./Components/Header";
import FileMenuDropDown from "./Components/FileContainer/FileMenu.jsx";
import './index.css';
import ChatContainer from "./Components/ChatContainer/ChatContainer.jsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/Components/ui/dialog"
import { Button } from "./Components/ui/button.jsx";
import { useWebSocketStore } from './store/store.js';
import { WEBSOCKET_ROUTE } from "./routes/routes.js";


function App() {
  const [showPolicy, setShowPolicy] = useState(true);
  const { ws, connection, connID, setWebSocket, setConnected, setConnID } = useWebSocketStore();

  useEffect(() => {
    setConnected(null);
    if (ws) return;
    const timer = setTimeout(() => {
      console.log("Starting WS");
      const socket = new WebSocket(WEBSOCKET_ROUTE);
  
      setWebSocket(socket);
  
      socket.onopen = () => {
        console.log('WebSocket connection established');
        setConnected(true);
  
        const { ws, connected } = useWebSocketStore.getState();
        console.log("After socket open ws: ", ws, " | connState: ", connected);
      };
  
      socket.onmessage = (event) => {
        console.log('Received message:', event.data);
        try {
          const data = JSON.parse(event.data)
          if (data.connID != "") {
            console.log("Recieved connection ID:", data.connID)
            setConnID(data.connID)
          }
        } catch (error) {
          console.log("Parse error:", error)
        }
      };
  
      socket.onclose = () => {
        console.log('WebSocket connection closed');
        setConnected(false);
      };
  
      return () => {
        console.log("Cleaning up websocket");
        socket.close();
        setConnected(false);
      };
  
    }, 3000);
  
    return () => clearTimeout(timer);
  }, [ws]);
  // }, [setConnected, setWebSocket]);

	return (
    <div className="min-h-screen w-full bg-midnight flex flex-col justify-between items-center">
      <div className="flex flex-col-reverse lg:flex-row w-full h-full items-center space-x-10 pr-8 border border-red-400">
        <FileMenuDropDown />
        <ChatContainer />
        <div className="hidden">
          <Dialog open={showPolicy}>
            <DialogContent hideClose>
              <DialogHeader>
                <DialogTitle>Session will not be saved</DialogTitle>
                <DialogDescription className="pb-5">
                  Whatever changes you make will not be persisted and your uploaded files(if any) will be removed at the end of the session. Saved states will be added in the future. For now, upload your own documents or explore Parth's in the preset menu
                </DialogDescription>
                <Button className="bg-black" onClick={() => setShowPolicy(false)}>I understand</Button>
              </DialogHeader>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </div>
  );
}

export default App;
