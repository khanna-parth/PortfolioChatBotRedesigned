import React, {useEffect, useRef, useState} from 'react'
import { HiDocumentPlus } from "react-icons/hi2";
import axios from 'axios';
import FileListDisplay from './components/FileListDisplay.jsx';
import { toast } from 'sonner';
import { BiSolidErrorAlt } from "react-icons/bi";
import { PRESETS_ROUTE, UPLOAD_DOC_ROUTE } from '../../routes/routes.js';
import { useStore, useWebSocketStore } from '../../store/store.js';
import { hourglass } from 'ldrs';
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Button } from "@/Components/ui/button"
import { apiClient } from '../../lib/client.js';


export default function FileMenu() {
    const [filePreview, setFilePreview] = useState(null);
    const [fileName, setFileName] = useState("");
    const fileInputRef = useRef(null);
    const { allowUploads, isProcessing, setIsProcessing, setChatHistory, resetChatHistory } = useStore();
    const { connID, connected, docCount, setDocCount, setConnID, setWebSocket } = useWebSocketStore();

    const [selectedPreset, setSelectedPreset] = useState("Custom");
    const [presets, setPresets] = useState([]);

    hourglass.register()

    const triggerFileInput = () => {
      console.log("Triggering file input ref");
      fileInputRef.current.click();
    }

    const getPresets = async () => {
      try {
        console.log("Trying to get presets")
        const response = await apiClient.get(PRESETS_ROUTE)
        if (response.status === 200 && response.data != null) {
          console.log(response);
          setPresets(response.data);
          console.log(`Set presets to ${response.data}`);
        }
      } catch (error) {
        console.log(error);
      }
    }

    const changePreset = (selection) => {
      if (selection == selectedPreset) return;
      resetChatHistory();
      if (selection == "Custom") {
        setWebSocket(null);
        setSelectedPreset("Custom")
      } else {
        console.log(`Preset changed to ${selection}`);
        setSelectedPreset(selection)
        setConnID(selection);
      }
    }

    const isIDPreset = (id) => {
      if (id.startsWith("@")) {
        return true;
      }
      return false;
    }

    useEffect(() => {
      if (connected) {
        getPresets()
      }
    }, [connected])

    const handleFileUpload = (event) => {
        console.log("Triggered file upload");
    
        const file = event.target.files[0];
        if (!file) {
          console.log("No file selected");
          return;
        }

        setFileName(file.name);
        console.log(`File size: ${file.size} and of type: ${file.type}`)
    
        const formData = new FormData();
        formData.append("myfile", file);
        setIsProcessing(true);
    
        axios
          .post(UPLOAD_DOC_ROUTE, formData, {
            headers: {
              "Content-Type": "multipart/form-data",
              "X-Connection-ID": connID
            },
          })
          .then((response) => {
            setIsProcessing(false);
            console.log("File uploaded:", response);
            setChatHistory({docChange: `You added ${file.name} into the chat context`})
            setDocCount(docCount+1);
            
          })
          .catch((error) => {
            setIsProcessing(false);
            const errorMsg = `Error uploading document: ${error.response.data}`
            toast.error(errorMsg)
            console.error("Error uploading file:", error);
          });
      };

    return (
        <div className='flex flex-col w-[20vw] h-[87vh] items-start justify-between ml-10'>
          <div className='flex flex-col w-full relative top-0 justify-center items-center'>
            {
              presets.length > 0 && (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button className="text-white bg-[#141e37] shadow-xl p-5 border border-white/50" variant="primary">Presets</Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent className="w-56">
                    <DropdownMenuLabel>Options</DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    {
                      presets.map ((preset, index) => (
                        <DropdownMenuCheckboxItem key={index} checked={selectedPreset == preset} onCheckedChange={() => changePreset(preset)}>{preset}</DropdownMenuCheckboxItem>
                      ))
                    }
                    <DropdownMenuCheckboxItem checked={selectedPreset == "Custom"} onCheckedChange={() => changePreset("Custom")}>Custom</DropdownMenuCheckboxItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              )
            }
            
          </div>
          <div className='flex flex-col w-full items-center justify-start'>
            <div className='items-center justify-center p-5 bg-[#272e3f] rounded-3xl space-y-2'>
              <div className='relative items-center justify-center'>
                {
                  <>
                      { allowUploads == null ? (
                          <div className='flex flex-col items-center justify-center'>
                          </div>
                      ) : allowUploads == true ?
                      (
                        <>
                          {
                            isIDPreset(connID) && (
                              <div className='w-full flex flex-col justify-center items-center space-y-4 p-3'>
                                <BiSolidErrorAlt size={60} color='white'/>
                                <h4 className='w-full text-lg flex text-white'>Presets Active</h4>
                              </div>
                            )
                          }
                          {
                            !isIDPreset(connID) && (
                              isProcessing ? (
                                <div className='w-full flex flex-col items-center justify-center space-y-4'>
                                  <h1 className='text-white text-xl font-medium'>Processing</h1>
                                  <l-hourglass size={40} speed={'1.75'} color={'white'}></l-hourglass>
                                </div>
                              ) : (
                                <div className='flex flex-col space-y-4'>
                                  <HiDocumentPlus size={60} color='white' onClick={triggerFileInput} onChange={triggerFileInput} className='w-full items-center justify-center cursor-pointer' />
                                  <input className="hidden" ref={fileInputRef} type="file" onChange={handleFileUpload} accept='.pdf, .txt, .doc, .docx, .word' />
                                  <h1 className='text-white text-md'>Upload a file here</h1>
                                </div>
                              )
                            )
                          }
                        </>
                      ) : (
                        <div className='flex flex-col items-center justify-center'>
                            <BiSolidErrorAlt size={60} className='text-red-600'/>
                            <h1 className='text-red-600'>Server Refused</h1>
                            <h1 className='text-white'>Capacity limit reached on server</h1>
                        </div>
                      )
                      }
                  </>
                }
              </div>
            </div>

            <div className='mt-10 pt-10'>
              <FileListDisplay />
            </div>
          </div>
          <div>

          </div>
        </div>
    );
}
