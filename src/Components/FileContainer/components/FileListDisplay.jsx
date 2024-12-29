import React, { useEffect, useState } from 'react'
import { FaFilePdf, FaFileWord} from "react-icons/fa";
import { FiFileText } from "react-icons/fi";
import { TbFileTypeDocx } from "react-icons/tb";
import { apiClient } from '../../../lib/client';

import { FaTrash } from "react-icons/fa";
import { DELETE_DOC_ROUTE, LIST_DOCS_ROUTE, PROMPT_ROUTE } from '../../../routes/routes.js';



const FileListDisplay= () => {
    // const files = ["resume.pdf", "coverLetter.pdf", "reference.txt", "document.docx"];
    const [files, setFiles] = useState([]);

    const getFiles = async () => {
        const data = { userID: "pkhanna"};

        try {
            const response = await apiClient.post(LIST_DOCS_ROUTE, data)
            console.log(response)
            setFiles(response.data.documents)

        } catch (error) {
            console.log(error)
        }
    }

    const deleteDocument = async (documentName) => {
        const data = {
            userID: "pkhanna",
            document: "resume.pdf"
        }

        try {
            const response = await apiClient.post(DELETE_DOC_ROUTE, data);
            console.log(response);
            if (response.status == 200) {
                getFiles()
            }
        } catch (error) {
            console.log(error);
        }
    }

    useEffect(() => {
        getFiles()
    }, []);


  return (
    <div className='flex flex-col flex-grow h-full space-y-4'>
        {
            files.map((eachFile) => (
                <div key={eachFile} className='flex items-center justify-start space-x-4'>
                    {

                    }
                    <FaFilePdf color='white' />
                    <span className='text-[#DBEBC0]'>{eachFile}</span>
                    <FaTrash onClick={() => deleteDocument(eachFile)} color='red' />
                </div>
            )
        )}
    </div>
  )
}

export default FileListDisplay