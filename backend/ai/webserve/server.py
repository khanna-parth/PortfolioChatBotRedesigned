from fastapi import FastAPI, Request
from pydantic import BaseModel
from chain import generator

app = FastAPI()

class DocumentRequest(BaseModel):
    userID: str
    docPath: str

class PromptRequest(BaseModel):
    userID: str
    prompt: str

@app.post("/add-doc")
async def add_document(documentRequest: DocumentRequest):
    userID = documentRequest.userID
    docPath = documentRequest.docPath
    print(userID)
    print(docPath)
    print("Adding")
    userStore = generator.UserDocumentStore(userID)
    userStore.add_file(docPath)
    return


@app.post("/remove-doc")
def remove_doc(documentRequest: DocumentRequest):
    userID = documentRequest.userID
    docPath = documentRequest.docPath
    print(userID)
    print(docPath)
    print("Removing") 
    userStore = generator.UserDocumentStore(userID)
    userStore.remove_documents(docPath)

@app.post("/prompt")
def handle_prompt(promptRequest: PromptRequest):
    userID = promptRequest.userID
    prompt = promptRequest.prompt

    userStore = generator.UserDocumentStore(userID)
    resp = generator.answer_question(userStore.get_rag_chain(), prompt)
    return resp