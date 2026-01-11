
# Movement

Mirrors `internal/model.Movement` JSON. NOTE: in Go it serializes `amount` as decimal (shopspring/decimal) which is typically a JSON string/number depending on config. TODO: confirm runtime JSON encoding for decimal.Decimal and adjust if needed. 

## Properties

Name | Type
------------ | -------------
`id` | number
`accountId` | string
`amount` | string
`type` | string
`description` | string
`occurredAt` | Date

## Example

```typescript
import type { Movement } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "id": 1,
  "accountId": null,
  "amount": 12.34,
  "type": null,
  "description": null,
  "occurredAt": null,
} satisfies Movement

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Movement
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


