/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly BACKEND_HOST: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}