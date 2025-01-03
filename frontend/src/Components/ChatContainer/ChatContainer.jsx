import React, { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner';
import { IoSend } from "react-icons/io5";
import { Avatar, AvatarFallback, AvatarImage } from "@/Components/ui/avatar"
import { apiClient } from '../../lib/client.js';
import { PRESET_PROMPT_ROUTE, PROMPT_ROUTE } from '../../routes/routes.js';
import { FaCircle } from "react-icons/fa";
import docbotImage from '../../assets/docubot.png'
import { useStore, useWebSocketStore } from '../../store/store.js';
import { dotWave, quantum } from 'ldrs';

const ChatContainer = () => {
    const [inputMessage, setInputMessage] = useState("");
    const [sendPressed, setSendPressed] = useState(false);
    const [isActive, setIsActive] = useState(null);
    const { chatHistory, setChatHistory } = useStore();
    const { ws, connected, connID } = useWebSocketStore();
    const [generating, setGenerating] = useState(false);
    const chatContainerRef = useRef(null);

    dotWave.register();
    quantum.register();

    const handleKeyDown = (e) => {
        if (!connected) return;
        if (e.key === "Enter" && inputMessage.trim() !== "") {
          handleSendChat();
        }
      };

    const handleSendChat = (e) => {
        if (!connected) return;
        if (inputMessage.trim() != "") {
            setChatHistory({sender: "client", content: inputMessage});
            sendPrompt(inputMessage);
            setInputMessage("");
            setSendPressed(false);
        }
    }

    const sendPrompt = async (prompt) => {
        try {
            const data = {
                userID: "pkhanna",
                prompt: prompt
            };

            const headers = {
                "X-Connection-ID": connID,
            };

            setGenerating(true);
            let response;

            if (connID.startsWith("@")) {
                response = await apiClient.post(PRESET_PROMPT_ROUTE, data, { headers }) 
            } else {
                response = await apiClient.post(PROMPT_ROUTE, data, { headers }) 
            }

            if (response.status === 200) {
                console.log(response.data);
                setChatHistory({sender: "server", content: response.data.response})
                setGenerating(false);
            }
        } catch (error) {
            console.log(error)
            const errorMsg = `Error creating response: ${error.response.data}`
            toast.error(errorMsg)
            setGenerating(false);
        }
    }

    useEffect(() => {
        if (connected === true) {
          toast.info("Connection established");
          setIsActive(true);
        } else if (connected === false) {
          toast.info("Connection invalid");
          setIsActive(false);
        } else {
            setIsActive(null); 
        }
    }, [connected])

    useEffect(() => {
        if (chatContainerRef.current) {
            chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
        }
    }, [chatHistory]);


    useEffect(() => {
        if (sendPressed) {
          handleSendChat();
          setSendPressed(false);
        }
      }, [sendPressed]);

  return (
    <div className="bg-chat-black w-full px-7 h-[90vh] flex flex-col rounded-lg shadow-2xl mt-14 mb-10">
        <div className='flex w-full items-center justify-center mt-3 space-x-2'>
            {
                isActive !== null && (
                    <div className='flex flex-row space-x-4'>
                        <Avatar>
                            <AvatarImage src={docbotImage} alt="docubot image" className="rounded-full ring-1 ring-blue-500 scale-90" />
                            <AvatarFallback>CN</AvatarFallback>
                        </Avatar> 
                        <div className='flex flex-col'>
                            <h1 className='text-white font-bold text-lg'>DocuBot</h1>
                            <div className='flex flex-row space-x-2 items-center justify-center'>
                                <FaCircle className={`${isActive ? 'text-green-600' : 'text-red-600'} text-sm scale-75`} />
                                <h1 className={`${isActive ? 'text-green-600' : 'text-red-600'} text-xs`} >{isActive ? "Online" : "Offline"}</h1>
                            </div>
                        </div>
                    </div>
                )}
            
        </div>

        <div className="flex-grow overflow-y-auto px-4 py-2 rounded-lg" ref={chatContainerRef}>
            <>
                {
                    isActive === null || isActive === false ? (
                        <div className='flex flex-col items-center justify-center h-full'>
                            <h1 className='text-3xl font-bold text-white/50'>Connecting</h1>
                            <l-quantum size={50} speed={1} color={'white'} className="m-5"></l-quantum>
                        </div>
                    ) : (
                            chatHistory.length == 0 ? (
                                <div className='flex flex-col items-center justify-center h-full'>
                                    <h4 className="text-center items-center justify-center font-medium text-2xl text-white text-opacity-60">Start a conversation below</h4>
                                </div>
                            ) : (
                            <div className="flex flex-col space-y-4">
                                {chatHistory.map((message, index) => (
                                    <div key={index} className='flex w-full'>
                                        {
                                            message.sender && message.content ? (
                                                <div className={`flex w-full ${message.sender == "client" ? 'items-end justify-end' : 'items-start justify-start'} overflow-hidden`}>
                                                    <Avatar>
                                                        <AvatarImage src={message.sender == "client" ? 'https://github.com/shadcn.png' : docbotImage} alt="avatar" />
                                                        <AvatarFallback>CN</AvatarFallback>
                                                    </Avatar>
                                                    <p className={`text-white my-10 py-4 ${message.sender == "client" ? 'bg-blue-800' : 'bg-green-800'} px-5 rounded-3xl break-words max-w-[80%]`}>{message.content}</p>
                                                </div>
                                            ) : null
                                        }
                                        {
                                            message.docChange ? (
                                                <div className='flex w-full items-center justify-center'>
                                                    <h6 className='text-sm text-white/50'>{message.docChange}</h6>
                                                </div>
                                            ) : null
                                        }
                                    </div>
                                ))}
                            </div>
                        )
                    )
                }
                {
                    generating && (
                        <div className='w-full relative left-0 bottom-0'>
                            <l-dot-wave size={30} speed={1} color={'white'} className="mb-5 pb-5 ml-5"></l-dot-wave>
                        </div>
                )}
            </>
        </div>

        <div className='flex flex-row w-full min-w-full space-x-3 items-center justify-start align-middle mb-5'>
            <span className='flex relative left-0 text-gray-50 items-start justify-center'>{inputMessage.length}/150</span>
            <input
            onChange={(e) => {
                if (e.target.value.length <= 150) {
                    setInputMessage(e.target.value)
                }
            }}
            onKeyDown={handleKeyDown}
            className="relative w-full px-10 py-2 rounded-lg bg-white text-gray-800 placeholder-gray-500 focus:outline-none focus:ring-2 overflow-hidden focus:ring-blue-500 placeholder-style"
            placeholder="Start chatting..."
            value={inputMessage}
            maxLength={150}
            />
            <IoSend size={44} className='relative bg-blue-600 text-white rounded-xl p-1.5 cursor-pointer' onClick={() => setSendPressed(true)} />
        </div>
        
    </div>
  )
}

export default ChatContainer