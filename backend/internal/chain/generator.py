import os
import json
from dotenv import load_dotenv
from langchain_community.document_loaders import PyPDFLoader
from langchain.text_splitter import CharacterTextSplitter
from langchain_openai import OpenAIEmbeddings
from langchain.chains.question_answering import load_qa_chain
from langchain_openai import ChatOpenAI
from langchain_community.vectorstores import Chroma
from dotenv import load_dotenv
from concurrent.futures import ProcessPoolExecutor
from sentence_transformers import SentenceTransformer
import time
from langchain.embeddings.base import Embeddings
from typing import List
import concurrent.futures
import argparse
import warnings
from langchain.schema import HumanMessage

'''
Serves as refined command line standalone version of generator
'''


warnings.filterwarnings("ignore", category=DeprecationWarning)

class CustomEmbeddings(Embeddings):
    def __init__(self, model_name: str = "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"):
        self.model = SentenceTransformer(model_name)

    def embed_documents(self, texts: List[str]) -> List[List[float]]:
        return [self.model.encode(text).tolist() for text in texts]

    def embed_query(self, text: str) -> List[float]:
        return self.model.encode([text])[0].tolist()

def process_pdf(pdf_path):
    loader = PyPDFLoader(pdf_path)
    return loader.load()

def load_chunk_persist_pdf() -> Chroma:
    pdf_folder_path = os.getcwd()
    pdf_files = [os.path.join(pdf_folder_path, file) for file in os.listdir(pdf_folder_path) if file.endswith('.pdf')]
    
    documents = []
    # with ProcessPoolExecutor() as executor:
    with concurrent.futures.ThreadPoolExecutor() as executor:
        results = executor.map(process_pdf, pdf_files)
        for result in results:
            documents.extend(result)
    
    text_splitter = CharacterTextSplitter(chunk_size=1000, chunk_overlap=10)
    chunked_documents = text_splitter.split_documents(documents)
    
    persist_directory = os.path.join(os.getcwd(), 'vector_store')
    embedding_model = OpenAIEmbeddings()
    
    if os.path.exists(persist_directory):
        # print("Loading existing vector store...")
        vectordb = Chroma(persist_directory=persist_directory, embedding_function=embedding_model)
    else:
        # print("Creating a new vector store...")
        vectordb = Chroma.from_documents(
            documents=chunked_documents,
            embedding=embedding_model,
            persist_directory=persist_directory
        )
        vectordb.persist()
    
    return vectordb

def add_documents_to_vector_store(new_documents, persist_directory):
    text_splitter = CharacterTextSplitter(chunk_size=1000, chunk_overlap=10)
    chunked_documents = text_splitter.split_documents(new_documents)
    
    vectordb = Chroma(persist_directory=persist_directory, embedding_function=CustomEmbeddings())
    
    vectordb.add_documents(chunked_documents)
    
    vectordb.persist()
    print(f"Added {len(chunked_documents)} new documents.")
    
    return vectordb

def remove_documents_from_vector_store(documents_to_remove, persist_directory):
    vectordb = Chroma(persist_directory=persist_directory, embedding_function=CustomEmbeddings())
    
    existing_documents = vectordb.get_documents()
    
    updated_documents = [doc for doc in existing_documents if doc not in documents_to_remove]
    
    vectordb = Chroma.from_documents(
        documents=updated_documents,
        embedding=CustomEmbeddings(),
        persist_directory=persist_directory
    )
    
    vectordb.persist()
    print(f"Removed {len(documents_to_remove)} documents.")
    
    return vectordb

def create_agent_chain():
    model_name = "gpt-3.5-turbo"
    llm = ChatOpenAI(name=model_name)
    chain = load_qa_chain(llm, chain_type="stuff", verbose=False)
    return chain


def get_llm_response(query):
    vectordb = load_chunk_persist_pdf()

    chain = create_agent_chain()
    matching_docs = vectordb.similarity_search(query)
    answer = chain.run(input_documents=matching_docs, question=query)
    return answer

# while True:
#     prompt = input("Enter prompt: ")
#     start = time.perf_counter()
#     print(get_llm_response(prompt))
#     print(f"Took {time.perf_counter() - start}s")

def worker_function(query):
    response = get_llm_response(query)
    return response

def simple_prompt(query):
    chat = ChatOpenAI(model="gpt-3.5-turbo")
    response = chat.invoke(query)
    return response.content

# def simulate_multiple_users(prompts):
#     with multiprocessing.Pool(processes=multiprocessing.cpu_count()/2) as pool:
#         responses = pool.map(worker_function, prompts)

#     for prompt, response in zip(prompts, responses):
#         print(f"Prompt: {prompt}")
#         print(f"Response: {response}\n")

# if __name__ == '__main__':
    # prompts = [
    #     "How is percy's relationship with his dad?",
    #     "How is percy's relationship with his Gabe?",
    #     "How did Percy escape the Underworld from Hades",
    #     "How was Aphrodite described when compared to her relationship with Hades",
    #     "How is Grover's appearance described?"
    # ]

    # simulate_multiple_users(prompts)

parser = argparse.ArgumentParser()
parser.add_argument('--dir', type=str, required=True)
parser.add_argument('--query', type=str, required=True)
parser.add_argument("--key", type=str, required=True)
args = parser.parse_args()

os.environ["OPENAI_API_KEY"] = args.key

os.chdir(args.dir)

start = time.perf_counter()

if len(os.listdir()) <= 1:
    resp = simple_prompt(args.query)
else:
    resp = worker_function(args.query)

data = {
    "prompt": args.query,
    "response": resp,
    "elapsed": time.perf_counter() - start,
}
print(json.dumps(data))