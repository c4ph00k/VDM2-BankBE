# AccountsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**accountsCreateMovement**](AccountsApi.md#accountsCreateMovement) | **POST** /api/v1/accounts/movements | Create account movement |
| [**accountsGetBalance**](AccountsApi.md#accountsGetBalance) | **GET** /api/v1/accounts/balance | Get account balance |
| [**accountsListMovements**](AccountsApi.md#accountsListMovements) | **GET** /api/v1/accounts/movements | List account movements (paginated) |


<a id="accountsCreateMovement"></a>
# **accountsCreateMovement**
> Movement accountsCreateMovement(createMovementRequest)

Create account movement

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = AccountsApi()
val createMovementRequest : CreateMovementRequest =  // CreateMovementRequest | 
try {
    val result : Movement = apiInstance.accountsCreateMovement(createMovementRequest)
    println(result)
} catch (e: ClientException) {
    println("4xx response calling AccountsApi#accountsCreateMovement")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling AccountsApi#accountsCreateMovement")
    e.printStackTrace()
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **createMovementRequest** | [**CreateMovementRequest**](CreateMovementRequest.md)|  | |

### Return type

[**Movement**](Movement.md)

### Authorization


Configure BearerPASETO:
    ApiClient.accessToken = ""
Configure BearerJWT:
    ApiClient.accessToken = ""

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a id="accountsGetBalance"></a>
# **accountsGetBalance**
> BalanceResponse accountsGetBalance()

Get account balance

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = AccountsApi()
try {
    val result : BalanceResponse = apiInstance.accountsGetBalance()
    println(result)
} catch (e: ClientException) {
    println("4xx response calling AccountsApi#accountsGetBalance")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling AccountsApi#accountsGetBalance")
    e.printStackTrace()
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BalanceResponse**](BalanceResponse.md)

### Authorization


Configure BearerPASETO:
    ApiClient.accessToken = ""
Configure BearerJWT:
    ApiClient.accessToken = ""

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a id="accountsListMovements"></a>
# **accountsListMovements**
> PaginatedMovementsResponse accountsListMovements(page, limit)

List account movements (paginated)

### Example
```kotlin
// Import classes:
//import org.openapitools.client.infrastructure.*
//import org.openapitools.client.models.*

val apiInstance = AccountsApi()
val page : kotlin.Int = 56 // kotlin.Int | Page number (default: 1)
val limit : kotlin.Int = 56 // kotlin.Int | Items per page (default: 10, max: 100)
try {
    val result : PaginatedMovementsResponse = apiInstance.accountsListMovements(page, limit)
    println(result)
} catch (e: ClientException) {
    println("4xx response calling AccountsApi#accountsListMovements")
    e.printStackTrace()
} catch (e: ServerException) {
    println("5xx response calling AccountsApi#accountsListMovements")
    e.printStackTrace()
}
```

### Parameters
| **page** | **kotlin.Int**| Page number (default: 1) | [optional] [default to 1] |
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **limit** | **kotlin.Int**| Items per page (default: 10, max: 100) | [optional] [default to 10] |

### Return type

[**PaginatedMovementsResponse**](PaginatedMovementsResponse.md)

### Authorization


Configure BearerPASETO:
    ApiClient.accessToken = ""
Configure BearerJWT:
    ApiClient.accessToken = ""

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

