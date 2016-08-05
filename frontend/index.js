const router = require('./router');

if(typeof window !== 'undefined') {
  router.run();
}
