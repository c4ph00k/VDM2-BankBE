
# BalanceResponse


## Properties

Name | Type
------------ | -------------
`accountId` | string
`balance` | string
`currency` | string

## Example

```typescript
import type { BalanceResponse } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "accountId": null,
  "balance": 100.00,
  "currency": EUR,
} satisfies BalanceResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as BalanceResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


