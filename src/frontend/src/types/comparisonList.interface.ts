import { IInstrument } from './instrument.interface'

export interface IComparisonListResponse {
	comparisonList: IComparisonList
	comparisonListInstruments: IInstrument[]
}

export interface IComparisonList {
	ComparisonListId: number
	UserId: number
	TotalPrice: number
	Amount: number
}
