import { defineConfig } from 'vite'
import react, { reactCompilerPreset } from '@vitejs/plugin-react'
import babel from '@rolldown/plugin-babel'
import mkcert from 'vite-plugin-mkcert'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    mkcert({ savePath: './certs' }),
    babel({ presets: [reactCompilerPreset()] })
  ],
  server: {
    host: true,
    port: 8082,
    cors: true,
    https: true,
    proxy: {
      '/server': {
        target: 'https://localhost:8080',
        changeOrigin: true,
        secure: false,
        ws: true,
        rewrite: (path) => path.replace(/^\/server/, '')
      }
    },
    hmr: {
      host: 'localhost',
      protocol: 'wss',
    }
  }
})
