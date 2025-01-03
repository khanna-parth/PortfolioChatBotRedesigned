# from sentence_transformers import SentenceTransformer
# from langchain.vectorstores import Chroma
# from langchain.embeddings import Embedding
# from langchain.llms import LLM
# from langchain.prompts import PromptTemplate
# from langchain.chains import LLMChain
# from transformers import LlamaForCausalLM, LlamaTokenizer
# from langchain_openai import ChatOpenAI, OpenAIEmbeddings
# import torch


# def get_llm(model_type: str):
#     if model_type == "openai":
#         return ChatOpenAI(model="gpt-3.5-turbo")
#     elif model_type == "llama":
#         return LlamaModel(model_name="meta-llama/Llama-2-7b-hf")
#     else:
#         raise ValueError(f"Unsupported model type: {model_type}")

# class LlamaModel:
#     def __init__(self, model_name: str):
#         self.model_name = model_name
#         self.model = LlamaForCausalLM.from_pretrained(model_name)
#         self.tokenizer = LlamaTokenizer.from_pretrained(model_name)
#         self.device = "cuda" if torch.cuda.is_available() else "cpu"
#         self.model.to(self.device)

#     def _call(self, prompt: str) -> str:
#         inputs = self.tokenizer(prompt, return_tensors="pt").to(self.device)
#         outputs = self.model.generate(**inputs, max_length=200, num_return_sequences=1, temperature=0.7)
#         return self.tokenizer.decode(outputs[0], skip_special_tokens=True)

# class LlamaEmbeddings:
#     def __init__(self, model_name: str = 'all-MiniLM-L6-v2'):
#         self.model = SentenceTransformer(model_name)
    
#     def embed_documents(self, texts: list) -> list:
#         return self.model.encode(texts, show_progress_bar=True).tolist()
    
#     def embed_query(self, query: str) -> list:
#         return self.model.encode([query])[0].tolist()