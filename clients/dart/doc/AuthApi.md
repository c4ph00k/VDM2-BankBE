# vdm2_bank_sdk.api.AuthApi

## Load the API package
```dart
import 'package:vdm2_bank_sdk/api.dart';
```

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**authGoogle**](AuthApi.md#authgoogle) | **GET** /api/v1/auth/google | Start Google OAuth flow
[**authGoogleCallback**](AuthApi.md#authgooglecallback) | **GET** /api/v1/auth/google/callback | Handle Google OAuth callback
[**authLogin**](AuthApi.md#authlogin) | **POST** /api/v1/auth/login | Login
[**authSignUp**](AuthApi.md#authsignup) | **POST** /api/v1/auth/signup | Register a new user


# **authGoogle**
> authGoogle()

Start Google OAuth flow

Redirects the user-agent to Google and sets `oauth_state` cookie.

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';

final api_instance = AuthApi();

try {
    api_instance.authGoogle();
} catch (e) {
    print('Exception when calling AuthApi->authGoogle: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authGoogleCallback**
> String authGoogleCallback(code, state)

Handle Google OAuth callback

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';

final api_instance = AuthApi();
final code = code_example; // String | OAuth code
final state = state_example; // String | CSRF state

try {
    final result = api_instance.authGoogleCallback(code, state);
    print(result);
} catch (e) {
    print('Exception when calling AuthApi->authGoogleCallback: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **code** | **String**| OAuth code | 
 **state** | **String**| CSRF state | 

### Return type

**String**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/html, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authLogin**
> AuthResponse authLogin(loginRequest)

Login

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';

final api_instance = AuthApi();
final loginRequest = LoginRequest(); // LoginRequest | 

try {
    final result = api_instance.authLogin(loginRequest);
    print(result);
} catch (e) {
    print('Exception when calling AuthApi->authLogin: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginRequest** | [**LoginRequest**](LoginRequest.md)|  | 

### Return type

[**AuthResponse**](AuthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authSignUp**
> User authSignUp(signUpRequest)

Register a new user

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';

final api_instance = AuthApi();
final signUpRequest = SignUpRequest(); // SignUpRequest | 

try {
    final result = api_instance.authSignUp(signUpRequest);
    print(result);
} catch (e) {
    print('Exception when calling AuthApi->authSignUp: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **signUpRequest** | [**SignUpRequest**](SignUpRequest.md)|  | 

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

