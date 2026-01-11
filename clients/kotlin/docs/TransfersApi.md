# TransfersApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**transfersCreate**](TransfersApi.md#transfersCreate) | **POST** /api/v1/transfers | Create a transfer |
| [**transfersList**](TransfersApi.md#transfersList) | **GET** /api/v1/transfers | List transfers (paginated) |


<a id="transfersCreate"></a>
# **transfersCreate**
> Transfer transfersCreate(transferRequest)

Create a transfer

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = TransfersApi()
val transferRequest : TransferRequest =  // TransferRequest | 
try {
    val result : Transfer = apiInstance.transfersCreate(transferRequest)
    println(result)
} catch (e: ClientException) {
    println("4xx response calling TransfersApi#transfersCreate")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling TransfersApi#transfersCreate")
    e.printStackTrace()
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **transferRequest** | [**TransferRequest**](TransferRequest.md)|  | |

### Return type

[**Transfer**](Transfer.md)

### Authorization


Configure BearerPASETO:
    ApiClient.accessToken = ""
Configure BearerJWT:
    ApiClient.accessToken = ""

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a id="transfersList"></a>
# **transfersList**
> PaginatedTransfersResponse transfersList(page, limit)

List transfers (paginated)

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = TransfersApi()
val page : kotlin.Int = 56 // kotlin.Int | Page number (default: 1)
val limit : kotlin.Int = 56 // kotlin.Int | Items per page (default: 10, max: 100)
try {
    val result : PaginatedTransfersResponse = apiInstance.transfersList(page, limit)
    println(result)
} catch (e: ClientException) {
    println("4xx response calling TransfersApi#transfersList")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling TransfersApi#transfersList")
    e.printStackTrace()
}
```

### Parameters
| **page** | **kotlin.Int**| Page number (default: 1) | [optional] [default to 1] |
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **limit** | **kotlin.Int**| Items per page (default: 10, max: 100) | [optional] [default to 10] |

### Return type

[**PaginatedTransfersResponse**](PaginatedTransfersResponse.md)

### Authorization


Configure BearerPASETO:
    ApiClient.accessToken = ""
Configure BearerJWT:
    ApiClient.accessToken = ""

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

