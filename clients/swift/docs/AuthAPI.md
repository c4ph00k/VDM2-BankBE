# AuthAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**authGoogle**](AuthAPI.md#authgoogle) | **GET** /api/v1/auth/google | Start Google OAuth flow
[**authGoogleCallback**](AuthAPI.md#authgooglecallback) | **GET** /api/v1/auth/google/callback | Handle Google OAuth callback
[**authLogin**](AuthAPI.md#authlogin) | **POST** /api/v1/auth/login | Login
[**authSignUp**](AuthAPI.md#authsignup) | **POST** /api/v1/auth/signup | Register a new user


# **authGoogle**
```swift
    open class func authGoogle(completion: @escaping (_ data: Void?, _ error: Error?) -> Void)
```

Start Google OAuth flow

Redirects the user-agent to Google and sets `oauth_state` cookie.

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK


// Start Google OAuth flow
AuthAPI.authGoogle() { (response, error) in
    guard error == nil else {
        print(error)
        return
    }

    if (response) {
        dump(response)
    }
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

Void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authGoogleCallback**
```swift
    open class func authGoogleCallback(code: String, state: String, completion: @escaping (_ data: String?, _ error: Error?) -> Void)
```

Handle Google OAuth callback

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK

let code = "code_example" // String | OAuth code
let state = "state_example" // String | CSRF state

// Handle Google OAuth callback
AuthAPI.authGoogleCallback(code: code, state: state) { (response, error) in
    guard error == nil else {
        print(error)
        return
    }

    if (response) {
        dump(response)
    }
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **code** | **String** | OAuth code | 
 **state** | **String** | CSRF state | 

### Return type

**String**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/html, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authLogin**
```swift
    open class func authLogin(loginRequest: LoginRequest, completion: @escaping (_ data: AuthResponse?, _ error: Error?) -> Void)
```

Login

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK

let loginRequest = LoginRequest(email: "email_example", password: "password_example") // LoginRequest | 

// Login
AuthAPI.authLogin(loginRequest: loginRequest) { (response, error) in
    guard error == nil else {
        print(error)
        return
    }

    if (response) {
        dump(response)
    }
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginRequest** | [**LoginRequest**](LoginRequest.md) |  | 

### Return type

[**AuthResponse**](AuthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authSignUp**
```swift
    open class func authSignUp(signUpRequest: SignUpRequest, completion: @escaping (_ data: User?, _ error: Error?) -> Void)
```

Register a new user

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK

let signUpRequest = SignUpRequest(email: "email_example", password: "password_example", username: "username_example", firstName: "firstName_example", lastName: "lastName_example", fiscalCode: "fiscalCode_example") // SignUpRequest | 

// Register a new user
AuthAPI.authSignUp(signUpRequest: signUpRequest) { (response, error) in
    guard error == nil else {
        print(error)
        return
    }

    if (response) {
        dump(response)
    }
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **signUpRequest** | [**SignUpRequest**](SignUpRequest.md) |  | 

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

