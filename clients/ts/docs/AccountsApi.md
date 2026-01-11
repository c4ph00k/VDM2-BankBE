# AccountsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**accountsCreateMovement**](AccountsApi.md#accountscreatemovement) | **POST** /api/v1/accounts/movements | Create account movement |
| [**accountsGetBalance**](AccountsApi.md#accountsgetbalance) | **GET** /api/v1/accounts/balance | Get account balance |
| [**accountsListMovements**](AccountsApi.md#accountslistmovements) | **GET** /api/v1/accounts/movements | List account movements (paginated) |



## accountsCreateMovement

> Movement accountsCreateMovement(createMovementRequest)

Create account movement

### Example

```ts
import {
  Configuration,
  AccountsApi,
} from '@vdm2-bank/sdk';
import type { AccountsCreateMovementRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerPASETO
    accessToken: "YOUR BEARER TOKEN",
    // Configure HTTP bearer authorization: BearerJWT
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new AccountsApi(config);

  const body = {
    // CreateMovementRequest
    createMovementRequest: ...,
  } satisfies AccountsCreateMovementRequest;

  try {
    const data = await api.accountsCreateMovement(body);
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
| **createMovementRequest** | [CreateMovementRequest](CreateMovementRequest.md) |  | |

### Return type

[**Movement**](Movement.md)

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


## accountsGetBalance

> BalanceResponse accountsGetBalance()

Get account balance

### Example

```ts
import {
  Configuration,
  AccountsApi,
} from '@vdm2-bank/sdk';
import type { AccountsGetBalanceRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerPASETO
    accessToken: "YOUR BEARER TOKEN",
    // Configure HTTP bearer authorization: BearerJWT
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new AccountsApi(config);

  try {
    const data = await api.accountsGetBalance();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**BalanceResponse**](BalanceResponse.md)

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


## accountsListMovements

> PaginatedMovementsResponse accountsListMovements(page, limit)

List account movements (paginated)

### Example

```ts
import {
  Configuration,
  AccountsApi,
} from '@vdm2-bank/sdk';
import type { AccountsListMovementsRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerPASETO
    accessToken: "YOUR BEARER TOKEN",
    // Configure HTTP bearer authorization: BearerJWT
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new AccountsApi(config);

  const body = {
    // number | Page number (default: 1) (optional)
    page: 56,
    // number | Items per page (default: 10, max: 100) (optional)
    limit: 56,
  } satisfies AccountsListMovementsRequest;

  try {
    const data = await api.accountsListMovements(body);
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

[**PaginatedMovementsResponse**](PaginatedMovementsResponse.md)

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

