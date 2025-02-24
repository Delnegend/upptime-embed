import path from "node:path";

import { VitePWA } from "vite-plugin-pwa";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: "2024-11-01",
	devtools: { enabled: true },
	css: ["./app.css"],
	ssr: false,
	imports: { scan: false, autoImport: false },
	components: false,
	postcss: {
		plugins: {
			tailwindcss: {},
			autoprefixer: {},
		},
	},
	vite: {
		resolve: {
			alias: {
				"~": path.resolve(__dirname, "./"),
			},
		},
		plugins: [VitePWA({
			registerType: "autoUpdate",
		})],
	},
	devServer: {
		port: 3000,
	},
	modules: ["@vite-pwa/nuxt"],
	app: {
		head: {
			title: "Upptime Embed",
			meta: [
				{ charset: "utf-8" },
				{ name: "viewport", content: "width=device-width, initial-scale=1" },
				{ hid: "description", name: "description", content: "Uptime embed for Upptime" },
			],
		},
	},
});
