import { defineBoot } from '#q-app/wrappers'
import VueSanitize from 'vue-sanitize-directive'

// "async" is optional;
// more info on params: https://v2.quasar.dev/quasar-cli-vite/boot-files
export default defineBoot(async ({ app }/* { app, router, ... } */) => {
  // something to do
  app.use(VueSanitize)
})
