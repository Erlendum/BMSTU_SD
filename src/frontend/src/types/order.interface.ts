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
