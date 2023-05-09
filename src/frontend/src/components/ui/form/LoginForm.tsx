import {
	Button,
	createStyles,
	makeStyles,
	Paper,
	Typography
} from '@material-ui/core'
import React, { useState } from 'react'
import { IUser } from '../../../types/user.interface'
import { UserService } from '../../../services/user.service'
import Layout from '../layout/Layout'
import CustomTextField from './CustomTextField'
import MessageBox from '../message-box/MessageBox'

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

const LoginForm = () => {
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
	const [errors, setErrors] = useState({})

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setValues({ ...values, [event.target.name]: event.target.value })
	}

	const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault()
		if (validateFields()) {
			return
		}
		setError('no error')
		UserService.get(values.Login, values.Password).catch(error => {
			if (error.response) {
				setError(error.response.data.Error)
			}
		})
		console.log(values)
		console.log(error)
	}

	const validateFields = () => {
		let error = false

		setErrors(state => ({
			...state
		}))
		console.log('no errors')
		return error
	}

	return (
		<Layout>
			<Paper className={classes.container}>
				<Typography variant={'h4'} className={classes.title}>
					Login
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
					<Button
						type={'submit'}
						variant={'contained'}
						className={classes.button}
					>
						Login
					</Button>
					<MessageBox type={error == 'no error' ? 'success' : 'error'}>
						{error == 'no error' ? 'success' : error}
					</MessageBox>
				</form>
			</Paper>
		</Layout>
	)
}

export default LoginForm
