# Kolide Request Mocks

Request mocks are used to intercept API requests when running tests. Requests
are mocked to simulate valid and invalid requests. The naming convention is
similar to the [API client entity CRUD methods](https://github.com/kolide/kolide/tree/master/frontend/kolide/README.md).

## Using Mocks

```js
// import the mocks you want in the test file
import queryMocks from 'test/mocks/query_mocks';

// mock the API request before making the API call
queryMocks.load.valid(bearerToken, queryID); // valid request
queryMocks.load.invalid(bearerToken, queryID); // invalid request
```

Each entity with mocked requests has a dedicated file in this directory
containing the mocks for the entity, such as `queryMocks` in the example above. If requests
need to be mocked for multiple entities, consider importing all mocks:

```js
import mocks from 'test/mocks';

mocks.queries.load.valid(bearerToken, queryID);
mocks.packs.create.valid(bearerToken, params);
```

## Creating Mocks

Mocks are created using the [`createRequestMock`](https://github.com/kolide/kolide/tree/master/frontend/test/mocks/create_request_mock.js) function.

Example:

```js
// in /frontend/test/mocks/query_mocks.js
import createRequestMock from 'test/mocks/create_request_mock';
import { queryStub } from 'test/stubs';

const queryMocks = {
  load: {
    valid: (bearerToken, queryID) => {
      return createRequestMock({
        bearerToken,
        endpoint: `/api/v1/kolide/queries/${queryID}`,
        method: 'get',
        response: { query: { ...queryStub, id: queryID } },
        responseStatus: 200,
      });
    },
  },
}

export default queryMocks;
```

`createRequestMock` takes an options hash with the following options:

`bearerToken`

* Type: String
* Required?: False
* Default: None
* Purpose: Specifying the bearer token sets the Authorization header of the
  request and is often used when mocking authorized requests to the API.

`endpoint`

* Type: String
* Required?: True
* Default: None
* Purpose: The required endpoint option is the relative pathname of the request.

`method`

* Type: String (`get` | `post` | `patch` | `delete`)
* Required?: True
* Default: None
* Purpose: This string is the lower-cased request method. Options are `get`,
  `post`, `patch`, and `delete`.

`params`

* Type: Object
* Required?: False
* Default: None
* Purpose: This JS Object is for the parameters sent with a request. If the
  parameters are URL parameters, such as in a GET request, add the parameters to
the `endpoint` option.

`response`

* Type: Object
* Required?: True
* Default: None
* Purpose: This JS Object represents the response from the API

`responseStatus`

* Type: Number
* Required?: False
* Default: 200
* Purpose: This value is used for the response status of the API call.
