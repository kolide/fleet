# Kolide Front-End

The Kolide front-end is a Single Page Application using React and Redux.

## Running the Kolide web app

For details instruction on building and serving the Kolide web application
consult the [Development Documentation](https://docs.kolide.co/kolide/current/development/index.html)

## Directory Structure

Component directories in the Kolide front-end application encapsulate the entire
component, including files for the component, helper functions, styles, and tests. The
typical directory structure for a component is as follows:

```
|-- ComponentName
|  |-- _styles.scss
|  |-- ComponentName.jsx
|  |-- ComponentName.tests.jsx
|  |-- helpers.js
|  |-- helpers.tests.js
|  |-- index.js
```

* `_styles.scss`: The component css styles
* `ComponentName.jsx`: The React component
* `ComponentName.tests.jsx`: The React component tests
* `helpers.js`: Helper functions used by the component
* `helpers.tests.js`: Tests for the component's helper functions
* `index.js`: Exports the React component
  * This file is helpful as it allows other components to import the component
    by it's directory name. Without this file the component name would have to
be duplicated during imports (`components/ComponentName` vs. `components/ComponentName/ComponentName`).

### [app_constants](https://github.com/kolide/kolide/tree/master/frontend/app_constants)

The app_constants directory exports the constants used in the app. Examples
include the app's URL paths, settings, and http statuses. When building features
that require constants, the constants should be added here for accessibility
throughout the application.

### [components](https://github.com/kolide/kolide/tree/master/frontend/components)
The component directory contains the React components rendered by pages. They
are typically not connected to the redux state but receive props from their
parent components to render data and handle user interactions.

### [interfaces](https://github.com/kolide/kolide/tree/master/frontend/interfaces)

Files in the interfaces directory are used to specify the PropTypes for a reusable Kolide
entity. This is designed to DRY up the code and increase re-usability. These
interfaces are imported into component files and implemented when defining the
component's PropTypes.

### [kolide](https://github.com/kolide/kolide/tree/master/frontend/kolide)

The default export of the kolide directory is the API client. More info can be
found at the [API client documentation page](https://github.com/kolide/kolide/tree/master/frontend/kolide/README.md).

### [layouts](https://github.com/kolide/kolide/tree/master/frontend/layouts)

The Kolide application has only 1 layout, the [Core Layout](https://github.com/kolide/kolide/blob/master/frontend/layouts/CoreLayout/CoreLayout.jsx).
The Layout is rendered from the [router](https://github.com/kolide/kolide/blob/master/frontend/router/index.jsx) and are used to set up the general app UI (header, sidebar) and render child components.
The child components rendered by the layout are typically page components.

### [pages](https://github.com/kolide/kolide/tree/master/frontend/pages)

Page components are React components typically rendered from the [router](https://github.com/kolide/kolide/tree/master/frontend/router).
These components are connected to redux state and are used to gather data from
redux and pass that data to child components (located in the [components
directory](https://github.com/kolide/kolide/tree/master/frontend/components). As
connected components, Pages are also used to dispatch actions. Actions
dispatched from Pages are intended to update redux state and oftentimes include
making a call to the Kolide API.

### [redux](https://github.com/kolide/kolide/tree/master/frontend/redux)

The redux directory holds all of the application's redux middleware, actions,
and reducers. The redux directory also creates the [store](https://github.com/kolide/kolide/blob/master/frontend/redux/store.js) which is used in the router.
More information about the redux configuration can be found at the [Redux
Documentation page]().

### [router](https://github.com/kolide/kolide/tree/master/frontend/router)

The router directory is where the react router lives. The router decides which
component will render at a given URL. Components rendered from the router are
typically located in the [pages directory](https://github.com/kolide/kolide/tree/master/frontend/pages). The router directory also holds a `paths`
file which holds the application paths as string constants for reference
throughout the app. These paths are typically referenced from the [App
Constants](https://github.com/kolide/kolide/tree/master/frontend/app_constants) object.

### [styles](https://github.com/kolide/kolide/tree/master/frontend/styles)

The styles directory contains the general app style setup and variables. It
includes variables for the app color hex codes, fonts (families, weights and sizes), and padding.

### [templates](https://github.com/kolide/kolide/tree/master/frontend/templates)

The templates directory contains the HTML file that renders the React application via including the `bundle.js`
 and `bundle.css` files. The HTML page also includes the HTML element in which the React application is mounted.

### [test](https://github.com/kolide/kolide/tree/master/frontend/test)

The test directory includes test helpers, API request mocks, and stubbed data entities for use in test files.
More on request mocks [here](https://github.com/kolide/kolide/tree/master/frontend/test/mocks/README.md).

### [utilities](https://github.com/kolide/kolide/tree/master/frontend/utilities)

The utilities directory contains re-usable functions for use throughout the
application. The functions include helpers to convert an array of objects to
CSV, debounce functions to prevent multiple form submissions, format API errors,
etc.

## Forms

For details on creating a Kolide form visit the [Kolide Form Documentation](https://github.com/kolide/kolide/tree/master/frontend/components/forms/README.md).

## API Client

For details on the Kolide API Client visit the [Kolide API Client Documentation](https://github.com/kolide/kolide/tree/master/frontend/kolide/README.md).
