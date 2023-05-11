import React, { FC, useEffect, useState } from 'react'
import { IInstrument } from '../../../types/instrument.interface'
import { toast } from 'react-toastify'
import { ComparisonListService } from '../../../services/comparisonList.service'
import styles from '../discount-item/DiscountItem.module.scss'
import { IDiscount } from '../../../types/discount.interface'
import { DiscountService } from '../../../services/discount.service'
import { UserService } from '../../../services/user.service'
import UpdateInstrumentDBForm from '../form/UpdateInstrumentDBForm'
import Modal from '../modal/Modal'
import UpdateDiscountDBForm from '../form/UpdateDiscountDBForm'

const DiscountItem: FC<{
	discount: IDiscount
	updateQuery: boolean
	setUpdateQuery: any
}> = ({ discount , updateQuery, setUpdateQuery}) => {
	const [error, setError] = useState('no error')
	const [UpdateDiscountInDBModalActive, setUpdateDiscountInDBModalActive] =
		useState(false)
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
		await DiscountService.delete(discount.DiscountId).catch(error => {
			isError = true
			if (error.response) {
				setError(error.response.data.Error)
			}
		})
		if (!isError) {
			toast.success('Discount was successfully deleted', {
				position: toast.POSITION.BOTTOM_LEFT
			})
			setUpdateQuery(!updateQuery)
		}
		console.log(error)
	}

	return (
		<tr>
			<td className={styles.textLeft}>{discount.InstrumentId}</td>
			<td hidden={!isAdmin} className={styles.textLeft}>
				{discount.UserId}
			</td>
			<td className={styles.textLeft}>{discount.Amount}</td>
			<td className={styles.textLeft}>{discount.Type}</td>
			<td className={styles.textLeft}>{discount.DateBegin.substring(0, 10)}</td>
			<td className={styles.textLeft}>{discount.DateEnd.substring(0, 10)}</td>
			<td hidden={!isAdmin}>
				<a
					className={styles.links}
					href='javascript:void(0);'
					onClick={handleDelete}
				>
					<i className='fa fa-trash'></i>
				</a>
				<a
					className={styles.links}
					href='javascript:void(0);'
					onClick={() => setUpdateDiscountInDBModalActive(true)}
				>
					<i className='fa fa-edit'></i>
				</a>
			</td>
			<Modal
				active={UpdateDiscountInDBModalActive}
				setActive={setUpdateDiscountInDBModalActive}
			>
				<UpdateDiscountDBForm id={discount.DiscountId} />
			</Modal>
		</tr>
	)
}

export default DiscountItem
