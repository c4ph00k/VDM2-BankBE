# AccountsAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**accountsCreateMovement**](AccountsAPI.md#accountscreatemovement) | **POST** /api/v1/accounts/movements | Create account movement
[**accountsGetBalance**](AccountsAPI.md#accountsgetbalance) | **GET** /api/v1/accounts/balance | Get account balance
[**accountsListMovements**](AccountsAPI.md#accountslistmovements) | **GET** /api/v1/accounts/movements | List account movements (paginated)


# **accountsCreateMovement**
```swift
    open class func accountsCreateMovement(createMovementRequest: CreateMovementRequest, completion: @escaping (_ data: Movement?, _ error: Error?) -> Void)
```

Create account movement

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK

let createMovementRequest = CreateMovementRequest(amount: "amount_example", type: "type_example", description: "description_example") // CreateMovementRequest | 

// Create account movement
AccountsAPI.accountsCreateMovement(createMovementRequest: createMovementRequest) { (response, error) in
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
 **createMovementRequest** | [**CreateMovementRequest**](CreateMovementRequest.md) |  | 

### Return type

[**Movement**](Movement.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **accountsGetBalance**
```swift
    open class func accountsGetBalance(completion: @escaping (_ data: BalanceResponse?, _ error: Error?) -> Void)
```

Get account balance

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK


// Get account balance
AccountsAPI.accountsGetBalance() { (response, error) in
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

[**BalanceResponse**](BalanceResponse.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **accountsListMovements**
```swift
    open class func accountsListMovements(page: Int? = nil, limit: Int? = nil, completion: @escaping (_ data: PaginatedMovementsResponse?, _ error: Error?) -> Void)
```

List account movements (paginated)

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK

let page = 987 // Int | Page number (default: 1) (optional) (default to 1)
let limit = 987 // Int | Items per page (default: 10, max: 100) (optional) (default to 10)

// List account movements (paginated)
AccountsAPI.accountsListMovements(page: page, limit: limit) { (response, error) in
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
 **page** | **Int** | Page number (default: 1) | [optional] [default to 1]
 **limit** | **Int** | Items per page (default: 10, max: 100) | [optional] [default to 10]

### Return type

[**PaginatedMovementsResponse**](PaginatedMovementsResponse.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

