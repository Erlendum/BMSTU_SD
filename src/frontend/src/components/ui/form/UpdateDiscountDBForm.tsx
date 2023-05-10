import {
	Button,
	Checkbox,
	createStyles,
	FormControlLabel,
	makeStyles,
	Paper,
	Typography
} from '@material-ui/core'
import React, { useState } from 'react'
import CustomTextField from './CustomTextField'
import MessageBox from '../message-box/MessageBox'
import { IDiscount } from '../../../types/discount.interface'
import { DiscountService } from '../../../services/discount.service'

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

// @ts-ignore
const UpdateDiscountDBForm = ({ id }) => {
	const [error, setError] = useState('')
	const [isForAll, setIsForAll] = useState(false)
	const classes = useStyles()
	const [values, setValues] = useState<IDiscount>({
		DiscountId: id,
		InstrumentId: 0,
		UserId: 0,
		Amount: 0,
		Type: '',
		DateBegin: '',
		DateEnd: ''
	})
	const [errors, setErrors] = useState({
		instrumentIdNotNumber: '',
		instrumentIdNotPositiveNumber: '',
		userIdNotNumber: '',
		userIdNotPositiveNumber: '',
		amountNotNumber: '',
		amountNotPositiveNumber: '',
		invalidBeginDate: '',
		invalidEndDate: ''
	})

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		validateFields()

		if (
			event.target.name == 'InstrumentId' ||
			event.target.name == 'UserId' ||
			event.target.name == 'Amount'
		)
			setValues({
				...values,
				[event.target.name]: parseInt(event.target.value)
			})
		else setValues({ ...values, [event.target.name]: event.target.value })
	}

	const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault()
		if (validateFields()) {
			return
		}
		setError('no error')
		let oldBeginDate = values.DateBegin
		let oldEndDate = values.DateEnd
		if (values.DateBegin == '') values.DateBegin = '2006-01-02'
		values.DateBegin += 'T15:04:05Z'
		if (values.DateEnd == '') values.DateEnd = '2006-01-02'
		values.DateEnd += 'T15:04:05Z'
		DiscountService.update(values).catch(error => {
			if (error.response) {
				setError(error.response.data.Error)
			}
		})

		values.DateBegin = oldBeginDate
		values.DateEnd = oldEndDate
		console.log(values)
		console.log(error)
	}

	const validateFields = () => {
		let error = false
		if (isNaN(values.InstrumentId)) {
			error = true
			console.log('InstrumentId not a number')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: 'InstrumentId should be integer number',
				instrumentIdNotPositiveNumber: '',
				userIdNotNumber: '',
				userIdNotPositiveNumber: '',
				amountNotNumber: '',
				amountNotPositiveNumber: '',
				invalidBeginDate: '',
				invalidEndDate: ''
			}))
			console.log('instrumentIdNotNumber set')
			return error
		}
		if (values.InstrumentId < 0) {
			error = true
			console.log('InstrumentId a positive number')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: '',
				instrumentIdNotPositiveNumber: 'InstrumentId should be positive number',
				userIdNotNumber: '',
				userIdNotPositiveNumber: '',
				amountNotNumber: '',
				amountNotPositiveNumber: '',
				invalidBeginDate: '',
				invalidEndDate: ''
			}))
			console.log('instrumentIdNotPositiveNumber set')
			return error
		}

		if (isNaN(values.UserId)) {
			error = true
			console.log('UserId not a number')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: '',
				instrumentIdNotPositiveNumber: '',
				userIdNotNumber: 'UserId should be integer number',
				userIdNotPositiveNumber: '',
				amountNotNumber: '',
				amountNotPositiveNumber: '',
				invalidBeginDate: '',
				invalidEndDate: ''
			}))
			console.log('userIdNotNumber set')
			return error
		}
		if (values.UserId < 0) {
			error = true
			console.log('UserId not a positive number')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: '',
				instrumentIdNotPositiveNumber: '',
				userIdNotNumber: '',
				userIdNotPositiveNumber: 'userId should be positive number',
				amountNotNumber: '',
				amountNotPositiveNumber: '',
				invalidBeginDate: '',
				invalidEndDate: ''
			}))
			console.log('userIdNotPositiveNumber set')
			return error
		}

		if (isNaN(values.Amount)) {
			error = true
			console.log('Amount not a number')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: '',
				instrumentIdNotPositiveNumber: '',
				userIdNotNumber: '',
				userIdNotPositiveNumber: '',
				amountNotNumber: 'Amount should be integer number',
				amountNotPositiveNumber: '',
				invalidBeginDate: '',
				invalidEndDate: ''
			}))
			console.log('amountNotNumber set')
			return error
		}
		if (values.Amount < 0) {
			error = true
			console.log('Amount not a positive number')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: '',
				instrumentIdNotPositiveNumber: '',
				userIdNotNumber: '',
				userIdNotPositiveNumber: '',
				amountNotNumber: '',
				amountNotPositiveNumber: 'Amount should be positive number',
				invalidBeginDate: '',
				invalidEndDate: ''
			}))
			console.log('amountNotPositiveNumber set')
			return error
		}

		if (
			values.DateBegin != '' &&
			isNaN(Date.parse(values.DateBegin.replace(/-/g, '/')))
		) {
			error = true
			console.log('invalid dateBegin')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: '',
				instrumentIdNotPositiveNumber: '',
				userIdNotNumber: '',
				userIdNotPositiveNumber: '',
				amountNotNumber: '',
				amountNotPositiveNumber: '',
				invalidBeginDate: 'Invalid Begin Date',
				invalidEndDate: ''
			}))
			console.log('invalidBeginDate set')
			return error
		}

		if (
			values.DateEnd != '' &&
			isNaN(Date.parse(values.DateEnd.replace(/-/g, '/')))
		) {
			error = true
			console.log('invalid dateEnd')

			setErrors(state => ({
				...state,
				instrumentIdNotNumber: '',
				instrumentIdNotPositiveNumber: '',
				userIdNotNumber: '',
				userIdNotPositiveNumber: '',
				amountNotNumber: '',
				amountNotPositiveNumber: '',
				invalidBeginDate: '',
				invalidEndDate: 'Invalid End Date'
			}))
			console.log('invalidEndDate set')
			return error
		}

		setErrors(state => ({
			...state,
			instrumentIdNotNumber: '',
			instrumentIdNotPositiveNumber: '',
			userIdNotNumber: '',
			userIdNotPositiveNumber: '',
			amountNotNumber: '',
			amountNotPositiveNumber: '',
			invalidBeginDate: '',
			invalidEndDate: ''
		}))
		console.log('no errors')
		return error
	}

	return (
		<Paper className={classes.container}>
			<Typography variant={'h4'} className={classes.title}>
				Update discount in database
			</Typography>
			<form onSubmit={e => handleSubmit(e)} className={classes.form}>
				<CustomTextField
					changeHandler={handleChange}
					label={'Instrument ID'}
					error={
						Boolean(errors.instrumentIdNotNumber) ||
						Boolean(errors.instrumentIdNotPositiveNumber)
					}
					helperText={
						errors.instrumentIdNotNumber == ''
							? errors.instrumentIdNotPositiveNumber
							: errors.instrumentIdNotNumber
					}
					name={'InstrumentId'}
					hide={false}
				/>
				<CustomTextField
					changeHandler={handleChange}
					label={'User ID'}
					error={
						Boolean(errors.userIdNotNumber) ||
						Boolean(errors.userIdNotPositiveNumber)
					}
					helperText={
						errors.userIdNotNumber == ''
							? errors.userIdNotPositiveNumber
							: errors.userIdNotNumber
					}
					name={'UserId'}
					hide={false}
				/>
				<CustomTextField
					changeHandler={handleChange}
					label={'Amount'}
					error={
						Boolean(errors.amountNotNumber) ||
						Boolean(errors.amountNotPositiveNumber)
					}
					helperText={
						errors.amountNotNumber == ''
							? errors.amountNotPositiveNumber
							: errors.amountNotNumber
					}
					name={'Amount'}
					hide={false}
				/>
				<CustomTextField
					label={'Type'}
					name={'Type'}
					error={false}
					helperText={''}
					changeHandler={handleChange}
					hide={false}
				/>
				<CustomTextField
					label={'Date Begin'}
					name={'DateBegin'}
					error={Boolean(errors.invalidBeginDate)}
					helperText={errors.invalidBeginDate}
					changeHandler={handleChange}
					hide={false}
				/>
				<CustomTextField
					label={'Date End'}
					name={'DateEnd'}
					error={Boolean(errors.invalidEndDate)}
					helperText={errors.invalidEndDate}
					changeHandler={handleChange}
					hide={false}
				/>

				<Button
					type={'submit'}
					variant={'contained'}
					className={classes.button}
				>
					Update
				</Button>
				<MessageBox type={error == 'no error' ? 'success' : 'error'}>
					{error == 'no error' ? 'success' : error}
				</MessageBox>
			</form>
		</Paper>
	)
}

export default UpdateDiscountDBForm
