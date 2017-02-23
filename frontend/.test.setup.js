import jsdom from 'jsdom';
import webdriver from 'selenium-webdriver';
import chrome from 'selenium-webdriver/chrome';
import chromedriver from 'chromedriver';

const doc = jsdom.jsdom(
  `<!doctype html>
  <html>
    <body>
      <input id="method1" value="hello world" />
      <input id="method2" value="hello world" />
    </body>
  </html>`,
  {
    url: 'http://localhost:8080/foo'
  },
);

global.document = doc;
global.document.queryCommandEnabled = () => { return true; };
global.document.execCommand = () => { return true; };
global.window = doc.defaultView;
global.window.getSelection = () => {
  return {
    removeAllRanges: () => { return true; },
  };
};
global.navigator = global.window.navigator;
chrome.setDefaultService(new chrome.ServiceBuilder(chromedriver.path).build());
global.webDriver = new webdriver.Builder()
  .withCapabilities(webdriver.Capabilities.chrome())
  .build();

function mockStorage() {
  let storage = {};

  return {
    setItem(key, value = '') {
      storage[key] = value;
    },
    getItem(key) {
      return storage[key];
    },
    removeItem(key) {
      delete storage[key];
    },
    get length() {
      return Object.keys(storage).length;
    },
    key(i) {
      return Object.keys(storage)[i] || null;
    },
    clear () {
      storage = {};
    },
  };
}

global.localStorage = window.localStorage = mockStorage();
