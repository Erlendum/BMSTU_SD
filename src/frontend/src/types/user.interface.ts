export interface IUserResponse {
	user: IUser
}

export interface IUser {
	UserId: number
	Login: string
	Password: string
	Fio: string
	DateBirth: string
	Gender: string
	IsAdmin: boolean
}
