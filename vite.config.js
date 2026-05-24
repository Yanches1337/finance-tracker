import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
    plugins: [react()],
    // Указываем явный корень для Vite, где лежит index.html и исходники
    root: 'frontend/finance-app',
    base: '/',
    build: {
        // Говорим сборщику складывать готовый билд в общую папку dist в корне,
        // чтобы Dockerfile (или Nginx) легко его забрал
        outDir: '../../dist',
        emptyOutDir: true,
    }
})