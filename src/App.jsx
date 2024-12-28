import { useState } from "react";
import Header from "./Components/Header";
import FileMenuDropDown from "./Components/FileMenu";
import './index.css';

function App() {
	const [count, setCount] = useState(0);
  const [chatHistory, setChatHistory] = useState([]);
  const [inputMessage, setInputMessage] = useState("");

  const handleSendChat = (e) => {
    if (e.key === "Enter" && inputMessage.trim() != "") {
      setChatHistory((prevHistory) => [...prevHistory, inputMessage])
    }
  }

	return (
    <div className="min-h-screen w-full bg-midnight flex flex-col justify-between items-center">
      <Header />
      {/* <FileMenuDropDown /> */}
      <div className="bg-[#A2BCE0] w-full sm:max-w-md md:max-w-lg lg:max-w-2xl xl:max-w-6xl h-[80vh] flex-grow flex flex-col justify-start rounded-lg shadow-2xl border border-cyan-200 px-2 mt-14">
        <div>
          <h1 className="text-black text-2xl font-bold flex justify-center">Hey there, start a conversation</h1>
        </div>

        <div className="flex-grow flex flex-col items-center justify-end overflow-y-auto">
          {chatHistory.length == 0 ? (
            <h4 className="text-center font-medium text-2xl">No previous chat history</h4>
          ) : (
            <div className="max-w-full px-4 py-2">
              {chatHistory.map((message, index) => (
                <p key={index} className="text-white my-10 py-8 px-2 bg-blue-800 rounded-xl">{message}</p>
              ))}
            </div>
          )}
        </div>

        <input 
          onChange={(e) => setInputMessage(e.target.value)} onKeyDown={handleSendChat}
          className="relative max-w-full px-10 py-2 rounded-lg bg-white text-gray-800 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 mt-auto mb-4 placeholder-style"
          placeholder="Start chatting..."
          value={inputMessage}
        />
      </div>
    </div>
  );
}

export default App;
