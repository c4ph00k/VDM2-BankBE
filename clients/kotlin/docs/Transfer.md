
# Transfer

## Properties
| Name | Type | Description | Notes |
| ------------ | ------------- | ------------- | ------------- |
| **id** | **kotlin.Long** |  |  |
| **fromAccount** | [**java.util.UUID**](java.util.UUID.md) |  |  |
| **toAccount** | [**java.util.UUID**](java.util.UUID.md) |  |  |
| **amount** | **kotlin.String** | Decimal encoded as string (shopspring/decimal) |  |
| **status** | [**inline**](#Status) |  |  |
| **initiatedAt** | [**java.time.OffsetDateTime**](java.time.OffsetDateTime.md) |  |  |
| **completedAt** | [**java.time.OffsetDateTime**](java.time.OffsetDateTime.md) |  |  |


<a id="Status"></a>
## Enum: status
| Name | Value |
| ---- | ----- |
| status | pending, completed, failed |



