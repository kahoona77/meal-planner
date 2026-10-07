import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
    base: "/meal-planner/assets",
    plugins: [],
    build: {
        manifest: 'vite-manifest.json',
        rollupOptions: {
            // overwrite default .html entry
            input: [
                '/src/index.css',
            ]
        },
    },
    server: {
        origin: 'http://localhost:8080'
    },

})
