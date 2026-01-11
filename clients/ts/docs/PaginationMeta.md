
# PaginationMeta


## Properties

Name | Type
------------ | -------------
`currentPage` | number
`totalPages` | number
`totalItems` | number
`perPage` | number

## Example

```typescript
import type { PaginationMeta } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "currentPage": null,
  "totalPages": null,
  "totalItems": null,
  "perPage": null,
} satisfies PaginationMeta

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PaginationMeta
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


