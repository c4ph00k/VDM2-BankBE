# AuthApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**authGoogle**](AuthApi.md#authGoogle) | **GET** /api/v1/auth/google | Start Google OAuth flow |
| [**authGoogleCallback**](AuthApi.md#authGoogleCallback) | **GET** /api/v1/auth/google/callback | Handle Google OAuth callback |
| [**authLogin**](AuthApi.md#authLogin) | **POST** /api/v1/auth/login | Login |
| [**authSignUp**](AuthApi.md#authSignUp) | **POST** /api/v1/auth/signup | Register a new user |


<a id="authGoogle"></a>
# **authGoogle**
> authGoogle()

Start Google OAuth flow

Redirects the user-agent to Google and sets &#x60;oauth_state&#x60; cookie.

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = AuthApi()
try {
    apiInstance.authGoogle()
} catch (e: ClientException) {
    println("4xx response calling AuthApi#authGoogle")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling AuthApi#authGoogle")
    e.printStackTrace()
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a id="authGoogleCallback"></a>
# **authGoogleCallback**
> kotlin.String authGoogleCallback(code, state)

Handle Google OAuth callback

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = AuthApi()
val code : kotlin.String = code_example // kotlin.String | OAuth code
val state : kotlin.String = state_example // kotlin.String | CSRF state
try {
    val result : kotlin.String = apiInstance.authGoogleCallback(code, state)
    println(result)
} catch (e: ClientException) {
    println("4xx response calling AuthApi#authGoogleCallback")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling AuthApi#authGoogleCallback")
    e.printStackTrace()
}
```

### Parameters
| **code** | **kotlin.String**| OAuth code | |
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **state** | **kotlin.String**| CSRF state | |

### Return type

**kotlin.String**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a id="authLogin"></a>
# **authLogin**
> AuthResponse authLogin(loginRequest)

Login

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = AuthApi()
val loginRequest : LoginRequest =  // LoginRequest | 
try {
    val result : AuthResponse = apiInstance.authLogin(loginRequest)
    println(result)
} catch (e: ClientException) {
    println("4xx response calling AuthApi#authLogin")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling AuthApi#authLogin")
    e.printStackTrace()
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **loginRequest** | [**LoginRequest**](LoginRequest.md)|  | |

### Return type

[**AuthResponse**](AuthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a id="authSignUp"></a>
# **authSignUp**
> User authSignUp(signUpRequest)

Register a new user

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = AuthApi()
val signUpRequest : SignUpRequest =  // SignUpRequest | 
try {
    val result : User = apiInstance.authSignUp(signUpRequest)
    println(result)
} catch (e: ClientException) {
    println("4xx response calling AuthApi#authSignUp")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling AuthApi#authSignUp")
    e.printStackTrace()
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **signUpRequest** | [**SignUpRequest**](SignUpRequest.md)|  | |

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

