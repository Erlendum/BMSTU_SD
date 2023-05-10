import React, { FC, useEffect, useState } from 'react'
import { IInstrument } from '../../../types/instrument.interface'

import styles from './InstrumentItem.module.scss'
import { UserService } from '../../../services/user.service'
import { InstrumentService } from '../../../services/instrument.service'
import UpdateInstrumentDBForm from '../form/UpdateInstrumentDBForm'
import Modal from '../modal/Modal'
import { ComparisonListService } from '../../../services/comparisonList.service'
import { toast, ToastContainer } from 'react-toastify'

const InstrumentItem: FC<{
	instrument: IInstrument
	isComparisonList: boolean
}> = ({ instrument, isComparisonList }) => {
	const [error, setError] = useState('no error')
	let isLogin = UserService.getCurrentLogin() != null
	let isAdmin =
		UserService.getCurrentIsAdmin() != null &&
		UserService.getCurrentIsAdmin() == 'true'

	useEffect(() => {
		if (error != 'no error') {
			toast.error('ERROR ' + error, {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
	}, [error])

	const handleDelete = async () => {
		setError('no error')
		let isError = false
		await InstrumentService.delete(instrument.InstrumentId).catch(error => {
			isError = true
			if (error.response) {
				setError(error.response.data.Error)
				console.log(error.response.data.Error)
			}
		})
		if (!isError) {
			toast.success('Instrument was successfully deleted', {
				position: toast.POSITION.BOTTOM_LEFT
			})
		}
	}

	const handleAddInstrumentToCart = async () => {
		setError('no error')
		let isError = false
		ComparisonListService.addInstrument(instrument.InstrumentId).catch(
			error => {
				isError = true
				if (error.response) {
					setError(error.response.data.Error)
				}
			}
		)
		if (!isError) {
			toast.success('Instrument was successfully added', {
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
						hidden={isComparisonList}
						onClick={handleAddInstrumentToCart}
						className='fa fa-heart'
					></i>
				</a>
				<a href='javascript:void(0);' onClick={handleDelete}>
					<i hidden={!isAdmin || isComparisonList} className='fa fa-trash'></i>
				</a>
				<a
					href='javascript:void(0);'
					onClick={() => setUpdateInstrumentInDBModalActive(true)}
				>
					<i hidden={!isAdmin || isComparisonList} className='fa fa-edit'></i>
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
