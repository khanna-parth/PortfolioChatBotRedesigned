from dotenv import load_dotenv
from PyPDF2 import PdfReader
from langchain.text_splitter import CharacterTextSplitter
from langchain_openai import OpenAIEmbeddings
from sentence_transformers import SentenceTransformer
from langchain_community.vectorstores import FAISS
from langchain_openai import ChatOpenAI
from langchain.chains import RetrievalQA
import time
import faiss
import numpy as np
import concurrent.futures


load_dotenv()

def get_pdf_text(pdf_docs):
    """Extract text from the provided PDFs."""
    text = ""
    for pdf in pdf_docs:
        pdf_reader = PdfReader(pdf)
        for page in pdf_reader.pages:
            text += page.extract_text()
    return text

def encode_chunk(model, chunk):
    return model.encode(chunk)

def get_text_chunks(text):
    """Split the extracted text into smaller chunks."""
    text_splitter = CharacterTextSplitter(
        separator="\n",
        # chunk_size=1000,
        # chunk_size=2000,
        # chunk_overlap=200,
        chunk_overlap=400,
        length_function=len
    )
    chunks = text_splitter.split_text(text)
    return chunks

def get_vectorstore(text_chunks):
    """Convert text chunks into a vector store using sentence-transformers and a custom FAISS index."""
    model = SentenceTransformer('all-MiniLM-L6-v2')

    with concurrent.futures.ThreadPoolExecutor() as executor:
        embedded_texts = list(executor.map(lambda chunk: encode_chunk(model, chunk), text_chunks))

    vectors = np.array(embedded_texts).astype('float32')

    faiss_index = faiss.IndexFlatL2(vectors.shape[1])
    faiss_index.add(vectors)

    docstore = {i: text for i, text in enumerate(text_chunks)}
    index_to_docstore_id = {i: i for i in range(len(text_chunks))}

    vectorstore = FAISS(
        embedding_function=model.encode, 
        index=faiss_index,
        docstore=docstore,
        index_to_docstore_id=index_to_docstore_id
    )


    return vectorstore

def get_conversation_chain(vectorstore):
    """Create a conversational chain using the vectorstore and an LLM."""
    llm = ChatOpenAI(model="gpt-3.5-turbo")
    conversation_chain = RetrievalQA.from_llm(
        llm=llm,
        retriever=vectorstore.as_retriever(search_kwargs={"k": 5}),
    )
    return conversation_chain

def run(pdf_docs, prompt):
    """Main function to process PDFs and answer a prompt."""
    start = time.perf_counter()

    text = get_pdf_text(pdf_docs)
    
    text_chunks = get_text_chunks(text)
    
    vectorstore = get_vectorstore(text_chunks)
    
    qa_chain = get_conversation_chain(vectorstore)

    print(f"Setup took {time.perf_counter() - start}s")
    
    # response = conversation_chain.run({"question": prompt})
    # response = qa_chain.run(prompt)
    # response = qa_chain.invoke({"query": prompt})
    response = qa_chain.run(prompt)

    return response
    
    return response['query']

if __name__ == "__main__":
    # pdf_files = ["resume.pdf", "percyJackson.pdf"]
    pdf_files = ["resume.pdf"]
    
    prompt = "Give a thorough description of the applicant's skillset"
    # prompt = "How did he beat Medusa"
    
    response = run(pdf_files, prompt)
    
    print("Response:", response)
