import { IInstrument } from './instrument.interface'

export interface ICartResponse {
	cart: ICart
	cartInstruments: IInstrument[]
}

export interface ICart {
	CartId: number
	UserId: number
	TotalPrice: number
	Amount: number
}
