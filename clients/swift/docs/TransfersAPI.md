# TransfersAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**transfersCreate**](TransfersAPI.md#transferscreate) | **POST** /api/v1/transfers | Create a transfer
[**transfersList**](TransfersAPI.md#transferslist) | **GET** /api/v1/transfers | List transfers (paginated)


# **transfersCreate**
```swift
    open class func transfersCreate(transferRequest: TransferRequest, completion: @escaping (_ data: Transfer?, _ error: Error?) -> Void)
```

Create a transfer

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK

let transferRequest = TransferRequest(toAccount: 123, amount: "amount_example", description: "description_example") // TransferRequest | 

// Create a transfer
TransfersAPI.transfersCreate(transferRequest: transferRequest) { (response, error) in
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
 **transferRequest** | [**TransferRequest**](TransferRequest.md) |  | 

### Return type

[**Transfer**](Transfer.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **transfersList**
```swift
    open class func transfersList(page: Int? = nil, limit: Int? = nil, completion: @escaping (_ data: PaginatedTransfersResponse?, _ error: Error?) -> Void)
```

List transfers (paginated)

### Example
```swift
// The following code samples are still beta. For any issue, please report via http://github.com/OpenAPITools/openapi-generator/issues/new
import VDM2BankSDK

let page = 987 // Int | Page number (default: 1) (optional) (default to 1)
let limit = 987 // Int | Items per page (default: 10, max: 100) (optional) (default to 10)

// List transfers (paginated)
TransfersAPI.transfersList(page: page, limit: limit) { (response, error) in
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

[**PaginatedTransfersResponse**](PaginatedTransfersResponse.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

