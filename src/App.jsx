import { useEffect, useState } from "react";
import Header from "./Components/Header";
import FileMenuDropDown from "./Components/FileContainer/FileMenu.jsx";
import './index.css';
import ChatContainer from "./Components/ChatContainer/ChatContainer.jsx";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "./Components/ui/button.jsx";

function App() {
  const [showPolicy, setShowPolicy] = useState(true);

  // useEffect(() => {
  //     const handleBeforeUnload = (event) => {
  //         event.preventDefault();
  //         event.returnValue = "";
  //     };

  //     window.addEventListener("beforeunload", handleBeforeUnload);

  //     return () => {
  //         window.removeEventListener("beforeunload", handleBeforeUnload);
  //     };
  // }, []);

	return (
    <div className="min-h-screen w-full bg-midnight flex flex-col justify-between items-center">
      <div className="flex flex-col lg:flex-row w-full items-center space-x-10 pr-8">
        <FileMenuDropDown />
        <Dialog open={showPolicy}>
        <DialogContent hideClose>
          <DialogHeader>
            <DialogTitle>Session will not be saved</DialogTitle>
            <DialogDescription>
              Whatever changes you make will not be persisted and your uploaded files(if any) will be removed at the end of the session. Saved states will be added in the future.
            </DialogDescription>
            <Button className="bg-black" onClick={() => setShowPolicy(false)}>I understand</Button>
          </DialogHeader>
        </DialogContent>
      </Dialog>
        <ChatContainer />
      </div>
    </div>
  );
}

export default App;
