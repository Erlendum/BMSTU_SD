export interface IInstrumentsResponse {
	instruments: IInstrument[]
	limit: number
	skip: number
}

export interface IInstrument {
	InstrumentId: number
	Brand: string
	Name: string
	Price: number
	Material: string
	Type: string
	Img: string
}
