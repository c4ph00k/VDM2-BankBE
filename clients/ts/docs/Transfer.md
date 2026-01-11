
# Transfer


## Properties

Name | Type
------------ | -------------
`id` | number
`fromAccount` | string
`toAccount` | string
`amount` | string
`status` | string
`initiatedAt` | Date
`completedAt` | Date

## Example

```typescript
import type { Transfer } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "id": 1,
  "fromAccount": null,
  "toAccount": null,
  "amount": 12.34,
  "status": null,
  "initiatedAt": null,
  "completedAt": null,
} satisfies Transfer

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Transfer
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


