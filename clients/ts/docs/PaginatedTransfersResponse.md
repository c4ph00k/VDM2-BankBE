
# PaginatedTransfersResponse

Concrete shape of `util.PaginatedResponse` as returned by `TransferService.GetByAccountID()`. 

## Properties

Name | Type
------------ | -------------
`data` | [Array&lt;Transfer&gt;](Transfer.md)
`pagination` | [PaginationMeta](PaginationMeta.md)

## Example

```typescript
import type { PaginatedTransfersResponse } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "data": null,
  "pagination": null,
} satisfies PaginatedTransfersResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PaginatedTransfersResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


