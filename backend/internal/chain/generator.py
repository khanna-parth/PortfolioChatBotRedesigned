from decouple import config
from langchain_openai import ChatOpenAI, OpenAIEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_chroma import Chroma
from langchain.chains import create_retrieval_chain
from langchain.chains.combine_documents import create_stuff_documents_chain
from langchain_core.prompts import ChatPromptTemplate
from langchain_community.document_loaders import PyPDFLoader 
import os
import argparse
from dotenv import load_dotenv

load_dotenv()

llm = ChatOpenAI(model="gpt-3.5-turbo")

class BadPDF(Exception):
    pass

contextualize_system_prompt = """Given a chat history and the latest user question \
which might reference context in the chat history, formulate a standalone question \
which can be understood without the chat history. Do NOT answer the question, \
just reformulate it if needed and otherwise return it as is."""
contextualize_prompt = ChatPromptTemplate.from_messages(
    [
        ("system", contextualize_system_prompt),
        ("human", "{input}"),
    ]
)

system_prompt = (
    "You are an assistant for question-answering tasks. "
    "Use the following pieces of retrieved context to answer "
    "the question. If you don't know the answer, say that you "
    "don't know. Use three sentences maximum and keep the "
    "answer concise."
    "\n\n"
    "{context}"
)

prompt = ChatPromptTemplate.from_messages(
    [
        ("system", system_prompt),
        ("human", "{input}"),
    ]
)

question_answer_chain = create_stuff_documents_chain(llm, prompt)

class UserDocumentStore:
    def __init__(self, user_id):
        self.user_id = user_id
        self.vector_store = None
        self.retriever = None
        self.rag_chain = None
        self.embeddings = OpenAIEmbeddings()
        self.storage_dir = f"storage/{user_id}"
        self.docs_stored = set()

        # Ensure the storage directory exists
        if not os.path.exists(self.storage_dir):
            os.makedirs(self.storage_dir)

        print(f"Checking for vector store in: {self.storage_dir}")
        self.load_vector_store()
    
    def update_refs(self):
        docs = []
        if not self.vector_store:
            return
        for x in range(len(self.vector_store.get()["ids"])):
                doc = self.vector_store.get()["metadatas"][x]
                source = doc["source"]
                docs.append(source)
        self.docs_stored = set(docs)

    def load_vector_store(self):
        vector_store_path = f"{self.storage_dir}/chroma.sqlite3"
        
        if os.path.exists(vector_store_path):
            print(f"Loading existing vector store for user {self.user_id} from {vector_store_path}")
            self.vector_store = Chroma(persist_directory=self.storage_dir, embedding_function=self.embeddings)

            self.retriever = self.vector_store.as_retriever()
            self.rag_chain = create_retrieval_chain(self.retriever, question_answer_chain)
        else:
            print(f"No existing vector store for user {self.user_id} at {vector_store_path}, starting fresh.")
        self.update_refs()

    def add_file(self, file_path):
        print(f"Processing file for user {self.user_id}:", file_path)
        loader = PyPDFLoader(file_path)
        docs = loader.load()

        if not docs or all(len(doc.page_content.strip()) == 0 for doc in docs):
            raise BadPDF("PDF parsed content is empty")

        text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=200)
        chunks = text_splitter.split_documents(docs)

        if not chunks or all(len(chunk.page_content.strip()) == 0 for chunk in chunks):
            raise BadPDF("PDF chunked content is empty")

        # Create or update vector store with documents for this user
        if self.vector_store is None:
            print(f"Creating new vector store for user {self.user_id}")
            self.vector_store = Chroma.from_documents(chunks, self.embeddings, persist_directory=self.storage_dir)
        else:
            print(f"Adding documents to existing vector store for user {self.user_id}")
            self.vector_store.add_documents(chunks)

        self.retriever = self.vector_store.as_retriever()
        self.rag_chain = create_retrieval_chain(self.retriever, question_answer_chain)
        self.update_refs()

    def remove_documents(self, document_ids):
        print(f"Removing documents for user {self.user_id}: {document_ids}")
        print(dir(self.vector_store))
        self.vector_store.delete_document(document_ids)
        self.retriever = self.vector_store.as_retriever()
        self.rag_chain = create_retrieval_chain(self.retriever, question_answer_chain)
        self.update_refs()

    def get_rag_chain(self):
        self.update_refs()
        return self.rag_chain
    
    def list_docs(self):
        return self.docs_stored

def answer_question(user_store: UserDocumentStore, question):
    rag_chain = user_store.get_rag_chain()
    if not rag_chain:
        return "[No documents added]"
    response = rag_chain.invoke({"input": question})
    print(f"User has stored docs: {user_store.docs_stored}")
    print("Answer:", response["answer"])
    return response['answer']