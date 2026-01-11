
# CreateMovementRequest


## Properties

Name | Type
------------ | -------------
`amount` | string
`type` | string
`description` | string

## Example

```typescript
import type { CreateMovementRequest } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "amount": 12.34,
  "type": null,
  "description": null,
} satisfies CreateMovementRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateMovementRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


