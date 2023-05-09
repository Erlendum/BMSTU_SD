import axios from 'axios'
import {
	IInstrument,
	IInstrumentsResponse
} from '../types/instrument.interface'
import { UserService } from './user.service'

axios.defaults.baseURL = 'http://localhost:8080'

export const InstrumentService = {
	async getList() {
		const response = await axios.get<IInstrumentsResponse>('/instruments', {
			params: {}
		})
		console.log(response.data.instruments.length)
		return response.data
	},

	async create(instrument: IInstrument) {
		const json = JSON.stringify(instrument)
		const res = await axios.post('/create_instrument', json, {
			params: { login: UserService.getCurrentLogin() }
		})
		return res.data
	},

	async delete(id: number) {
		const res = await axios.post('/delete_instrument', '', {
			params: { id: id, login: UserService.getCurrentLogin() }
		})
		return res.data
	},
	async update(instrument: IInstrument) {
		const json = JSON.stringify(instrument)
		const res = await axios.post('/update_instrument', json, {
			params: { login: UserService.getCurrentLogin() }
		})
		return res.data
	}
}
