# AuthApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**authGoogle**](AuthApi.md#authgoogle) | **GET** /api/v1/auth/google | Start Google OAuth flow |
| [**authGoogleCallback**](AuthApi.md#authgooglecallback) | **GET** /api/v1/auth/google/callback | Handle Google OAuth callback |
| [**authLogin**](AuthApi.md#authlogin) | **POST** /api/v1/auth/login | Login |
| [**authSignUp**](AuthApi.md#authsignup) | **POST** /api/v1/auth/signup | Register a new user |



## authGoogle

> authGoogle()

Start Google OAuth flow

Redirects the user-agent to Google and sets &#x60;oauth_state&#x60; cookie.

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@vdm2-bank/sdk';
import type { AuthGoogleRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const api = new AuthApi();

  try {
    const data = await api.authGoogle();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **307** | Temporary Redirect |  * Location - Google OAuth authorization URL <br>  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authGoogleCallback

> string authGoogleCallback(code, state)

Handle Google OAuth callback

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@vdm2-bank/sdk';
import type { AuthGoogleCallbackRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const api = new AuthApi();

  const body = {
    // string | OAuth code
    code: code_example,
    // string | CSRF state
    state: state_example,
  } satisfies AuthGoogleCallbackRequest;

  try {
    const data = await api.authGoogleCallback(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **code** | `string` | OAuth code | [Defaults to `undefined`] |
| **state** | `string` | CSRF state | [Defaults to `undefined`] |

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `text/html`, `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | HTML page that stores token in localStorage and redirects. |  -  |
| **400** | Bad request |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authLogin

> AuthResponse authLogin(loginRequest)

Login

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@vdm2-bank/sdk';
import type { AuthLoginRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const api = new AuthApi();

  const body = {
    // LoginRequest
    loginRequest: ...,
  } satisfies AuthLoginRequest;

  try {
    const data = await api.authLogin(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **loginRequest** | [LoginRequest](LoginRequest.md) |  | |

### Return type

[**AuthResponse**](AuthResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## authSignUp

> User authSignUp(signUpRequest)

Register a new user

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@vdm2-bank/sdk';
import type { AuthSignUpRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const api = new AuthApi();

  const body = {
    // SignUpRequest
    signUpRequest: ...,
  } satisfies AuthSignUpRequest;

  try {
    const data = await api.authSignUp(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **signUpRequest** | [SignUpRequest](SignUpRequest.md) |  | |

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | Bad request |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

