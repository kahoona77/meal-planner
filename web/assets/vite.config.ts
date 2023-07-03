import { defineConfig } from 'vite'

// https://vitejs.dev/config/
// @ts-ignore
export default defineConfig({
    base: "/meal-planner/assets",
    plugins: [],
    build: {
        // generate manifest.json in outDir
        manifest: true,
        rollupOptions: {
            // overwrite default .html entry
            input: [
                '/src/index.css',
                '/src/my-element.ts',
            ]
        },
    },
    server: {
        origin: 'http://localhost:8080'
    },

})
