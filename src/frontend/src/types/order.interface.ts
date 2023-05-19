export interface IOrdersResponse {
	orders: IOrder[]
}

export interface IOrder {
	OrderId: number
	Time: string
	Price: number
	Status: string
	UserId: number
}

export interface IOrderElementsResponse {
	order_elements: IOrderElement[]
}

export interface IOrderElement {
	OrderElementId: number
	InstrumentId:   number
	OrderId:        number
	Amount:         number
	Price:          number
}
