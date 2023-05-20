import React, { FC, useEffect, useState } from 'react'
import { IDiscount } from '../../../types/discount.interface'
import { UserService } from '../../../services/user.service'
import { toast } from 'react-toastify'
import { DiscountService } from '../../../services/discount.service'
import styles from '../order-item/OrderItem.module.scss'
import Modal from '../modal/Modal'
import UpdateDiscountDBForm from '../form/UpdateDiscountDBForm'
import { IOrder } from '../../../types/order.interface'
import { InstrumentService } from '../../../services/instrument.service'
import { OrderService } from '../../../services/order.service'
import UpdateInstrumentDBForm from '../form/UpdateInstrumentDBForm'
import OrderElementItem from '../orderElement-item/OrderElementItem'
import { useQuery } from 'react-query'

const OrderItem: FC<{
	order: IOrder
	updateQuery: boolean
	setUpdateQuery: any
}> = ({ order , updateQuery, setUpdateQuery}) => {
	const [OrderInfoModalActive, setOrderInfoModalActive] =
		useState(false)

	const [error, setError] = useState('no error')
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
			order.UserId = 0
			order.Price = 0
			order.Time = "2006-01-02T15:04:05Z"
			order.Status = status
			OrderService.update(order).catch(error => {
				if (error.response) {
					setError(error.response.data.Error)
				}
			})
			setUpdateQuery(!updateQuery)
		}
	}, [status])

	const { data: order_elements,  refetch} = useQuery(
		['order_elements'],
		() => OrderService.getOrderElements(order.OrderId),
		{
			select: ({ order_elements }) => order_elements,
			refetchOnWindowFocus: true,
			staleTime: 0,
			cacheTime: 0,
			refetchInterval: 0,
		}
	)

	const handleClick = () => {
		refetch();
		setOrderInfoModalActive(true)
	};

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
			<td className={styles.price}>{new Intl.NumberFormat('ru-RU', {
				style: 'currency',
				currency: 'RUB',
				maximumFractionDigits: 0
			}).format(order.Price)}</td>
			<td className={styles.textLeft}>{displayStatus}</td>
			<td hidden={!isAdmin} className={styles.textLeft}>
				{order.UserId}
			</td>
			<td>
				<a
					className={styles.links}
					href='javascript:void(0);'
					onClick={handleClick}
				>
					<i className='fa fa-info-circle'></i>
				</a>
			</td>
			<Modal
				active={OrderInfoModalActive}
				setActive={setOrderInfoModalActive}
			>
				<table className={styles.table}>
					<thead>
					<tr>
						<th className={styles.textLeft}>Instrument Id</th>
						<th className={styles.textLeft}>Price</th>
						<th className={styles.textLeft}>Amount</th>
					</tr>
					</thead>
					<tbody className='table-hover'>
					{order_elements?.map(
						order_element => {
							return (
								<OrderElementItem
									key={order_element.OrderElementId}
									orderElement={order_element}
								/>
							)
						}
					)}
					</tbody>
				</table>
			</Modal>
		</tr>

	)
}

export default OrderItem
