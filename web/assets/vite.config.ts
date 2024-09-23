import { defineConfig } from 'vite'

// https://vitejs.dev/config/
// @ts-ignore
export default defineConfig({
    base: "/meal-planner/assets",
    plugins: [],
    build: {
        manifest: 'vite-manifest.json',
        rollupOptions: {
            // overwrite default .html entry
            input: [
                '/src/index.css',
                '/src/index.ts',
            ]
        },
    },
    server: {
        origin: 'http://localhost:8080'
    },

})
