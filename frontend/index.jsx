import ReactDOM from 'react-dom';
import runtime from 'serviceworker-webpack-plugin/lib/runtime';

import routes from './router';
import './index.scss';

if (typeof window !== 'undefined') {
  const { document } = global;
  const app = document.getElementById('app');

  ReactDOM.render(routes, app);
}

if ('serviceWorker' in navigator && process.env.NODE_ENV === 'production') {
  const registration = runtime.register();
}
