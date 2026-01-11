# vdm2_bank_sdk.api.AccountsApi

## Load the API package
```dart
import 'package:vdm2_bank_sdk/api.dart';
```

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**accountsCreateMovement**](AccountsApi.md#accountscreatemovement) | **POST** /api/v1/accounts/movements | Create account movement
[**accountsGetBalance**](AccountsApi.md#accountsgetbalance) | **GET** /api/v1/accounts/balance | Get account balance
[**accountsListMovements**](AccountsApi.md#accountslistmovements) | **GET** /api/v1/accounts/movements | List account movements (paginated)


# **accountsCreateMovement**
> Movement accountsCreateMovement(createMovementRequest)

Create account movement

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';
// TODO Configure HTTP Bearer authorization: BearerPASETO
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerPASETO').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerPASETO').setAccessToken(yourTokenGeneratorFunction);
// TODO Configure HTTP Bearer authorization: BearerJWT
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerJWT').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerJWT').setAccessToken(yourTokenGeneratorFunction);

final api_instance = AccountsApi();
final createMovementRequest = CreateMovementRequest(); // CreateMovementRequest | 

try {
    final result = api_instance.accountsCreateMovement(createMovementRequest);
    print(result);
} catch (e) {
    print('Exception when calling AccountsApi->accountsCreateMovement: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createMovementRequest** | [**CreateMovementRequest**](CreateMovementRequest.md)|  | 

### Return type

[**Movement**](Movement.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **accountsGetBalance**
> BalanceResponse accountsGetBalance()

Get account balance

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';
// TODO Configure HTTP Bearer authorization: BearerPASETO
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerPASETO').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerPASETO').setAccessToken(yourTokenGeneratorFunction);
// TODO Configure HTTP Bearer authorization: BearerJWT
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerJWT').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerJWT').setAccessToken(yourTokenGeneratorFunction);

final api_instance = AccountsApi();

try {
    final result = api_instance.accountsGetBalance();
    print(result);
} catch (e) {
    print('Exception when calling AccountsApi->accountsGetBalance: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BalanceResponse**](BalanceResponse.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **accountsListMovements**
> PaginatedMovementsResponse accountsListMovements(page, limit)

List account movements (paginated)

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';
// TODO Configure HTTP Bearer authorization: BearerPASETO
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerPASETO').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerPASETO').setAccessToken(yourTokenGeneratorFunction);
// TODO Configure HTTP Bearer authorization: BearerJWT
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerJWT').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('BearerJWT').setAccessToken(yourTokenGeneratorFunction);

final api_instance = AccountsApi();
final page = 56; // int | Page number (default: 1)
final limit = 56; // int | Items per page (default: 10, max: 100)

try {
    final result = api_instance.accountsListMovements(page, limit);
    print(result);
} catch (e) {
    print('Exception when calling AccountsApi->accountsListMovements: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int**| Page number (default: 1) | [optional] [default to 1]
 **limit** | **int**| Items per page (default: 10, max: 100) | [optional] [default to 10]

### Return type

[**PaginatedMovementsResponse**](PaginatedMovementsResponse.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

