# MetaApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**healthCheck**](MetaApi.md#healthCheck) | **GET** /health | Health check |
| [**metrics**](MetaApi.md#metrics) | **GET** /metrics | Prometheus metrics |


<a id="healthCheck"></a>
# **healthCheck**
> healthCheck()

Health check

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = MetaApi()
try {
    apiInstance.healthCheck()
} catch (e: ClientException) {
    println("4xx response calling MetaApi#healthCheck")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling MetaApi#healthCheck")
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
 - **Accept**: Not defined

<a id="metrics"></a>
# **metrics**
> kotlin.String metrics()

Prometheus metrics

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = MetaApi()
try {
    val result : kotlin.String = apiInstance.metrics()
    println(result)
} catch (e: ClientException) {
    println("4xx response calling MetaApi#metrics")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling MetaApi#metrics")
    e.printStackTrace()
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

**kotlin.String**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain

