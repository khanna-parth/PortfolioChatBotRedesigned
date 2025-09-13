import { useEffect, useState } from "react";
import Header from "./components/Header";
import FileMenuDropDown from "./components/FileContainer/FileMenu.jsx";
import './index.css';
import ChatContainer from "./components/ChatContainer/ChatContainer.jsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "./components/ui/button.jsx";
import { useWebSocketStore, useStore } from './store/store.js';
import { SUGGESTIONS_ROUTE, WEBSOCKET_ROUTE } from "./routes/routes.js";
import { apiClient } from "./lib/client.js";

function App() {
  const [showPolicy, setShowPolicy] = useState(true);
  const { ws, connection, connID, setWebSocket, setConnected, setConnID } = useWebSocketStore();
  const { showSuggestions, setShowSuggestions } = useStore();
  const [suggestionsReady, setSuggestionsReady] = useState(false);
  const [suggestions, setSuggestions] = useState([]);

  const getSuggestions = async () => {
    try {
      const response = await apiClient.get(SUGGESTIONS_ROUTE)
      console.log(response);
      setSuggestions(response.data.suggestions);
      setTimeout(3000);
      setSuggestionsReady(true);
      console.log(`Set suggestions to ${response.data.suggestions}`)
    } catch (error) {
      setSuggestions(["Something went wrong with getting you suggestions."]);
      setSuggestionsReady(true);
      console.log(`Suggestions error: ${error}`)
    }
  }

  useEffect(() => {
    if (!showSuggestions) return;
    getSuggestions();
  }, [showSuggestions])

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
          <Dialog open={suggestionsReady}>
            <DialogContent hideClose>
              <DialogHeader>
                <DialogTitle>What you may be curious about...</DialogTitle>
              </DialogHeader>
              <div className="flex flex-col">
                <ul>
                  {
                    suggestions.length > 0 && (
                      suggestions.map((suggestion, index) => {
                        <li key={index}>{suggestion}21212121</li>
                      })
                    )
                  }
                  {
                    suggestions.length < 0 && (
                      <div>No content</div>
                    )
                  }
                </ul>
              </div>
              <Button className="bg-black" onClick={() => setSuggestionsReady(false)}>Got it</Button>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </div>
  );
}

export default App;
