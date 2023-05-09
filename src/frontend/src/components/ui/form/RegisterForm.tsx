import React, { useState } from 'react'
import { IInstrument } from '../../../types/instrument.interface'
import { InstrumentService } from '../../../services/instrument.service'
import {
	Button,
	createStyles,
	makeStyles,
	Paper,
	Typography
} from '@material-ui/core'
import CustomTextField from './CustomTextField'
import MessageBox from '../message-box/MessageBox'
import { IUser } from '../../../types/user.interface'
import { UserService } from '../../../services/user.service'
import Layout from '../layout/Layout'

const useStyles = makeStyles(() =>
	createStyles({
		form: {
			display: 'flex',
			flexDirection: 'column'
		},
		container: {
			backgroundColor: '#ffffff',
			position: 'absolute',
			top: '50%',
			left: '50%',
			transform: 'translate(-50%,-50%)',
			padding: 30,
			textAlign: 'center'
		},
		title: {
			margin: '0px 0 20px 0'
		},
		button: {
			margin: '20px 0'
		}
	})
)

const RegisterForm = () => {
	const [error, setError] = useState('')
	const classes = useStyles()
	const [values, setValues] = useState<IUser>({
		UserId: 0,
		Login: '',
		Password: '',
		Fio: '',
		DateBirth: '',
		Gender: '',
		IsAdmin: false
	})
	const [errors, setErrors] = useState({
		invalidDate: ''
	})

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		validateFields()

		setValues({ ...values, [event.target.name]: event.target.value })
	}

	const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault()
		if (validateFields()) {
			return
		}
		setError('no error')
		let oldDate = values.DateBirth
		values.DateBirth += 'T15:04:05Z'
		UserService.create(values).catch(error => {
			if (error.response) {
				setError(error.response.data.Error)
			}
		})
		console.log(values)
		values.DateBirth = oldDate
		console.log(error)
	}

	const validateFields = () => {
		let error = false
		if (isNaN(Date.parse(values.DateBirth))) {
			error = true
			console.log('invalid date')

			setErrors(state => ({
				...state,
				invalidDate: 'Invalid date'
			}))
			console.log('invalidDate set')
			return error
		}

		setErrors(state => ({
			...state,
			invalidDate: ''
		}))
		console.log('no errors')
		return error
	}

	return (
		<Layout>
			<Paper className={classes.container}>
				<Typography variant={'h4'} className={classes.title}>
					Registation
				</Typography>
				<form onSubmit={e => handleSubmit(e)} className={classes.form}>
					<CustomTextField
						changeHandler={handleChange}
						label={'Login'}
						error={false}
						helperText={''}
						name={'Login'}
						hide={false}
					/>
					<CustomTextField
						changeHandler={handleChange}
						label={'Password'}
						error={false}
						helperText={''}
						name={'Password'}
						hide={true}
					/>
					<CustomTextField
						changeHandler={handleChange}
						label={'Fio'}
						error={false}
						helperText={''}
						name={'Fio'}
						hide={false}
					/>
					<CustomTextField
						label={'DateBirth'}
						name={'DateBirth'}
						error={Boolean(errors.invalidDate)}
						helperText={errors.invalidDate}
						changeHandler={handleChange}
						hide={false}
					/>
					<CustomTextField
						label={'Gender'}
						name={'Gender'}
						error={false}
						helperText={''}
						changeHandler={handleChange}
						hide={false}
					/>
					<Button
						type={'submit'}
						variant={'contained'}
						className={classes.button}
					>
						Register
					</Button>
					<MessageBox type={error == 'no error' ? 'success' : 'error'}>
						{error == 'no error' ? 'success' : error}
					</MessageBox>
				</form>
			</Paper>
		</Layout>
	)
}

export default RegisterForm
