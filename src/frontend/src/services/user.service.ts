import axios from 'axios'
import { IUser, IUserResponse } from '../types/user.interface'
import { IInstrumentsResponse } from '../types/instrument.interface'
import {
	IComparisonList,
	IComparisonListResponse
} from '../types/comparisonList.interface'

axios.defaults.baseURL = 'http://localhost:8080'

export const UserService = {
	async create(user: IUser) {
		const json = JSON.stringify(user)
		const res = await axios.post('/create_user', json)
		return res.data
	},

	async getComparisonList() {
		const response = await axios.get<IComparisonListResponse>(
			'/comparison_list',
			{
				params: { id: this.getCurrentUserId() }
			}
		)
		if (response.data.comparisonList != null) {
			localStorage.setItem(
				'comparisonListId',
				response.data.comparisonList.ComparisonListId.toString()
			)
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
			link.href = '/comparison_list'
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
	getComparisonListId() {
		return localStorage.getItem('comparisonListId')
	},
	logout() {
		localStorage.removeItem('user')
		localStorage.removeItem('userId')
		localStorage.removeItem('isAdmin')
		localStorage.removeItem('comparisonListId')
		const link = document.createElement('a')
		link.href = '/'
		document.body.appendChild(link)
		link.click()
	}
}
