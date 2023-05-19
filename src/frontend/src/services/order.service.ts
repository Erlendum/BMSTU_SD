import axios from 'axios'
import { UserService } from './user.service'
import { IInstrumentsResponse } from '../types/instrument.interface'
import { IOrder, IOrderElementsResponse, IOrdersResponse } from '../types/order.interface'
import { IDiscount } from '../types/discount.interface'

export const OrderService = {
	async create() {
		const res = await axios.post('/create_order', '', {
			params: {
				id: UserService.getCurrentUserId()
			}
		})
		return res.data
	},
	async getList() {
		const response = await axios.get<IOrdersResponse>('/orders', {
			params: {user_id: UserService.getCurrentUserId()}
		})
		return response.data
	},
	async getListForAll() {
		const response = await axios.get<IOrdersResponse>('/ordersForAll', {
			params: {}
		})
		return response.data
	},
	async update(order: IOrder) {
		const json = JSON.stringify(order)
		const res = await axios.post('/update_order', json, {
			params: { login: UserService.getCurrentLogin() }
		})
		return res.data
	},

	async getOrderElements(id: number) {
		const response = await axios.get<IOrderElementsResponse>('/order_elements', {
			params: {order_id: id}
		})
		return response.data
	}
}
