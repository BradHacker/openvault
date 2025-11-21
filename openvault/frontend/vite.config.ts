import { defineConfig } from 'vite';
import viteReact from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import wails from '@wailsio/runtime/plugins/vite';

import { tanstackRouter } from '@tanstack/router-plugin/vite';
import { fileURLToPath, URL } from 'node:url';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    tanstackRouter({
      target: 'react',
      autoCodeSplitting: true
    }),
    viteReact(),
    tailwindcss(),
    wails('./bindings')
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@wails': fileURLToPath(
        new URL(
          './bindings/github.com/wailsapp/wails/v3/internal',
          import.meta.url
        )
      ),
      '@openvault': fileURLToPath(
        new URL('./bindings/github.com/BradHacker/openvault', import.meta.url)
      )
    }
  }
});
