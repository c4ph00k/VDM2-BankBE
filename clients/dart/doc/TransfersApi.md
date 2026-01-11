# vdm2_bank_sdk.api.TransfersApi

## Load the API package
```dart
import 'package:vdm2_bank_sdk/api.dart';
```

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**transfersCreate**](TransfersApi.md#transferscreate) | **POST** /api/v1/transfers | Create a transfer
[**transfersList**](TransfersApi.md#transferslist) | **GET** /api/v1/transfers | List transfers (paginated)


# **transfersCreate**
> Transfer transfersCreate(transferRequest)

Create a transfer

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

final api_instance = TransfersApi();
final transferRequest = TransferRequest(); // TransferRequest | 

try {
    final result = api_instance.transfersCreate(transferRequest);
    print(result);
} catch (e) {
    print('Exception when calling TransfersApi->transfersCreate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **transferRequest** | [**TransferRequest**](TransferRequest.md)|  | 

### Return type

[**Transfer**](Transfer.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **transfersList**
> PaginatedTransfersResponse transfersList(page, limit)

List transfers (paginated)

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

final api_instance = TransfersApi();
final page = 56; // int | Page number (default: 1)
final limit = 56; // int | Items per page (default: 10, max: 100)

try {
    final result = api_instance.transfersList(page, limit);
    print(result);
} catch (e) {
    print('Exception when calling TransfersApi->transfersList: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int**| Page number (default: 1) | [optional] [default to 1]
 **limit** | **int**| Items per page (default: 10, max: 100) | [optional] [default to 10]

### Return type

[**PaginatedTransfersResponse**](PaginatedTransfersResponse.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

