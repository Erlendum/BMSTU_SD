import axios from 'axios'
import { IUser, IUserResponse } from '../types/user.interface'
import { IComparisonListResponse } from '../types/comparisonList.interface'
import { UserService } from './user.service'

export const ComparisonListService = {
	async addInstrument(instrumentId: number) {
		const res = await axios.post('/add_instrument_to_comparison_list', '', {
			params: {
				id: UserService.getComparisonListId(),
				instrumentId: instrumentId
			}
		})
		return res.data
	},

	async deleteInstrument(instrumentId: number) {
		const res = await axios.post(
			'/delete_instrument_from_comparison_list',
			'',
			{
				params: {
					id: UserService.getComparisonListId(),
					instrumentId: instrumentId
				}
			}
		)
		return res.data
	}
}
