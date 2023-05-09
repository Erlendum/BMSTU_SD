import axios from 'axios'
import { IUser, IUserResponse } from '../types/user.interface'
import { ICartResponse } from '../types/cart.interface'
import { UserService } from './user.service'

export const CartService = {
	async addInstrument(instrumentId: number) {
		const res = await axios.post('/add_instrument_to_cart', '', {
			params: {
				id: UserService.getCartId(),
				instrumentId: instrumentId
			}
		})
		return res.data
	},

	async deleteInstrument(instrumentId: number) {
		const res = await axios.post('/delete_instrument_from_cart', '', {
			params: {
				id: UserService.getCartId(),
				instrumentId: instrumentId
			}
		})
		return res.data
	}
}
