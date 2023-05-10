export interface IDiscountsResponse {
	discounts: IDiscount[]
	limit: number
	skip: number
}

export interface IDiscount {
	DiscountId: number
	InstrumentId: number
	UserId: number
	Amount: number
	Type: string
	DateBegin: string
	DateEnd: string
}
