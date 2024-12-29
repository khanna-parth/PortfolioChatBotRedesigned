from fastapi import FastAPI, Request, Depends, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Union
from chain import generator

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class UserRequest(BaseModel):
    userID: str

def get_user_store(request: UserRequest)-> generator.UserDocumentStore:
    print(f"Request body: {request}")
    userID = request.userID
    if userID == "":
        raise HTTPException(
            status_code=400,
            detail="userID is required"
        )
    return generator.UserDocumentStore(userID)

@app.get("/ping")
async def ping():
    return "Pong"

@app.post("/add-doc")
async def add_document(request: Request, userStore: generator.UserDocumentStore = Depends(get_user_store)):
    data = await request.json()
    verify_params(["userID", "docPath"], data)

    userID = data.get("userID")
    docPath = data.get("docPath")
    print(userID)
    print(docPath)
    print("Adding")
    # userStore = generator.UserDocumentStore(userID)
    try:
        userStore.add_file(docPath)
    except generator.BadPDF:
        raise HTTPException(
            status_code=422,
            detail="The PDF could not be processed"
        )
    return {"Message": "Document was added"}


@app.post("/remove-doc")
async def remove_doc(request: Request, userStore: generator.UserDocumentStore = Depends(get_user_store)):
    data = await request.json()
    verify_params(["userID", "docPath"], data)

    userID = data.get("userID")
    docPath = data.get("docPath")
    print(userID)
    print(docPath)
    print("Removing") 
    # userStore = generator.UserDocumentStore(userID)
    userStore.remove_documents(docPath)

    return {"Document was removed"}

@app.post("/list-docs")
async def list_docs(request: Request, userStore: generator.UserDocumentStore = Depends(get_user_store)):
    data = await request.json()
    verify_params(["userID"], data)

    userID = data.get("userID")
    # userStore = generator.UserDocumentStore(userID)
    docs = userStore.list_docs()
    print(f"Returning {userID}'s docs: {docs}")
    return {"documents": list(docs)}

@app.post("/prompt")
async def handle_prompt(request: Request, userStore: generator.UserDocumentStore = Depends(get_user_store)):
    data = await request.json()
    verify_params(["userID", "prompt"], data)

    userID = data.get("userID")
    prompt = data.get("prompt")

    # userStore = generator.UserDocumentStore(userID)
    resp = generator.answer_question(userStore, prompt)
    return {"Message": resp}

def verify_params(paramKeys, source) -> bool:
    for paramKey in paramKeys:
        if source.get(paramKey) == None or source.get(paramKey) == "":
            print("raising exception")
            raise HTTPException(
                status_code=400,
                detail=f"Parameter '{paramKey}' is required"
            )
            return false
    print("All checks passed")
    return True