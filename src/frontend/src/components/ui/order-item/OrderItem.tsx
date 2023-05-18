import React, { FC, useEffect, useState } from 'react'
import { IDiscount } from '../../../types/discount.interface'
import { UserService } from '../../../services/user.service'
import { toast } from 'react-toastify'
import { DiscountService } from '../../../services/discount.service'
import styles from '../discount-item/DiscountItem.module.scss'
import Modal from '../modal/Modal'
import UpdateDiscountDBForm from '../form/UpdateDiscountDBForm'
import { IOrder } from '../../../types/order.interface'
import { InstrumentService } from '../../../services/instrument.service'
import { OrderService } from '../../../services/order.service'

const OrderItem: FC<{
	order: IOrder
	updateQuery: boolean
	setUpdateQuery: any
}> = ({ order , updateQuery, setUpdateQuery}) => {
	const [error, setError] = useState('no error')
	const [UpdateDiscountInDBModalActive, setUpdateDiscountInDBModalActive] =
		useState(false)
	const [status, setStatus] = useState(order.Status)
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

	useEffect(() => {
		if (order.Status != status) {
			setError('no error')
			order.Status = status
			OrderService.update(order).catch(error => {
				if (error.response) {
					setError(error.response.data.Error)
				}
			})
			setUpdateQuery(!updateQuery)
		}
	}, [status])

	const displayStatus = isAdmin ? <select
		form="dorama-form"
		required={true}
		className="w-auto"
		onChange={(e)=>{setStatus(e.target.value)}}
		defaultValue={order.Status}
	>
		<option value="Created">Created</option>
		<option value="Delivering">Delivering</option>
		<option value="Delivered">Delivered</option>
	</select> : order.Status

	return (
		<tr>
			<td className={styles.textLeft}>{order.OrderId}</td>
			<td className={styles.textLeft}>{order.Time.substring(0, 10)}</td>
			<td className={styles.price}>{order.Price}</td>
			<td className={styles.textLeft}>{displayStatus}</td>
			<td hidden={!isAdmin} className={styles.textLeft}>
				{order.UserId}
			</td>
		</tr>
	)
}

export default OrderItem
