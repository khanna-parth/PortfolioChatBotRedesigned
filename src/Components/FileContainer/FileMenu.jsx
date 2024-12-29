import React, {useEffect, useRef, useState} from 'react'
import { HiDocumentPlus } from "react-icons/hi2";
import axios from 'axios';
import FileListDisplay from './components/FileListDisplay.jsx';
import { Button } from '../ui/button.jsx';
import { toast } from 'sonner';

export default function FileMenu() {
    const [filePreview, setFilePreview] = useState(null);
    const [fileName, setFileName] = useState("");
    const fileInputRef = useRef(null);

    const triggerFileInput = () => {
      console.log("Triggering file input ref");
      fileInputRef.current.click();
    }

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
    
        axios
          .post("http://localhost:8080/upload", formData, {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          })
          .then((response) => {
            console.log("File uploaded:", response);
          })
          .catch((error) => {
            console.error("Error uploading file:", error);
          });
      };

    return (
        <div className='w-[20vw] items-center justify-center ml-10'>
          <div className='flex flex-col w-full items-center justify-center'>
            <div className='items-center justify-center p-5 bg-[#272e3f] rounded-3xl space-y-2'>
              <div className='relative items-center justify-center'>
                <HiDocumentPlus size={60} color='white' onClick={triggerFileInput} onChange={triggerFileInput} className='w-full items-center justify-center cursor-pointer' />
                <input className="hidden" ref={fileInputRef} type="file" onChange={handleFileUpload} accept='.pdf, .txt, .doc, .docx, .word' />
              </div>
              <div>
                <span className='text-white text-xl text-opacity-70 items-center justify-center'>Add a file here</span>
              </div>
            </div>

            <div className='mt-10 pt-10'>
              <FileListDisplay />
            </div>
          </div>
          
        </div>
    );
}
