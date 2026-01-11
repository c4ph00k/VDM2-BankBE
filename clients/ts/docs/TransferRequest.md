
# TransferRequest


## Properties

Name | Type
------------ | -------------
`toAccount` | string
`amount` | string
`description` | string

## Example

```typescript
import type { TransferRequest } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "toAccount": null,
  "amount": 12.34,
  "description": null,
} satisfies TransferRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as TransferRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


