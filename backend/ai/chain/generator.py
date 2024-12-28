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

        # Ensure the storage directory exists
        if not os.path.exists(self.storage_dir):
            os.makedirs(self.storage_dir)

        # Check the contents of the storage directory before loading
        print(f"Checking for vector store in: {self.storage_dir}")
        self.load_vector_store()

    def load_vector_store(self):
        vector_store_path = f"{self.storage_dir}/chroma.sqlite3"
        
        if os.path.exists(vector_store_path):
            print(f"Loading existing vector store for user {self.user_id} from {vector_store_path}")
            self.vector_store = Chroma(persist_directory=self.storage_dir, embedding_function=self.embeddings)
            self.retriever = self.vector_store.as_retriever()
            self.rag_chain = create_retrieval_chain(self.retriever, question_answer_chain)
        else:
            print(f"No existing vector store for user {self.user_id} at {vector_store_path}, starting fresh.")

    def add_file(self, file_path):
        print(f"Processing file for user {self.user_id}:", file_path)
        loader = PyPDFLoader(file_path)
        docs = loader.load()

        text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=200)
        chunks = text_splitter.split_documents(docs)

        # Create or update vector store with documents for this user
        if self.vector_store is None:
            print(f"Creating new vector store for user {self.user_id}")
            self.vector_store = Chroma.from_documents(chunks, self.embeddings, persist_directory=self.storage_dir)
        else:
            print(f"Adding documents to existing vector store for user {self.user_id}")
            self.vector_store.add_documents(chunks)

        self.retriever = self.vector_store.as_retriever()
        self.rag_chain = create_retrieval_chain(self.retriever, question_answer_chain)

    def remove_documents(self, document_ids):
        print(f"Removing documents for user {self.user_id}: {document_ids}")
        print(dir(self.vector_store))
        self.vector_store.delete_document(document_ids)
        self.retriever = self.vector_store.as_retriever()
        self.rag_chain = create_retrieval_chain(self.retriever, question_answer_chain)

    def get_rag_chain(self):
        return self.rag_chain

def answer_question(rag_chain, question):
    response = rag_chain.invoke({"input": question})
    print("Answer:", response["answer"])
    return response['answer']

# if __name__ == "__main__":
#     user_id = input("Enter your user ID: ")

#     user_store = UserDocumentStore(user_id)

#     while True:
#         add_doc = input("Do you want to add a new document? (y/n): ")
#         if add_doc.lower() == 'y':
#             new_file_path = input("Enter the new file path: ")
#             user_store.add_file(new_file_path)

#         remove_doc = input("Do you want to remove some documents? (y/n): ")
#         if remove_doc.lower() == 'y':
#             document_ids = input("Enter document IDs to remove (comma separated): ").split(',')
#             user_store.remove_documents([doc_id.strip() for doc_id in document_ids])

#         question = input("Enter your question: ")
#         if question == "":
#             break

#         answer_question(user_store.get_rag_chain(), question)

