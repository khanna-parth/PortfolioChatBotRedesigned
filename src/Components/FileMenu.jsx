import React, {useState} from 'react'
import { HiDocumentAdd } from "react-icons/hi";
import axios from 'axios';

export default function FileMenu() {
    const [filePreview, setFilePreview] = useState(null);
    const [fileName, setFileName] = useState("");
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
        <div className='flex flex-col relative'>
          <input type="file" onChange={handleFileUpload} />
        </div>
    );
}
