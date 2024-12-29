import React, { useEffect, useState } from 'react'
import { toast } from 'sonner';
import { IoSend } from "react-icons/io5";
import {
    Avatar,
    AvatarFallback,
    AvatarImage,
  } from "@/components/ui/avatar"

const ChatContainer = () => {
    const [count, setCount] = useState(0);
    const [chatHistory, setChatHistory] = useState([]);
    const [inputMessage, setInputMessage] = useState("");
    const [sendPressed, setSendPressed] = useState(false);

    const handleKeyDown = (e) => {
        if (e.key === "Enter" && inputMessage.trim() !== "") {
          handleSendChat();
        }
      };

    const handleSendChat = (e) => {
        if (inputMessage.trim() != "") {
            setChatHistory((prevHistory) => [...prevHistory, inputMessage]);
            setInputMessage("");
            setSendPressed(false);
        }
    }

    useEffect(() => {
        if (sendPressed) {
          handleSendChat();
          setSendPressed(false);
        }
      }, [sendPressed]);

    function pickRandom(element1, element2) {
        return Math.random() < 0.5 ? element1 : element2;
    }

  return (
    <div className="bg-chat-black w-full px-10 h-[90vh] flex flex-col rounded-lg shadow-2xl mt-14">
        <div>
        </div>

        <div className="flex-grow overflow-y-auto px-4 py-2 rounded-lg">
        {
        chatHistory.length == 0 ? (
            <div className='flex flex-col items-center justify-center h-full'>
                <h4 className="text-center items-center justify-center font-medium text-2xl text-white text-opacity-60">Start a conversation below</h4>
            </div>
        ) : (
            <div className="flex flex-col space-y-4">
            {chatHistory.map((message, index) => (
                <div key={index} className={`flex ${index % 2 == 0 ? 'items-end justify-end' : 'items-start justify-start'} overflow-hidden`}>
                    <Avatar>
                        <AvatarImage src={index % 2 == 0 ? 'https://github.com/shadcn.png' : 'https://www.pngfind.com/pngs/m/2-24642_imagenes-random-png-cosas-random-png-transparent-png.png'} alt="toad" />
                        <AvatarFallback>CN</AvatarFallback>
                    </Avatar>
                    <p key={index} className={`text-white my-10 py-4 ${index % 2 == 0 ? 'bg-blue-800' : 'bg-green-800'} px-5 rounded-3xl break-words max-w-[80%]`}>{message}</p>
                </div>
            ))}
            </div>
        )}
        </div>

        <div className='flex flex-row space-x-3'>
            <input 
            onChange={(e) => setInputMessage(e.target.value)} onKeyDown={handleKeyDown}
            className="relative w-full px-10 py-2 rounded-lg bg-white text-gray-800 placeholder-gray-500 focus:outline-none focus:ring-2 overflow-hidden focus:ring-blue-500 mt-auto mb-4 placeholder-style"
            placeholder="Start chatting..."
            value={inputMessage}
            />
            <IoSend size={44} className='relative bg-blue-600 text-white rounded-xl p-1.5 cursor-pointer' onClick={() => setSendPressed(true)} />
        </div>
        
    </div>
  )
}

export default ChatContainer