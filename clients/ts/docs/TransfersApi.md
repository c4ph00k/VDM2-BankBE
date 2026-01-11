# TransfersApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**transfersCreate**](TransfersApi.md#transferscreate) | **POST** /api/v1/transfers | Create a transfer |
| [**transfersList**](TransfersApi.md#transferslist) | **GET** /api/v1/transfers | List transfers (paginated) |



## transfersCreate

> Transfer transfersCreate(transferRequest)

Create a transfer

### Example

```ts
import {
  Configuration,
  TransfersApi,
} from '@vdm2-bank/sdk';
import type { TransfersCreateRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerPASETO
    accessToken: "YOUR BEARER TOKEN",
    // Configure HTTP bearer authorization: BearerJWT
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new TransfersApi(config);

  const body = {
    // TransferRequest
    transferRequest: ...,
  } satisfies TransfersCreateRequest;

  try {
    const data = await api.transfersCreate(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **transferRequest** | [TransferRequest](TransferRequest.md) |  | |

### Return type

[**Transfer**](Transfer.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | Bad request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## transfersList

> PaginatedTransfersResponse transfersList(page, limit)

List transfers (paginated)

### Example

```ts
import {
  Configuration,
  TransfersApi,
} from '@vdm2-bank/sdk';
import type { TransfersListRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerPASETO
    accessToken: "YOUR BEARER TOKEN",
    // Configure HTTP bearer authorization: BearerJWT
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new TransfersApi(config);

  const body = {
    // number | Page number (default: 1) (optional)
    page: 56,
    // number | Items per page (default: 10, max: 100) (optional)
    limit: 56,
  } satisfies TransfersListRequest;

  try {
    const data = await api.transfersList(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **page** | `number` | Page number (default: 1) | [Optional] [Defaults to `1`] |
| **limit** | `number` | Items per page (default: 10, max: 100) | [Optional] [Defaults to `10`] |

### Return type

[**PaginatedTransfersResponse**](PaginatedTransfersResponse.md)

### Authorization

[BearerPASETO](../README.md#BearerPASETO), [BearerJWT](../README.md#BearerJWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

