import axios from 'axios'
import { UserService } from './user.service'

export const OrderService = {
	async create() {
		const res = await axios.post('/create_order', '', {
			params: {
				id: UserService.getCurrentUserId()
			}
		})
		return res.data
	}
}
