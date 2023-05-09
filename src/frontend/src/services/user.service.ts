import axios from 'axios'
import { IUser, IUserResponse } from '../types/user.interface'
import { IInstrumentsResponse } from '../types/instrument.interface'
import { ICart, ICartResponse } from '../types/cart.interface'

axios.defaults.baseURL = 'http://localhost:8080'

export const UserService = {
	async create(user: IUser) {
		const json = JSON.stringify(user)
		const res = await axios.post('/create_user', json)
		return res.data
	},

	async getCart() {
		const response = await axios.get<ICartResponse>('/cart', {
			params: { id: this.getCurrentUserId() }
		})
		if (response.data.cart != null) {
			localStorage.setItem('cartId', response.data.cart.CartId.toString())
		}
		return response.data
	},
	async get(login: string, password: string) {
		const res = await axios.get<IUserResponse>('/get_user', {
			params: {
				login: login,
				password: password
			}
		})
		if (res.data.user.Login == login) {
			localStorage.setItem('userId', res.data.user.UserId.toString())
			localStorage.setItem('user', res.data.user.Login)
			localStorage.setItem('isAdmin', res.data.user.IsAdmin ? 'true' : 'false')
			const link = document.createElement('a')
			link.href = '/cart'
			document.body.appendChild(link)
			link.click()
		}
		return res.data
	},
	getCurrentLogin() {
		return localStorage.getItem('user')
	},
	getCurrentUserId() {
		return localStorage.getItem('userId')
	},
	getCurrentIsAdmin() {
		return localStorage.getItem('isAdmin')
	},
	getCartId() {
		return localStorage.getItem('cartId')
	},
	logout() {
		localStorage.removeItem('user')
		localStorage.removeItem('userId')
		localStorage.removeItem('isAdmin')
		localStorage.removeItem('cartId')
		const link = document.createElement('a')
		link.href = '/'
		document.body.appendChild(link)
		link.click()
	}
}
