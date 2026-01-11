# MetaApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**healthCheck**](MetaApi.md#healthcheck) | **GET** /health | Health check |
| [**metrics**](MetaApi.md#metrics) | **GET** /metrics | Prometheus metrics |



## healthCheck

> healthCheck()

Health check

### Example

```ts
import {
  Configuration,
  MetaApi,
} from '@vdm2-bank/sdk';
import type { HealthCheckRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const api = new MetaApi();

  try {
    const data = await api.healthCheck();
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

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## metrics

> string metrics()

Prometheus metrics

### Example

```ts
import {
  Configuration,
  MetaApi,
} from '@vdm2-bank/sdk';
import type { MetricsRequest } from '@vdm2-bank/sdk';

async function example() {
  console.log("🚀 Testing @vdm2-bank/sdk SDK...");
  const api = new MetaApi();

  try {
    const data = await api.metrics();
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

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `text/plain`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

