import React, { useEffect, useState } from 'react'
import { FaFilePdf, FaFileWord} from "react-icons/fa";
import { apiClient } from '../../../lib/client.js';
import { FaTrash } from "react-icons/fa";
import { DELETE_DOC_ROUTE, LIST_DOCS_ROUTE, PROMPT_ROUTE } from '../../../routes/routes.js';
import { LoaderIcon } from 'lucide-react';
import { useStore } from '../../../store/store.js';
import { useWebSocketStore } from '../../../store/store.js';
import { toast } from 'sonner';


const FileListDisplay= () => {
    const [files, setFiles] = useState([]);
    const [fetchComplete, setFetchComplete] = useState(false);
    const { setAllowUploads, setIsProcessing, setChatHistory } = useStore();
    const { connID, docCount } = useWebSocketStore();

    const getFiles = async () => {
        const data = { userID: "pkhanna"};
        setFetchComplete(false);

        try {
            const headers = {
                'X-Connection-ID': connID
            };
            console.log("Prepared Headers:", headers);
        
            const response = await apiClient.post(LIST_DOCS_ROUTE, data, { headers });
        
            console.log(response);

            if (response.status === 200) {
                setFiles(response.data.documents)
                setFetchComplete(true);
                setAllowUploads(true);
            } else if (response.status === 503) {
                setAllowUploads(false);
                setFetchComplete(true);
                toast.error("Server capacity full")
            }            

        } catch (error) {
            const errorMsg = `Error retrieving documents: ${error.response.data}`
            toast.error(errorMsg)
            console.log(error)
        }
    }

    const deleteDocument = async (documentName) => {
        const data = {
            document: documentName
        };

        const headers = {
            "X-Connection-ID": connID,
        };

        try {
            setIsProcessing(true);
            const response = await apiClient.post(DELETE_DOC_ROUTE, data, { headers });
            console.log(response);
            if (response.status == 200) {
                setChatHistory({docChange: `You removed ${documentName} from the chat context`});
                setIsProcessing(false);
                getFiles()
            }
            setIsProcessing(false);
        } catch (error) {
            setIsProcessing(false);
            console.log(error);
            const errorMsg = `Error deleting document: ${error.response.data}`
            toast.error(errorMsg)
        }
    }

    useEffect(() => {
        if (connID === null || connID === undefined) return;
        console.log("ConnID state in useEffect:", connID)
        console.log("Running useEffect getFiles")
        getFiles()
    }, [connID, docCount]);

    const isIDPreset = (id) => {
        if (id.startsWith("@")) {
            return true;
          }
          return false;
    }


  return (
    <div className='flex flex-col flex-grow h-full space-y-4'>
        {!fetchComplete && (
                <div className='flex flex-col items-center justify-center space-y-8'>
                    <h1 className='text-white text-lg'>Retrieving your files</h1>
                    <LoaderIcon className='text-white text-xl animate-spin duration-50'></LoaderIcon>
                </div>
            )
        }
        {fetchComplete && (
            <>
                {files.length === 0 ? (
                    <div className='w-full text-white text-md flex flex-col items-center justify-center space-y-2'>
                        <h1 className='w-full text-center'>No uploads</h1>
                        <span className='text-center text-white/50'>Add a file above</span>
                    </div>
                ) : (
                    files.map((eachFile) => (
                    <div key={eachFile} className='flex items-center justify-start space-x-4'>
                        <FaFilePdf color='white' />
                        <span className='text-[#DBEBC0]'>{eachFile}</span>
                        <FaTrash className={`cursor-pointer ${isIDPreset(connID) ? "hidden" : ""}`} onClick={() => deleteDocument(eachFile)} color='red' />
                    </div>
                    ))
                )}
            </>
        )}
    </div>
  )
}

export default FileListDisplay