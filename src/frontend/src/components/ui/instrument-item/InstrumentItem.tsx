import React, { FC, useState } from 'react'
import { IInstrument } from '../../../types/instrument.interface'

import styles from './InstrumentItem.module.scss'
import { UserService } from '../../../services/user.service'
import { InstrumentService } from '../../../services/instrument.service'
import UpdateInstrumentDBForm from '../form/UpdateInstrumentDBForm'
import Modal from '../modal/Modal'
import { CartService } from '../../../services/cart.service'
import { toast, ToastContainer } from 'react-toastify'

const InstrumentItem: FC<{ instrument: IInstrument; isCart: boolean }> = ({
	instrument,
	isCart
}) => {
	const [error, setError] = useState('no error')
	let isLogin = UserService.getCurrentLogin() != null
	let isAdmin =
		UserService.getCurrentIsAdmin() != null &&
		UserService.getCurrentIsAdmin() == 'true'

	const handleDelete = async () => {
		await setError('no error')
		InstrumentService.delete(instrument.InstrumentId).catch(error => {
			if (error.response) {
				setError(error.response.data.Error)
			}
		})
		if (error == 'no error') {
			toast.success('Instrument was successfully deleted', {
				position: toast.POSITION.BOTTOM_LEFT
			})
		} else {
			toast.error('ERROR ' + error, {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
		console.log(error)
	}

	const handleDeleteCart = async () => {
		await setError('no error')
		CartService.deleteInstrument(instrument.InstrumentId).catch(error => {
			if (error.response) {
				setError(error.response.data.Error)
			}
		})
		if (error == 'no error') {
			toast.success('Instrument was successfully deleted from cart', {
				position: toast.POSITION.BOTTOM_LEFT
			})
		} else {
			toast.error('ERROR ' + error, {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
		console.log(error)
	}

	const handleAddInstrumentToCart = async () => {
		await setError('no error')
		CartService.addInstrument(instrument.InstrumentId).catch(error => {
			if (error.response) {
				setError(error.response.data.Error)
			}
		})
		if (error == 'no error') {
			toast.success('Instrument was successfully added to cart', {
				position: toast.POSITION.BOTTOM_LEFT
			})
		} else {
			toast.error('ERROR ' + error, {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
		console.log(error)
	}

	const [UpdateInstrumentInDBModalActive, setUpdateInstrumentInDBModalActive] =
		useState(false)

	return (
		<div className={styles.item}>
			<img src={instrument.Img} alt={instrument.Name} />
			<div className={styles.type}>{instrument.Type}</div>

			<div className={styles.heading}>{instrument.Name} </div>
			<div className={styles.price}>
				{new Intl.NumberFormat('ru-RU', {
					style: 'currency',
					currency: 'RUB',
					maximumFractionDigits: 0
				}).format(instrument.Price)}
			</div>
			<div className={styles.links} hidden={!isLogin}>
				<a href='javascript:void(0);'>
					<i
						hidden={isCart}
						onClick={handleAddInstrumentToCart}
						className='fa fa-shopping-cart'
					></i>
				</a>
				<a href='javascript:void(0);' onClick={handleDelete}>
					<i hidden={!isAdmin || isCart} className='fa fa-trash'></i>
				</a>
				<a href='javascript:void(0);' onClick={handleDeleteCart}>
					<i hidden={!isCart} className='fa fa-trash'></i>
				</a>
				<a
					href='javascript:void(0);'
					onClick={() => setUpdateInstrumentInDBModalActive(true)}
				>
					<i hidden={!isAdmin || isCart} className='fa fa-edit'></i>
				</a>
			</div>
			<Modal
				active={UpdateInstrumentInDBModalActive}
				setActive={setUpdateInstrumentInDBModalActive}
			>
				<UpdateInstrumentDBForm id={instrument.InstrumentId} />
			</Modal>
			<ToastContainer />
		</div>
	)
}

export default InstrumentItem
