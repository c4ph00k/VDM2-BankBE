
# ErrorResponse

Current error envelope from `internal/util/errors.go`. Note: for non-*util.APIError errors, the server responds with code=500 and message=\"internal server error\". 

## Properties

Name | Type
------------ | -------------
`error` | [APIError](APIError.md)

## Example

```typescript
import type { ErrorResponse } from '@vdm2-bank/sdk'

// TODO: Update the object below with actual values
const example = {
  "error": null,
} satisfies ErrorResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ErrorResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


