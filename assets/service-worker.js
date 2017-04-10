var serviceWorkerOption = {
  "assets": [
    "/kolide-logo-condensed@60a5addbf42470fd54a117be483ed491.svg",
    "/key@520ac59e6c6c189b94f4723c4e91d266.svg",
    "/404@f3676fd679816c107348ed2737a4897c.svg",
    "/500@b18abd32b7900af901327c55495fb0c7.svg",
    "/avatar@b3cfa572c321bac1e0bb50bfc9181d5f.svg",
    "/footer-logo@bd8b92e34e99f955afdd993acf667060.svg",
    "/kolide-logo@311a81c6f65bf0fbf98caf5770e1499c.svg",
    "/laptop-plus@50e51c3e80307960656836f2eb591c12.svg",
    "/osquery-certificate@8b14a15c345627a0faf7eb0648403fe8.svg",
    "/sign-up-pencil@81e812eadfadb2abf0471c6be52af5d8.svg",
    "/swoop-arrow@5a3a2402459ad28ca1f8f1de0ed29808.svg",
    "/bundle.js",
    "/bundle.css"
  ]
};
        
        /******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};

/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {

/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId])
/******/ 			return installedModules[moduleId].exports;

/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};

/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);

/******/ 		// Flag the module as loaded
/******/ 		module.l = true;

/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}


/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;

/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;

/******/ 	// identity function for calling harmony imports with the correct context
/******/ 	__webpack_require__.i = function(value) { return value; };

/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, {
/******/ 				configurable: false,
/******/ 				enumerable: true,
/******/ 				get: getter
/******/ 			});
/******/ 		}
/******/ 	};

/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};

/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };

/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "/assets/";

/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = 0);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/*
 Copyright 2015 Google Inc. All Rights Reserved.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/



// Incrementing CACHE_VERSION will kick off the install event and force previously cached
// resources to be cached again.

var CACHE_VERSION = 1;
var CURRENT_CACHES = {
  offline: 'offline-v' + CACHE_VERSION
};
var OFFLINE_URL = '/offline';

function createCacheBustedRequest(url) {
  var request = new Request(url, { cache: 'reload' });
  // See https://fetch.spec.whatwg.org/#concept-request-mode
  // This is not yet supported in Chrome as of M48, so we need to explicitly check to see
  // if the cache: 'reload' option had any effect.
  if ('cache' in request) {
    return request;
  }

  // If {cache: 'reload'} didn't have any effect, append a cache-busting URL parameter instead.
  var bustedUrl = new URL(url, self.location.href);
  bustedUrl.search += (bustedUrl.search ? '&' : '') + 'cachebust=' + Date.now();
  return new Request(bustedUrl);
}

self.addEventListener('install', function (event) {
  event.waitUntil(
  // We can't use cache.add() here, since we want OFFLINE_URL to be the cache key, but
  // the actual URL we end up requesting might include a cache-busting parameter.
  fetch(createCacheBustedRequest(OFFLINE_URL)).then(function (response) {
    return caches.open(CURRENT_CACHES.offline).then(function (cache) {
      return cache.put(OFFLINE_URL, response);
    });
  }));
});

self.addEventListener('activate', function (event) {
  // Delete all caches that aren't named in CURRENT_CACHES.
  // While there is only one cache in this example, the same logic will handle the case where
  // there are multiple versioned caches.
  var expectedCacheNames = Object.keys(CURRENT_CACHES).map(function (key) {
    return CURRENT_CACHES[key];
  });

  event.waitUntil(caches.keys().then(function (cacheNames) {
    return Promise.all(cacheNames.map(function (cacheName) {
      if (expectedCacheNames.indexOf(cacheName) === -1) {
        // If this cache name isn't present in the array of "expected" cache names,
        // then delete it.
        console.log('Deleting out of date cache:', cacheName);
        return caches.delete(cacheName);
      }
    }));
  }));
});

self.addEventListener('fetch', function (event) {
  // We only want to call event.respondWith() if this is a navigation request
  // for an HTML page.
  // request.mode of 'navigate' is unfortunately not supported in Chrome
  // versions older than 49, so we need to include a less precise fallback,
  // which checks for a GET request with an Accept: text/html header.
  if (event.request.mode === 'navigate' || event.request.method === 'GET' && event.request.headers.get('accept').includes('text/html')) {
    console.log('Handling fetch event for', event.request.url);
    event.respondWith(fetch(event.request).catch(function (error) {
      // The catch is only triggered if fetch() throws an exception, which will most likely
      // happen due to the server being unreachable.
      // If fetch() returns a valid HTTP response with an response code in the 4xx or 5xx
      // range, the catch() will NOT be called. If you need custom handling for 4xx or 5xx
      // errors, see https://github.com/GoogleChrome/samples/tree/gh-pages/service-worker/fallback-response
      console.log('Fetch failed; returning offline page instead.', error);
      return caches.match(OFFLINE_URL);
    }));
  }

  // If our if() condition is false, then this fetch handler won't intercept the request.
  // If there are any other fetch handlers registered, they will get a chance to call
  // event.respondWith(). If no fetch handlers call event.respondWith(), the request will be
  // handled by the browser as if there were no service worker involvement.
});

/***/ })
/******/ ]);