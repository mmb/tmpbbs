import { defineConfig } from "eslint/config"
import globals from "globals"
import js from "@eslint/js"

export default defineConfig([
  {
    extends: ["js/all"],
    files: ["**/*.{js,mjs,cjs}"],
    languageOptions: { globals: globals.browser },
    linterOptions: {
      reportUnusedDisableDirectives: "error",
      reportUnusedInlineConfigs: "error"
    },
    plugins: { js },
    rules: {
      "comma-dangle":  ["error", "never"],
      "no-magic-numbers": [ "error", { "ignore" : [0, 1] }],
      "semi":  ["error", "never"]
    }
  }
])
