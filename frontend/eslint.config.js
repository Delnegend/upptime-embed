import hagemanto from "eslint-plugin-hagemanto"
import vue from "eslint-plugin-vue"

export default [
	{ files: ["**/*.{vue,ts}"] },
	{ ignores: ["**/*.d.ts", "**/*.config.js", ".nuxt/**/*"] },
	...hagemanto({ vueConfig: vue.configs["flat/recommended"] }),
]