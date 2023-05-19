import React, { FC, useEffect, useState } from 'react'
import { IInstrument } from '../../../types/instrument.interface'
import { UserService } from '../../../services/user.service'
import { toast, ToastContainer } from 'react-toastify'
import { InstrumentService } from '../../../services/instrument.service'
import { ComparisonListService } from '../../../services/comparisonList.service'
import styles from '../orderElement-item/OrderElementItem.module.scss'
import Modal from '../modal/Modal'
import UpdateInstrumentDBForm from '../form/UpdateInstrumentDBForm'
import { IOrder, IOrderElement } from '../../../types/order.interface'
import { OrderService } from '../../../services/order.service'

const OrderElementItem: FC<{
	orderElement: IOrderElement
}> = ({ orderElement}) => {
	const [UpdateDiscountInDBModalActive, setUpdateDiscountInDBModalActive] =
		useState(false)
	let isAdmin =
		UserService.getCurrentIsAdmin() != null &&
		UserService.getCurrentIsAdmin() == 'true'



	return (
		<tr>
			<td className={styles.textLeft}>{orderElement.InstrumentId}</td>
			<td className={styles.price}>{new Intl.NumberFormat('ru-RU', {
				style: 'currency',
				currency: 'RUB',
				maximumFractionDigits: 0
			}).format(orderElement.Price)}</td>
			<td className={styles.textLeft}>{orderElement.Amount}</td>
		</tr>
	)
}

export default OrderElementItem
