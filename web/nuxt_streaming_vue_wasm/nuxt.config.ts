// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  nitro: {
    routeRules: {
      '/api': { proxy: 'http://localhost:3001' },
    }
  }
})
