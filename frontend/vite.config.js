import react from '@vitejs/plugin-react'
import path from "path"
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, path.resolve(__dirname, 'config'), '')


  const defineEnv = {}
  for (const key in env) {
    if (key.startsWith('VITE_')) {
      defineEnv[`import.meta.env.${key}`] = JSON.stringify(env[key])
    }
  }

  return {
    plugins: [react()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    define: defineEnv,
  }
})