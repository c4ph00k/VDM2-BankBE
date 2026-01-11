# vdm2_bank_sdk.api.MetaApi

## Load the API package
```dart
import 'package:vdm2_bank_sdk/api.dart';
```

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**healthCheck**](MetaApi.md#healthcheck) | **GET** /health | Health check
[**metrics**](MetaApi.md#metrics) | **GET** /metrics | Prometheus metrics


# **healthCheck**
> healthCheck()

Health check

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';

final api_instance = MetaApi();

try {
    api_instance.healthCheck();
} catch (e) {
    print('Exception when calling MetaApi->healthCheck: $e\n');
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
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **metrics**
> String metrics()

Prometheus metrics

### Example
```dart
import 'package:vdm2_bank_sdk/api.dart';

final api_instance = MetaApi();

try {
    final result = api_instance.metrics();
    print(result);
} catch (e) {
    print('Exception when calling MetaApi->metrics: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

**String**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

