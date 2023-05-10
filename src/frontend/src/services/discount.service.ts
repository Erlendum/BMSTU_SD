import axios from 'axios'
import { IDiscount, IDiscountsResponse } from '../types/discount.interface'
import { UserService } from './user.service'

axios.defaults.baseURL = 'http://localhost:8080'

export const DiscountService = {
	async getList() {
		const response = await axios.get<IDiscountsResponse>('/discounts', {
			params: {}
		})
		console.log('discounts: ' + response.data.discounts.length)
		return response.data
	},

	async create(discount: IDiscount) {
		const json = JSON.stringify(discount)
		const res = await axios.post('/create_discount', json, {
			params: { login: UserService.getCurrentLogin() }
		})
		return res.data
	},

	async createForAll(discount: IDiscount) {
		const json = JSON.stringify(discount)
		const res = await axios.post('/create_for_all_discount', json, {
			params: { login: UserService.getCurrentLogin() }
		})
		return res.data
	},

	async delete(id: number) {
		const res = await axios.post('/delete_discount', '', {
			params: { id: id, login: UserService.getCurrentLogin() }
		})
		return res.data
	},
	async update(discount: IDiscount) {
		const json = JSON.stringify(discount)
		const res = await axios.post('/update_discount', json, {
			params: { login: UserService.getCurrentLogin() }
		})
		return res.data
	}
}
