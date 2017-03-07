# Kolide Redux Middleware

The Kolide Redux Middleware handles actions before they hit the reducers. The
current middleware does the following:

## [Authentication Middleware](https://github.com/kolide/kolide/blob/master/frontend/redux/middlewares/auth.js)

The authentication middleware handles logging a user in/out and handles logging out a user when the API responds
with an unauthenticated error.

## [Redirect Middleware](https://github.com/kolide/kolide/blob/master/frontend/redux/middlewares/redirect/index.js)

The redirect middleware transitions the user to the 500 page when an API call
fails with a 500 status.

## [Nag Message Middleware](https://github.com/kolide/kolide/blob/master/frontend/redux/middlewares/nag_message/index.js)

The Nag Message Middleware handles displaying a persistent flash message when a
Kolide user should update their Kolide license.
