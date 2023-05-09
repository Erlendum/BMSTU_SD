import React, { useEffect, useState } from 'react'
import {
	createStyles,
	makeStyles,
	Typography,
	Paper,
	Button
} from '@material-ui/core'

import CustomTextField from './CustomTextField'
import CustomDropDown from './CustomDropDown'
import { IInstrument } from '../../../types/instrument.interface'
import { InstrumentService } from '../../../services/instrument.service'
import CustomNumberInput from './CustomNumberInput'
import { useQuery } from 'react-query'
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

// @ts-ignore
const UpdateInstrumentDBForm = ({ id }) => {
	const [error, setError] = useState('')
	const classes = useStyles()
	const [values, setValues] = useState<IInstrument>({
		InstrumentId: id,
		BrandId: 0,
		Name: '',
		Price: 0,
		Material: '',
		Type: '',
		Img: ''
	})
	const [errors, setErrors] = useState({
		priceNotNumber: '',
		priceNotPositiveNumber: '',
		brandNotNumber: '',
		brandNotPositiveNumber: ''
	})

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		validateFields()

		if (event.target.name == 'BrandId' || event.target.name == 'Price')
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
		InstrumentService.update(values).catch(error => {
			if (error.response) {
				setError(error.response.data.Error)
			}
		})
		console.log(values)

		console.log(error)
	}

	const validateFields = () => {
		let error = false
		if (isNaN(values.Price)) {
			error = true
			console.log('price not a number')

			setErrors(state => ({
				...state,
				priceNotNumber: 'Price should be integer number',
				priceNotPositiveNumber: ''
			}))
			console.log('priceNotNumber set')
			return error
		}
		if (values.Price < 0) {
			error = true
			console.log('price not a positive number')

			setErrors(state => ({
				...state,
				priceNotNumber: '',
				priceNotPositiveNumber: 'Price should be positive number'
			}))
			console.log('priceNotPositiveNumber set')
			return error
		}

		if (isNaN(values.BrandId)) {
			error = true
			console.log('brandId not a number')

			setErrors(state => ({
				...state,
				brandNotNumber: 'Brand id should be integer number',
				brandNotPositiveNumber: ''
			}))
			console.log('brandNotNumber set')
			return error
		}
		if (values.BrandId < 0) {
			error = true
			console.log('brandId not a positive number')

			setErrors(state => ({
				...state,
				brandNotNumber: '',
				brandNotPositiveNumber: 'Brand id should be positive number'
			}))
			console.log('brandNotPositiveNumber set')
			return error
		}
		setErrors(state => ({
			...state,
			priceNotNumber: '',
			priceNotPositiveNumber: '',
			brandNotNumber: '',
			brandNotPositiveNumber: ''
		}))
		console.log('no errors')
		return error
	}

	return (
		<Paper className={classes.container}>
			<Typography variant={'h4'} className={classes.title}>
				Update instrument in database
			</Typography>
			<form onSubmit={e => handleSubmit(e)} className={classes.form}>
				<CustomTextField
					changeHandler={handleChange}
					label={'Brand id'}
					error={
						Boolean(errors.brandNotNumber) ||
						Boolean(errors.brandNotPositiveNumber)
					}
					helperText={
						errors.brandNotNumber == ''
							? errors.brandNotPositiveNumber
							: errors.brandNotNumber
					}
					name={'BrandId'}
					hide={false}
				/>
				<CustomTextField
					changeHandler={handleChange}
					label={'Name'}
					error={false}
					helperText={''}
					name={'Name'}
					hide={false}
				/>
				<CustomTextField
					changeHandler={handleChange}
					label={'Price'}
					error={
						Boolean(errors.priceNotNumber) ||
						Boolean(errors.priceNotPositiveNumber)
					}
					helperText={
						errors.priceNotNumber == ''
							? errors.priceNotPositiveNumber
							: errors.priceNotNumber
					}
					name={'Price'}
					hide={false}
				/>
				<CustomTextField
					label={'Material'}
					name={'Material'}
					error={false}
					helperText={''}
					changeHandler={handleChange}
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
					label={'Img'}
					name={'Img'}
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
					Add
				</Button>
				<MessageBox type={error == 'no error' ? 'success' : 'error'}>
					{error == 'no error' ? 'success' : error}
				</MessageBox>
			</form>
		</Paper>
	)
}

export default UpdateInstrumentDBForm
